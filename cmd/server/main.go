package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"log"
	"math/big"
	"time"

	"github.com/marten-seemann/masque-go"
)

func main() {
	addr := flag.String("addr", "", "MASQUE server (IP:port)")
	flag.Parse()

	if len(*addr) == 0 {
		log.Fatal("missing MASQUE server")
	}

	tlsConf, err := generateTLSConfig()
	if err != nil {
		log.Fatalf("Generating TLS config failed: %s", err)
	}
	server := masque.NewServer(*addr, tlsConf)
	server.Serve()
}

func generateTLSConfig() (*tls.Config, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Acme Co"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{{
			PrivateKey:  priv,
			Certificate: [][]byte{derBytes}},
		},
	}, nil
}
