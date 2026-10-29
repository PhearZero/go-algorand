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
	"strconv"

	"github.com/algorand/go-algorand/util/tlsutil"
)

const pqKeyExchange = tlsutil.PQKeyExchange

func makeTLSConfig(requirePQ, requireSig bool) *tls.Config {
	return tlsutil.MakeServerConfig(requirePQ, requireSig)
}

func makeClientTLSConfig(requirePQ, requireSig bool) *tls.Config {
	cfg := tlsutil.MakeClientConfig(requirePQ, requireSig)
	if cfg != nil && requireSig {
		cfg.VerifyConnection = checkTLSPeerCertificate
	}
	return cfg
}

func recordTLSKeyExchange(state *tls.ConnectionState) {
	if state == nil {
		return
	}
	networkTLSKeyExchange.Add(strconv.Itoa(int(state.CurveID)), 1)
}

func recordTLSConnection(state *tls.ConnectionState) {
	if state == nil {
		return
	}
	recordTLSKeyExchange(state)
	if len(state.PeerCertificates) != 0 {
		networkTLSPeerCertificate.Add(state.PeerCertificates[0].PublicKeyAlgorithm.String(), 1)
	}
}

func checkTLSKeyExchange(state *tls.ConnectionState, requirePQ bool) error {
	err := tlsutil.CheckKeyExchange(state, requirePQ)
	if err != nil && requirePQ {
		networkTLSKeyExchangeRejected.Inc(nil)
	}
	return err
}

func checkTLSPeerCertificate(state tls.ConnectionState) error {
	return tlsutil.CheckPeerCertificate(state)
}

func checkTLSServerCertificate(certFile, keyFile string) error {
	return tlsutil.CheckServerCertificate(certFile, keyFile)
}
