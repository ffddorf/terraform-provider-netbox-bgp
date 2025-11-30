package provider

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

func loadCACertsFromFile(path string) (*x509.CertPool, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(bytes); !ok {
		return nil, errors.New("unable to parse any PEM encoded certificates from file")
	}
	return pool, nil
}
