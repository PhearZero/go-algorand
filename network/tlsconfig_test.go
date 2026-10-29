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

package network

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand/test/partitiontest"
	"github.com/algorand/go-algorand/util/tlsutil/tlstest"
)

func TestMakeTLSConfig(t *testing.T) {
	partitiontest.PartitionTest(t)

	require.Nil(t, makeTLSConfig(false, false))

	cfg := makeTLSConfig(true, false)
	require.Equal(t, uint16(tls.VersionTLS13), cfg.MinVersion)
	require.Equal(t, []tls.CurveID{tls.X25519MLKEM768}, cfg.CurvePreferences)

	sigConfig := makeTLSConfig(false, true)
	require.Equal(t, uint16(tls.VersionTLS13), sigConfig.MinVersion)
	require.Empty(t, sigConfig.CurvePreferences)
	require.Nil(t, sigConfig.VerifyConnection)
	clientConfig := makeClientTLSConfig(false, true)
	if clientConfig == nil {
		t.Fatal("expected client TLS config")
	}
	require.NotNil(t, clientConfig.VerifyConnection)
}

func TestCheckTLSKeyExchange(t *testing.T) {
	partitiontest.PartitionTest(t)

	require.NoError(t, checkTLSKeyExchange(nil, false))
	require.Error(t, checkTLSKeyExchange(nil, true))
	require.Error(t, checkTLSKeyExchange(&tls.ConnectionState{CurveID: tls.X25519}, true))
	require.NoError(t, checkTLSKeyExchange(&tls.ConnectionState{CurveID: pqKeyExchange}, true))
}

func TestTLSRequirePostQuantumKeyExchange(t *testing.T) {
	partitiontest.PartitionTest(t)

	curve := make(chan tls.CurveID, 1)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		curve <- r.TLS.CurveID
		w.WriteHeader(http.StatusNoContent)
	}))
	server.TLS = makeTLSConfig(true, false)
	server.StartTLS()
	defer server.Close()

	roots := x509.NewCertPool()
	certificate, err := x509.ParseCertificate(server.TLS.Certificates[0].Certificate[0])
	require.NoError(t, err)
	roots.AddCert(certificate)

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots}}}
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusNoContent, response.StatusCode)
	require.Equal(t, pqKeyExchange, <-curve)

	classicalClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		RootCAs:          roots,
		CurvePreferences: []tls.CurveID{tls.X25519},
	}}}
	_, err = classicalClient.Get(server.URL)
	require.Error(t, err)

	classicalServer := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	classicalServer.TLS = &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519}}
	classicalServer.StartTLS()
	defer classicalServer.Close()

	classicalRoots := x509.NewCertPool()
	classicalCertificate, err := x509.ParseCertificate(classicalServer.TLS.Certificates[0].Certificate[0])
	require.NoError(t, err)
	classicalRoots.AddCert(classicalCertificate)
	pqClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		RootCAs:          classicalRoots,
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{pqKeyExchange},
	}}}
	_, err = pqClient.Get(classicalServer.URL)
	require.Error(t, err)
}

func TestTLSRequirePostQuantumSignature(t *testing.T) {
	partitiontest.PartitionTest(t)

	pair := tlstest.Certificate(t, "127.0.0.1", "localhost")
	require.Equal(t, x509.MLDSA, pair.Leaf.PublicKeyAlgorithm)

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	server.TLS = makeTLSConfig(true, true)
	server.TLS.Certificates = []tls.Certificate{pair}
	server.StartTLS()
	defer server.Close()

	roots := x509.NewCertPool()
	roots.AddCert(pair.Leaf)
	clientTLS := makeClientTLSConfig(true, true)
	if clientTLS == nil {
		t.Fatal("expected client TLS config")
	}
	clientTLS.RootCAs = roots
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: clientTLS}}
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
	require.NoError(t, response.Body.Close())

	classicalServer := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer classicalServer.Close()
	classicalRoots := x509.NewCertPool()
	classicalRoots.AddCert(classicalServer.Certificate())
	classicalClientTLS := makeClientTLSConfig(false, true)
	if classicalClientTLS == nil {
		t.Fatal("expected client TLS config")
	}
	classicalClientTLS.RootCAs = classicalRoots
	classicalClient := &http.Client{Transport: &http.Transport{TLSClientConfig: classicalClientTLS}}
	_, err = classicalClient.Get(classicalServer.URL)
	require.Error(t, err)
}

func TestCheckTLSServerCertificate(t *testing.T) {
	partitiontest.PartitionTest(t)

	certFile, keyFile := tlstest.CertificateFiles(t, "relay.example")
	require.NoError(t, checkTLSServerCertificate(certFile, keyFile))
	require.Error(t, checkTLSServerCertificate(certFile, filepath.Join(t.TempDir(), "missing.pem")))

	classicalServer := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer classicalServer.Close()
	classicalCertificate := classicalServer.TLS.Certificates[0]
	classicalKeyDER, err := x509.MarshalPKCS8PrivateKey(classicalCertificate.PrivateKey)
	require.NoError(t, err)
	classicalCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: classicalCertificate.Certificate[0]})
	classicalKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: classicalKeyDER})
	dir := t.TempDir()
	classicalCertFile := filepath.Join(dir, "cert.pem")
	classicalKeyFile := filepath.Join(dir, "key.pem")
	require.NoError(t, os.WriteFile(classicalCertFile, classicalCertPEM, 0600))
	require.NoError(t, os.WriteFile(classicalKeyFile, classicalKeyPEM, 0600))
	require.Error(t, checkTLSServerCertificate(classicalCertFile, classicalKeyFile))
}
