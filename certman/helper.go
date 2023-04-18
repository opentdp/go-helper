package certman

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"sync"
	"time"
)

type certKey struct {
	domain string
}

func (c certKey) String() string {
	return c.domain
}

type certState struct {
	sync.RWMutex
	locked bool              // locked for read/write
	key    crypto.Signer     // private key for cert
	cert   [][]byte          // DER encoding
	leaf   *x509.Certificate // parsed cert[0]; always non-nil if cert != nil
}

func (s *certState) tlscert() (*tls.Certificate, error) {

	if s.key == nil {
		return nil, errors.New("missing signer")
	}
	if len(s.cert) == 0 {
		return nil, errors.New("missing certificate")
	}

	return &tls.Certificate{
		PrivateKey:  s.key,
		Certificate: s.cert,
		Leaf:        s.leaf,
	}, nil

}

func encodeECDSAKey(w io.Writer, key *ecdsa.PrivateKey) error {

	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}

	pb := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	return pem.Encode(w, pb)

}

func parsePrivateKey(der []byte) (crypto.Signer, error) {

	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			return key, nil
		case *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("unknown private key type in PKCS#8 wrapping")
		}
	}

	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("failed to parse private key")

}

func validCertificate(ck certKey, der [][]byte, key crypto.Signer, now time.Time) (*x509.Certificate, error) {

	var n int
	for _, b := range der {
		n += len(b)
	}
	pub := make([]byte, n)
	n = 0
	for _, b := range der {
		n += copy(pub[n:], b)
	}

	x509Cert, err := x509.ParseCertificates(pub)
	if err != nil || len(x509Cert) == 0 {
		return nil, errors.New("no public key found")
	}

	leaf := x509Cert[0]
	if now.Before(leaf.NotBefore) {
		return nil, errors.New("certificate is not valid yet")
	}
	if now.After(leaf.NotAfter) {
		return nil, errors.New("expired certificate")
	}
	if err := leaf.VerifyHostname(ck.domain); err != nil {
		return nil, err
	}

	// ensure the leaf corresponds to the private key and matches the certKey type
	switch pub := leaf.PublicKey.(type) {
	case *rsa.PublicKey:
		prv, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key type does not match public key type")
		}
		if pub.N.Cmp(prv.N) != 0 {
			return nil, errors.New("private key does not match public key")
		}
		return nil, errors.New("key type does not match expected value")
	case *ecdsa.PublicKey:
		prv, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key type does not match public key type")
		}
		if pub.X.Cmp(prv.X) != 0 || pub.Y.Cmp(prv.Y) != 0 {
			return nil, errors.New("private key does not match public key")
		}
	default:
		return nil, errors.New("unknown public key algorithm")
	}

	return leaf, nil

}
