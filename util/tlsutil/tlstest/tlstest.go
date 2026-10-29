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

package tlstest

import (
	"crypto/mldsa"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Certificate(t testing.TB, hosts ...string) tls.Certificate {
	t.Helper()
	var dnsNames []string
	var ipAddresses []net.IP
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		if ip := net.ParseIP(host); ip != nil {
			ipAddresses = append(ipAddresses, ip)
		} else {
			dnsNames = append(dnsNames, host)
		}
	}
	if len(dnsNames)+len(ipAddresses) == 0 {
		t.Fatal("at least one certificate host is required")
	}

	privateKey, err := mldsa.GenerateKey(mldsa.MLDSA65())
	if err != nil {
		t.Fatalf("generating ML-DSA key: %v", err)
	}
	now := time.Now()
	template := &x509.Certificate{
		Subject:               pkix.Name{CommonName: commonName(dnsNames, ipAddresses)},
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.Add(24 * time.Hour),
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certificateDER, err := x509.CreateCertificate(rand.Reader, template, template, privateKey.Public(), privateKey)
	if err != nil {
		t.Fatalf("creating ML-DSA certificate: %v", err)
	}
	leaf, err := x509.ParseCertificate(certificateDER)
	if err != nil {
		t.Fatalf("parsing ML-DSA certificate: %v", err)
	}
	return tls.Certificate{Certificate: [][]byte{certificateDER}, PrivateKey: privateKey, Leaf: leaf}
}

func CertificateFiles(t testing.TB, hosts ...string) (certFile, keyFile string) {
	t.Helper()
	certificate := Certificate(t, hosts...)
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(certificate.PrivateKey)
	if err != nil {
		t.Fatalf("encoding ML-DSA private key: %v", err)
	}
	dir := t.TempDir()
	certFile = filepath.Join(dir, "cert.pem")
	keyFile = filepath.Join(dir, "key.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Certificate[0]})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})
	if err = os.WriteFile(certFile, certPEM, 0600); err != nil {
		t.Fatalf("writing ML-DSA certificate: %v", err)
	}
	if err = os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		t.Fatalf("writing ML-DSA private key: %v", err)
	}
	return certFile, keyFile
}

func commonName(dnsNames []string, ipAddresses []net.IP) string {
	if len(dnsNames) != 0 {
		return dnsNames[0]
	}
	return ipAddresses[0].String()
}
