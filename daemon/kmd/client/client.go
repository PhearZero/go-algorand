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

package client

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

const (
	timeoutSecs = 120
)

// KMDClient is the client used to interact with the kmd API.
type KMDClient struct {
	httpClient http.Client
	apiToken   string
	address    string
}

func makeHTTPClient() http.Client {
	client := http.Client{
		Timeout: timeoutSecs * time.Second,
	}
	return client
}

func makeHTTPClientWithTLS(tlsConfig *tls.Config) http.Client {
	client := makeHTTPClient()
	if tlsConfig != nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = tlsConfig
		client.Transport = transport
	}
	return client
}

// MakeKMDClient instantiates a KMDClient for the given sockFile and apiToken
func MakeKMDClient(address string, apiToken string) (KMDClient, error) {
	return MakeKMDClientWithTLS(address, apiToken, nil)
}

// MakeKMDClientWithTLS instantiates a KMDClient with the supplied TLS policy.
func MakeKMDClientWithTLS(address string, apiToken string, tlsConfig *tls.Config) (KMDClient, error) {
	kcl := KMDClient{
		httpClient: makeHTTPClientWithTLS(tlsConfig),
		apiToken:   apiToken,
		address:    strings.TrimRight(address, "/"),
	}
	return kcl, nil
}

func (kcl KMDClient) endpointURL(path string) string {
	address := kcl.address
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	return address + "/" + path
}
