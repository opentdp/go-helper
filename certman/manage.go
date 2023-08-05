package certman

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/libdns/libdns"
	"github.com/opentdp/go-helper/logman"
	"golang.org/x/crypto/acme"
	"golang.org/x/net/idna"
)

type Manager struct {
	Email                  string
	DirectoryUrl           string
	ExternalAccountBinding *acme.ExternalAccountBinding

	Cache  Cache
	Logger *logman.Logger

	DnsProvider interface {
		libdns.RecordAppender
		libdns.RecordDeleter
	}

	client   *acme.Client
	clientMu sync.Mutex

	state   map[certKey]*certState
	stateMu sync.Mutex
}

func (m *Manager) GetCertificate(name string) (*tls.Certificate, error) {

	if name == "" {
		return nil, errors.New("missing domain name")
	}

	name, err := idna.Lookup.ToASCII(name)
	if err != nil {
		return nil, errors.New("domain name contains invalid character")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ck := certKey{
		domain: name,
	}

	// read cache
	cert, err := m.loadCert(ctx, ck)
	if err == nil {
		return cert, nil
	}
	if err != ErrCacheMiss {
		return nil, err
	}

	// first time
	cert, err = m.createCert(ctx, ck)
	if err != nil {
		return nil, err
	}

	m.pemSave(ctx, ck, cert)
	return cert, nil

}

func (m *Manager) loadCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {

	m.stateMu.Lock()

	if s, ok := m.state[ck]; ok {
		m.stateMu.Unlock()
		s.RLock()
		defer s.RUnlock()
		return s.tlscert()
	}
	defer m.stateMu.Unlock()

	cert, err := m.pemLoad(ctx, ck)
	if err != nil {
		return nil, err
	}

	signer, ok := cert.PrivateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("private key cannot sign")
	}
	if m.state == nil {
		m.state = make(map[certKey]*certState)
	}

	s := &certState{
		key:  signer,
		cert: cert.Certificate,
		leaf: cert.Leaf,
	}
	m.state[ck] = s

	return cert, nil

}

func (m *Manager) createCert(ctx context.Context, ck certKey) (*tls.Certificate, error) {

	state, err := m.certState(ck)
	if err != nil {
		return nil, err
	}

	if !state.locked {
		state.RLock()
		defer state.RUnlock()
		return state.tlscert()
	}

	defer state.Unlock()
	state.locked = false

	der, leaf, err := m.authorizedCert(ctx, state.key, ck)
	if err != nil {
		time.AfterFunc(time.Minute, func() {
			m.stateMu.Lock()
			defer m.stateMu.Unlock()
			s, ok := m.state[ck]
			if !ok {
				return
			}
			if _, err := validCertificate(ck, s.cert, s.key, time.Now()); err == nil {
				return
			}
			delete(m.state, ck)
		})
		return nil, err
	}

	state.cert = der
	state.leaf = leaf

	return state.tlscert()

}

func (m *Manager) certState(ck certKey) (*certState, error) {

	m.stateMu.Lock()
	defer m.stateMu.Unlock()

	if m.state == nil {
		m.state = make(map[certKey]*certState)
	}
	if state, ok := m.state[ck]; ok {
		return state, nil
	}

	// new locked state
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	state := &certState{
		key:    key,
		locked: true,
	}
	state.Lock() // will be unlocked by m.certState caller

	m.state[ck] = state

	return state, nil

}

func (m *Manager) authorizedCert(ctx context.Context, key crypto.Signer, ck certKey) (der [][]byte, leaf *x509.Certificate, err error) {

	req := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: ck.domain},
		DNSNames: []string{ck.domain},
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, req, key)
	if err != nil {
		return nil, nil, err
	}

	m.Logger.Info("create client", "domain", ck.domain)
	client, err := m.acmeClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	m.Logger.Info("authorize order", "domain", ck.domain)
	order, err := m.authorizedOrder(ctx, ck)
	if err != nil {
		return nil, nil, err
	}

	m.Logger.Info("finalize order", "domain", ck.domain)
	chain, _, err := client.CreateOrderCert(ctx, order.FinalizeURL, csr, true)
	if err != nil {
		return nil, nil, err
	}

	m.Logger.Info("verify certificate", "domain", ck.domain)
	leaf, err = validCertificate(ck, chain, key, time.Now())
	if err != nil {
		return nil, nil, err
	}

	return chain, leaf, nil

}

