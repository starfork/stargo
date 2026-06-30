package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

// NewServerTransportCredentials creates server TLS credentials from cert/key files.
// If caFile is non-empty, enables mTLS (requires client certificate signed by the CA).
func NewServerTransportCredentials(certFile, keyFile, caFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.NoClientCert,
	}

	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, err
		}
		cfg.ClientCAs = caCertPool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return credentials.NewTLS(cfg), nil
}

// NewClientTransportCredentials creates client TLS credentials.
// If caFile is non-empty, verifies the server certificate against the CA.
// If certFile/keyFile are non-empty, provides client certificate for mTLS.
func NewClientTransportCredentials(caFile, certFile, keyFile, serverName string) (credentials.TransportCredentials, error) {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: serverName,
	}

	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, err
		}
		cfg.RootCAs = caCertPool
	}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = []tls.Certificate{cert}
	}

	return credentials.NewTLS(cfg), nil
}
