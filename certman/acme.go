package certman

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"golang.org/x/crypto/acme"

	"github.com/open-tdp/go-helper/logman"
)

func CreateCert() error {

	domains := []string{}

	ctx := context.Background()
	client := newClient(ctx, "--url--")

	for _, domain := range domains {
		authz, err := client.Authorize(ctx, domain)

		if err != nil {
			logman.Error("authorize failed", "domain", domain, "error", err)
			return err
		}

		// Already authorized.
		if authz.Status == acme.StatusValid {
			continue
		}

		// Pick the DNS challenge, if any.
		var chal *acme.Challenge
		for _, c := range authz.Challenges {
			if c.Type == "dns-01" {
				chal = c
				break
			}
		}

		if chal == nil {
			logman.Error("no dns-01 challenge", "domain", domain)
			return err
		}

		// Fulfill the challenge
		val, err := client.DNS01ChallengeRecord(chal.Token)
		if err != nil {
			logman.Error("get dns-01 value failed", "domain", domain, "error", err)
			return err
		}

		// Add a TXT record containing the val value under "_acme-challenge" name
		if err := createRecord(ctx, domain, val); err != nil {
			logman.Error("dns update failed", "domain", domain, "error", err)
			return err
		}

		// Let CA know we're ready. But are we? Is DNS propagated yet?
		if _, err := client.Accept(ctx, chal); err != nil {
			logman.Error("dns-01 accept failed", "domain", domain, "error", err)
			return err
		}

		// Wait for the CA to validate.
		if _, err := client.WaitAuthorization(ctx, authz.URI); err != nil {
			logman.Error("dns-01 authorization failed", "domain", domain, "error", err)
			return err
		}
	}

	// All authorizations are granted. Request the certificate.

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logman.Error("generate keypair failed", "error", err)
		return err
	}

	req := &x509.CertificateRequest{
		DNSNames: domains,
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, req, key)
	if err != nil {
		logman.Error("create csr failed", "error", err)
		return err
	}

	der, url, err := client.CreateOrderCert(ctx, "--url-", csr, true)
	if err != nil {
		logman.Error("create cert failed", "error", err)
		return err
	}

	fmt.Println(der, url)

	return nil

}

func newClient(ctx context.Context, url string) *acme.Client {

	akey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logman.Error("Generate keypair failed", "error", err)
	}

	client := &acme.Client{Key: akey, DirectoryURL: url}
	account := &acme.Account{}

	account, err = client.Register(ctx, account, acme.AcceptTOS)
	if err != nil {
		logman.Error("Creates account failed", "error", err)
	}

	fmt.Println(account)

	return client

}
