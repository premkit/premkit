package certs

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/premkit/premkit/log"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

// GenerateSelfSigned will generate a new, self-signed cert in memory and return
// the location on disk for the key and cert
func GenerateSelfSigned(tlsStore string) (string, string, error) {
	// If a key and cert exist in the store, return them, don't generate new
	if _, err := os.Stat(path.Join(tlsStore, "key.pem")); err == nil {
		if _, err := os.Stat(path.Join(tlsStore, "cert.pem")); err == nil {
			return path.Join(tlsStore, "key.pem"), path.Join(tlsStore, "cert.pem"), nil
		}
	}

	log.Infof("Generating a new self-signed cert to use for https connections")

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{""},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	if err := os.MkdirAll(tlsStore, 0755); err != nil {
		log.Error(err)
		return "", "", err
	}

	certOut, err := os.Create(path.Join(tlsStore, "cert.pem"))
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(path.Join(tlsStore, "key.pem"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()

	if err := ParseKeyPair(path.Join(tlsStore, "key.pem"), path.Join(tlsStore, "cert.pem")); err != nil {
		log.Error(err)
		return "", "", err
	}

	log.Debugf("Self-signed cert and key created in %s", tlsStore)
	return path.Join(tlsStore, "key.pem"), path.Join(tlsStore, "cert.pem"), nil
}

// ParseKeyPair will parse a key and cert from the filesystem and return
// an error or nil to determine if it validates
func ParseKeyPair(keyFile string, certFile string) error {
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Error(err)
		return err
	}

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Error(err)
		return err
	}

	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Error(err)
		return err
	}

	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Error(err)
		return err
	}

	// Validate this is an accepted key and cert
	if _, err := tls.X509KeyPair(certData, keyData); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