func (m *Manager) authorizedOrder(ctx context.Context, ck certKey) (*acme.Order, error) {

	order, err := m.client.AuthorizeOrder(ctx, acme.DomainIDs(ck.domain))
	if err != nil {
		return nil, err
	}

	defer func() { go m.revokePendingAuthz(order.AuthzURLs) }()

	if order.Status == acme.StatusReady {
		return order, nil
	}
	if order.Status != acme.StatusPending {
		return nil, errors.New("invalid new order status " + order.Status)
	}

	for _, zurl := range order.AuthzURLs {
		m.Logger.Info("authorizing domain", "domain", ck.domain, "authz_url", zurl)

		authz, err := m.client.GetAuthorization(ctx, zurl)
		if err != nil {
			return nil, err
		}
		if authz.Status != acme.StatusPending {
			continue
		}

		var chal *acme.Challenge
		for _, c := range authz.Challenges {
			if c.Type == "dns-01" {
				chal = c
				break
			}
		}
		if chal == nil {
			return nil, errors.New("no viable challenge type found")
		}

		if cleanup, err := m.fulfill(ctx, chal, ck.domain); err != nil {
			return nil, err
		} else {
			defer cleanup()
		}

		if _, err := m.client.Accept(ctx, chal); err != nil {
			return nil, err
		}

		m.Logger.Info("waiting for authorization to be valid", "uri", authz.URI)
		if _, err := m.client.WaitAuthorization(ctx, authz.URI); err != nil {
			return nil, err
		}
	}

	return m.client.WaitOrder(ctx, order.URI)

}

func (m *Manager) revokePendingAuthz(uri []string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, u := range uri {
		authz, err := m.client.GetAuthorization(ctx, u)
		if err == nil && authz.Status == acme.StatusPending {
			m.client.RevokeAuthorization(ctx, u)
		}
	}

}

func (m *Manager) fulfill(ctx context.Context, chal *acme.Challenge, domain string) (func(), error) {

	value, err := m.client.DNS01ChallengeRecord(chal.Token)
	if err != nil {
		return nil, err
	}

	record := []libdns.Record{
		{Type: "TXT", Name: "_acme-challenge", Value: value},
	}

	if _, err = m.DnsProvider.AppendRecords(ctx, domain, record); err != nil {
		return nil, err
	}

	m.Logger.Info("wait 30 seconds for dns to take effect")
	time.Sleep(30 * time.Second)

	return func() {
		go m.DnsProvider.DeleteRecords(ctx, domain, record)
	}, nil

}

func (m *Manager) accountKey(ctx context.Context) (crypto.Signer, error) {

	const keyName = "account.key"

	genKey := func() (*ecdsa.PrivateKey, error) {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	if m.Cache == nil {
		return genKey()
	}

	data, err := m.Cache.Get(ctx, keyName)
	if err != nil {
		if err == ErrCacheMiss {
			key, err := genKey()
			if err != nil {
				return nil, err
			}
			var buf bytes.Buffer
			if err := encodeECDSAKey(&buf, key); err != nil {
				return nil, err
			}
			if err := m.Cache.Put(ctx, keyName, buf.Bytes()); err != nil {
				return nil, err
			}
			return key, nil
		}
		return nil, err
	}

	priv, _ := pem.Decode(data)
	if priv == nil || !strings.Contains(priv.Type, "PRIVATE") {
		return nil, errors.New("invalid account key found in cache")
	}

	return parsePrivateKey(priv.Bytes)

}

func (m *Manager) acmeClient(ctx context.Context) (*acme.Client, error) {

	m.clientMu.Lock()
	defer m.clientMu.Unlock()

	if m.client != nil {
		return m.client, nil
	}

	accountKey, err := m.accountKey(ctx)
	if err != nil {
		return nil, err
	}

	client := &acme.Client{
		DirectoryURL: m.DirectoryUrl,
		UserAgent:    "autocert",
		Key:          accountKey,
	}

	account := &acme.Account{
		Contact:                []string{"mailto:" + m.Email},
		ExternalAccountBinding: m.ExternalAccountBinding,
	}

	_, err = client.Register(ctx, account, acme.AcceptTOS)
	if ae, ok := err.(*acme.Error); err == nil ||
		err == acme.ErrAccountAlreadyExists || (ok && ae.StatusCode == http.StatusConflict) {
		m.client = client
		err = nil
	}

	return m.client, err

}
