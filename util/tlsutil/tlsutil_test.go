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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand/util/tlsutil/tlstest"
)

func TestTLSConfigPolicies(t *testing.T) {
	require.Nil(t, MakeServerConfig(false, false))
	serverConfig := MakeServerConfig(true, true)
	require.Equal(t, uint16(tls.VersionTLS13), serverConfig.MinVersion)
	require.Equal(t, []tls.CurveID{PQKeyExchange}, serverConfig.CurvePreferences)
	require.Nil(t, serverConfig.VerifyConnection)
	clientConfig := MakeClientConfig(false, true)
	require.NotNil(t, clientConfig)
	require.NotNil(t, clientConfig.VerifyConnection)
}

func TestCertificatePinningAndSignaturePolicy(t *testing.T) {
	certFile, keyFile := tlstest.CertificateFiles(t, "127.0.0.1", "localhost")
	require.NoError(t, CheckServerCertificate(certFile, keyFile))

	serverConfig, err := LoadServerConfig(certFile, keyFile, true, true)
	require.NoError(t, err)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	server.TLS = serverConfig
	server.StartTLS()
	defer server.Close()

	clientConfig, err := WithTrustedCertificate(MakeClientConfig(true, true), certFile)
	require.NoError(t, err)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: clientConfig}}
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
	require.NoError(t, response.Body.Close())

	classicalServer := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer classicalServer.Close()
	classicalCertFile := filepath.Join(t.TempDir(), "classical.pem")
	classicalCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: classicalServer.Certificate().Raw})
	require.NoError(t, os.WriteFile(classicalCertFile, classicalCertPEM, 0600))
	classicalClientConfig, err := WithTrustedCertificate(MakeClientConfig(false, true), classicalCertFile)
	require.NoError(t, err)
	classicalClient := &http.Client{Transport: &http.Transport{TLSClientConfig: classicalClientConfig}}
	_, err = classicalClient.Get(classicalServer.URL)
	require.Error(t, err)
	require.Error(t, CheckServerCertificate(classicalCertFile, filepath.Join(t.TempDir(), "missing.pem")))
	_, err = LoadServerConfig(classicalCertFile, mustWriteClassicalKey(t, classicalServer), false, true)
	require.Error(t, err)
}

func mustWriteClassicalKey(t *testing.T, server *httptest.Server) string {
	t.Helper()
	keyDER, err := x509.MarshalPKCS8PrivateKey(server.TLS.Certificates[0].PrivateKey)
	require.NoError(t, err)
	keyFile := filepath.Join(t.TempDir(), "classical-key.pem")
	require.NoError(t, os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}), 0600))
	return keyFile
}
