package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// see readme for generating a local CA
func (a *App) loadCACert(caKeyFile, caCertFile string) (*tls.Certificate, error) {
	if caCertFile == "" {
		caCertFile = strings.TrimSuffix(caKeyFile, ".key") + ".crt"
	}
	ca, err := tls.LoadX509KeyPair(caCertFile, caKeyFile)
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func (a *App) loadTLSConfig(hostNames []string, caCert *tls.Certificate) (*tls.Config, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 30)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	// root
	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	leafTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"lightproxy"},
			CommonName:   "Proxy Cert",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              hostNames,
	}

	rootTemplate := leafTemplate
	var rootKey interface{} = leafKey
	if caCert != nil {
		rootTemplate, err = x509.ParseCertificate(caCert.Certificate[0])
		if err != nil {
			return nil, err
		}
		rootKey = caCert.PrivateKey
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader, leafTemplate, rootTemplate, &leafKey.PublicKey, rootKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	certBuf := &bytes.Buffer{}
	err = pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, err
	}
	keyBuf := &bytes.Buffer{}
	pem.Encode(keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(leafKey)})
	cert, err := tls.X509KeyPair(certBuf.Bytes(), keyBuf.Bytes())
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		NextProtos: []string{
			"http/1.1",
			"h2", // http/2
		},
		Certificates: []tls.Certificate{cert},
	}

	return tlsConfig, nil
}

func (a *App) startTLSProxy() error {
	hostnames := []string{}
	for _, e := range a.handlers {
		hostnames = append(hostnames, e.e.Source)
	}

	var ca *tls.Certificate
	var err error

	if a.config.UseMkcert {
		cmd := exec.Command("mkcert", "-CAROOT")
		mkcertOut, err := cmd.Output()
		if err != nil {
			return err
		}
		mkcertRoot := strings.TrimSpace(string(mkcertOut))
		if mkcertRoot == "" {
			return fmt.Errorf("No directory found for mkcert; have you run mkcert -install?")
		}
		caKeyFile := path.Join(mkcertRoot, "rootCA-key.pem")
		caCertFile := path.Join(mkcertRoot, "rootCA.pem")

		if _, err = os.Stat(caKeyFile); os.IsNotExist(err) {
			return fmt.Errorf("mkcert key file not found; have you run mkcert -install?")
		}
		if _, err = os.Stat(caCertFile); os.IsNotExist(err) {
			return fmt.Errorf("mkcert cert file not found; have you run mkcert -install?")
		}

		ca, err = a.loadCACert(caKeyFile, caCertFile)
		if err != nil {
			return err
		}
		fmt.Println("tls: using mkcert")
	} else if a.config.CAKeyFile != "" {
		ca, err = a.loadCACert(a.config.CAKeyFile, a.config.CACertFile)
		if err != nil {
			return err
		}
	}

	tlsConfig, err := a.loadTLSConfig(hostnames, ca)
	tlsAddr := a.config.TLSAddr
	if err != nil {
		return err
	}
	if tlsConfig == nil {
		return nil
	}
	l, err := tls.Listen("tcp", tlsAddr, tlsConfig)
	if err != nil {
		return err
	}
	tlsServer := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      a,
		TLSConfig: tlsConfig,
	}
	fmt.Println("tls: listening on", tlsAddr)
	go tlsServer.Serve(l)
	return nil
}
