// Package public provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package public

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	. "github.com/algorand/go-algorand/daemon/algod/api/server/v2/generated/model"
	"github.com/algorand/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get a list of unconfirmed transactions currently in the transaction pool by address.
	// (GET /v2/accounts/{address}/transactions/pending)
	GetPendingTransactionsByAddress(ctx echo.Context, address string, params GetPendingTransactionsByAddressParams) error
	// Broadcasts a raw transaction or transaction group to the network.
	// (POST /v2/transactions)
	RawTransaction(ctx echo.Context) error
	// Get a list of unconfirmed transactions currently in the transaction pool.
	// (GET /v2/transactions/pending)
	GetPendingTransactions(ctx echo.Context, params GetPendingTransactionsParams) error
	// Get a specific pending transaction.
	// (GET /v2/transactions/pending/{txid})
	PendingTransactionInformation(ctx echo.Context, txid string, params PendingTransactionInformationParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetPendingTransactionsByAddress converts echo context to params.
func (w *ServerInterfaceWrapper) GetPendingTransactionsByAddress(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameterWithLocation("simple", false, "address", runtime.ParamLocationPath, ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPendingTransactionsByAddressParams
	// ------------- Optional query parameter "max" -------------

	err = runtime.BindQueryParameter("form", true, false, "max", ctx.QueryParams(), &params.Max)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter max: %s", err))
	}

	// ------------- Optional query parameter "format" -------------

	err = runtime.BindQueryParameter("form", true, false, "format", ctx.QueryParams(), &params.Format)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter format: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPendingTransactionsByAddress(ctx, address, params)
	return err
}

// RawTransaction converts echo context to params.
func (w *ServerInterfaceWrapper) RawTransaction(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RawTransaction(ctx)
	return err
}

// GetPendingTransactions converts echo context to params.
func (w *ServerInterfaceWrapper) GetPendingTransactions(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPendingTransactionsParams
	// ------------- Optional query parameter "max" -------------

	err = runtime.BindQueryParameter("form", true, false, "max", ctx.QueryParams(), &params.Max)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter max: %s", err))
	}

	// ------------- Optional query parameter "format" -------------

	err = runtime.BindQueryParameter("form", true, false, "format", ctx.QueryParams(), &params.Format)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter format: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPendingTransactions(ctx, params)
	return err
}

// PendingTransactionInformation converts echo context to params.
func (w *ServerInterfaceWrapper) PendingTransactionInformation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "txid" -------------
	var txid string

	err = runtime.BindStyledParameterWithLocation("simple", false, "txid", runtime.ParamLocationPath, ctx.Param("txid"), &txid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter txid: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params PendingTransactionInformationParams
	// ------------- Optional query parameter "format" -------------

	err = runtime.BindQueryParameter("form", true, false, "format", ctx.QueryParams(), &params.Format)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter format: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PendingTransactionInformation(ctx, txid, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface, m ...echo.MiddlewareFunc) {
	RegisterHandlersWithBaseURL(router, si, "", m...)
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string, m ...echo.MiddlewareFunc) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v2/accounts/:address/transactions/pending", wrapper.GetPendingTransactionsByAddress, m...)
	router.POST(baseURL+"/v2/transactions", wrapper.RawTransaction, m...)
	router.GET(baseURL+"/v2/transactions/pending", wrapper.GetPendingTransactions, m...)
	router.GET(baseURL+"/v2/transactions/pending/:txid", wrapper.PendingTransactionInformation, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"6uOjk8mH/x8AAP//cWD98yQmAQA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
