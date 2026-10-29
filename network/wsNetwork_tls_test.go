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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand/network/phonebook"
	"github.com/algorand/go-algorand/test/partitiontest"
	"github.com/algorand/go-algorand/util/tlsutil/tlstest"
)

func TestWebsocketNetworkRequiresPQKeyExchange(t *testing.T) {
	partitiontest.PartitionTest(t)

	certFile, keyFile, roots := writeTestTLSFiles(t)
	serverConfig := defaultConfig
	serverConfig.GossipFanout = 1
	serverConfig.TLSCertFile = certFile
	serverConfig.TLSKeyFile = keyFile
	serverConfig.TLSRequirePostQuantumKeyExchange = true
	serverConfig.TLSRequirePostQuantumSignature = true
	server := makeTestWebsocketNodeWithConfig(t, serverConfig)
	require.NoError(t, server.Start())
	defer netStop(t, server, "PQ server")

	address, listening := server.Address()
	require.True(t, listening)

	clientConfig := defaultConfig
	clientConfig.GossipFanout = 1
	clientConfig.TLSRequirePostQuantumKeyExchange = true
	clientConfig.TLSRequirePostQuantumSignature = true
	client := makeTestWebsocketNodeWithConfig(t, clientConfig)
	client.tlsConfig.RootCAs = roots
	client.phonebook.ReplacePeerList([]string{address}, "default", phonebook.RelayRole)
	require.NoError(t, client.Start())
	defer netStop(t, client, "PQ client")
	waitReady(t, client, time.After(5*time.Second))

	peers := client.GetPeers(PeersConnectedOut)
	require.Len(t, peers, 1)
	peer := peers[0].(*wsPeer)
	conn := peer.conn.(wsPeerWebsocketConnImpl)
	tlsConn, ok := conn.UnderlyingConn().(*tls.Conn)
	require.True(t, ok)
	require.Equal(t, pqKeyExchange, tlsConn.ConnectionState().CurveID)

	classicalConfig := defaultConfig
	classicalConfig.GossipFanout = 1
	classical := makeTestWebsocketNodeWithConfig(t, classicalConfig)
	classical.tlsConfig = &tls.Config{
		RootCAs:          roots,
		CurvePreferences: []tls.CurveID{tls.X25519},
	}
	require.NoError(t, classical.Start())
	defer netStop(t, classical, "classical client")

	gossipAddress, err := classical.addrToGossipAddr(address)
	require.NoError(t, err)
	_, ok = classical.tryConnectReserveAddr(address)
	require.True(t, ok)
	classical.wg.Add(1)
	classical.tryConnect(address, gossipAddress)
	require.Empty(t, classical.GetPeers(PeersConnectedOut))
}

func writeTestTLSFiles(t *testing.T) (string, string, *x509.CertPool) {
	t.Helper()

	certFile, keyFile := tlstest.CertificateFiles(t, "127.0.0.1", "localhost")
	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	require.NoError(t, err)

	root, err := x509.ParseCertificate(certificate.Certificate[0])
	require.NoError(t, err)
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return certFile, keyFile, roots
}

func TestWebsocketNetworkRejectsClassicalSignature(t *testing.T) {
	partitiontest.PartitionTest(t)

	certFile, keyFile, roots := writeTestECDSATLSFiles(t)
	serverConfig := defaultConfig
	serverConfig.GossipFanout = 1
	serverConfig.TLSCertFile = certFile
	serverConfig.TLSKeyFile = keyFile
	server := makeTestWebsocketNodeWithConfig(t, serverConfig)
	require.NoError(t, server.Start())
	defer netStop(t, server, "classical certificate server")

	address, listening := server.Address()
	require.True(t, listening)
	clientConfig := defaultConfig
	clientConfig.GossipFanout = 1
	clientConfig.TLSRequirePostQuantumSignature = true
	client := makeTestWebsocketNodeWithConfig(t, clientConfig)
	client.tlsConfig.RootCAs = roots
	require.NoError(t, client.Start())
	defer netStop(t, client, "PQ signature client")

	gossipAddress, err := client.addrToGossipAddr(address)
	require.NoError(t, err)
	_, ok := client.tryConnectReserveAddr(address)
	require.True(t, ok)
	client.wg.Add(1)
	client.tryConnect(address, gossipAddress)
	require.Empty(t, client.GetPeers(PeersConnectedOut))
}

func writeTestECDSATLSFiles(t *testing.T) (string, string, *x509.CertPool) {
	t.Helper()

	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()
	certificate := server.TLS.Certificates[0]
	keyDER, err := x509.MarshalPKCS8PrivateKey(certificate.PrivateKey)
	require.NoError(t, err)

	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")
	require.NoError(t, os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Certificate[0]}), 0600))
	require.NoError(t, os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}), 0600))

	root, err := x509.ParseCertificate(certificate.Certificate[0])
	require.NoError(t, err)
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return certFile, keyFile, roots
}
