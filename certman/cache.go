package certman

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type Cache = autocert.Cache

type DirCache = autocert.DirCache

var ErrCacheMiss = autocert.ErrCacheMiss

func (m *Manager) pemLoad(ctx context.Context, ck certKey) (*tls.Certificate, error) {

	if m.Cache == nil {
		return nil, ErrCacheMiss
	}

	data, err := m.Cache.Get(ctx, ck.String()+".pem")
	if err != nil {
		return nil, err
	}

	priv, pub := pem.Decode(data)
	if priv == nil || !strings.Contains(priv.Type, "PRIVATE") {
		return nil, ErrCacheMiss
	}

	privKey, err := parsePrivateKey(priv.Bytes)
	if err != nil {
		return nil, err
	}

	var pubDer [][]byte
	for len(pub) > 0 {
		var b *pem.Block
		b, pub = pem.Decode(pub)
		if b == nil {
			break
		}
		pubDer = append(pubDer, b.Bytes)
	}
	if len(pub) > 0 {
		return nil, ErrCacheMiss
	}

	leaf, err := validCertificate(ck, pubDer, privKey, time.Now())
	if err != nil {
		return nil, ErrCacheMiss
	}

	tlscert := &tls.Certificate{
		Certificate: pubDer,
		PrivateKey:  privKey,
		Leaf:        leaf,
	}

	return tlscert, nil

}

func (m *Manager) pemSave(ctx context.Context, ck certKey, tlscert *tls.Certificate) error {

	if m.Cache == nil {
		return nil
	}

	var buf bytes.Buffer

	switch key := tlscert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		if err := encodeECDSAKey(&buf, key); err != nil {
			return err
		}
	case *rsa.PrivateKey:
		b := x509.MarshalPKCS1PrivateKey(key)
		pb := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
		if err := pem.Encode(&buf, pb); err != nil {
			return err
		}
	default:
		return errors.New("unknown private key type")
	}

	for _, b := range tlscert.Certificate {
		pb := &pem.Block{Type: "CERTIFICATE", Bytes: b}
		if err := pem.Encode(&buf, pb); err != nil {
			return err
		}
	}

	return m.Cache.Put(ctx, ck.String()+".pem", buf.Bytes())

}
