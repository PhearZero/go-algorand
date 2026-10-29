// Copyright (C) 2019-2026 Algorand Foundation Ltd.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const PQKeyExchange = tls.X25519MLKEM768

// MakeServerConfig creates the server-side TLS policy configuration.
func MakeServerConfig(requireKEX, requireSig bool) *tls.Config {
	if !requireKEX && !requireSig {
		return nil
	}
	cfg := &tls.Config{MinVersion: tls.VersionTLS13}
	if requireKEX {
		cfg.CurvePreferences = []tls.CurveID{PQKeyExchange}
	}
	return cfg
}

// MakeClientConfig creates the client-side TLS policy configuration.
func MakeClientConfig(requireKEX, requireSig bool) *tls.Config {
	cfg := MakeServerConfig(requireKEX, requireSig)
	if cfg != nil && requireSig {
		cfg.VerifyConnection = CheckPeerCertificate
	}
	return cfg
}

// CheckKeyExchange validates the negotiated TLS key exchange when required.
func CheckKeyExchange(state *tls.ConnectionState, requireKEX bool) error {
	if !requireKEX {
		return nil
	}
	if state == nil {
		return fmt.Errorf("TLS connection has no connection state")
	}
	if state.CurveID != PQKeyExchange {
		return fmt.Errorf("TLS key exchange %d is not %d", state.CurveID, PQKeyExchange)
	}
	return nil
}

// CheckPeerCertificate requires the peer leaf certificate to use ML-DSA.
func CheckPeerCertificate(state tls.ConnectionState) error {
	if len(state.PeerCertificates) == 0 {
		return fmt.Errorf("TLS connection has no peer certificate")
	}
	if state.PeerCertificates[0].PublicKeyAlgorithm != x509.MLDSA {
		return fmt.Errorf("TLS peer certificate key algorithm %s is not ML-DSA", state.PeerCertificates[0].PublicKeyAlgorithm)
	}
	return nil
}

// CheckServerCertificate validates the certificate pair and optionally requires ML-DSA.
func CheckServerCertificate(certFile, keyFile string) error {
	pair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("loading TLS certificate: %w", err)
	}
	if len(pair.Certificate) == 0 {
		return fmt.Errorf("TLS certificate file contains no certificate")
	}
	leaf, err := x509.ParseCertificate(pair.Certificate[0])
	if err != nil {
		return fmt.Errorf("parsing TLS certificate: %w", err)
	}
	return checkMLDSACertificate(leaf)
}

// LoadServerConfig loads a certificate and applies the requested TLS policies.
func LoadServerConfig(certFile, keyFile string, requireKEX, requireSig bool) (*tls.Config, error) {
	pair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading TLS certificate: %w", err)
	}
	if len(pair.Certificate) == 0 {
		return nil, fmt.Errorf("TLS certificate file contains no certificate")
	}
	if requireSig {
		leaf, parseErr := x509.ParseCertificate(pair.Certificate[0])
		if parseErr != nil {
			return nil, fmt.Errorf("parsing TLS certificate: %w", parseErr)
		}
		if err = checkMLDSACertificate(leaf); err != nil {
			return nil, err
		}
		pair.Leaf = leaf
	}
	cfg := MakeServerConfig(requireKEX, requireSig)
	if cfg == nil {
		cfg = &tls.Config{}
	}
	cfg.Certificates = []tls.Certificate{pair}
	return cfg, nil
}

// WithTrustedCertificate returns a client config that trusts only the given certificate.
func WithTrustedCertificate(cfg *tls.Config, certFile string) (*tls.Config, error) {
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("reading trusted TLS certificate: %w", err)
	}
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("trusted TLS certificate file contains no certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing trusted TLS certificate: %w", err)
	}
	if cfg == nil {
		cfg = &tls.Config{}
	} else {
		cfg = cfg.Clone()
	}
	roots := x509.NewCertPool()
	roots.AddCert(cert)
	cfg.RootCAs = roots
	return cfg, nil
}

func checkMLDSACertificate(cert *x509.Certificate) error {
	if cert.PublicKeyAlgorithm != x509.MLDSA {
		return fmt.Errorf("TLS certificate key algorithm %s is not ML-DSA", cert.PublicKeyAlgorithm)
	}
	return nil
}
