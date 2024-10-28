// Package experimental provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package experimental

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
	// Get a list of assets held by an account, inclusive of asset params.
	// (GET /v2/accounts/{address}/assets)
	AccountAssetsInformation(ctx echo.Context, address string, params AccountAssetsInformationParams) error
	// Returns OK if experimental API is enabled.
	// (GET /v2/experimental)
	ExperimentalCheck(ctx echo.Context) error
	// Fast track for broadcasting a raw transaction or transaction group to the network through the tx handler without performing most of the checks and reporting detailed errors. Should be only used for development and performance testing.
	// (POST /v2/transactions/async)
	RawTransactionAsync(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// AccountAssetsInformation converts echo context to params.
func (w *ServerInterfaceWrapper) AccountAssetsInformation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameterWithLocation("simple", false, "address", runtime.ParamLocationPath, ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params AccountAssetsInformationParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "next" -------------

	err = runtime.BindQueryParameter("form", true, false, "next", ctx.QueryParams(), &params.Next)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter next: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AccountAssetsInformation(ctx, address, params)
	return err
}

// ExperimentalCheck converts echo context to params.
func (w *ServerInterfaceWrapper) ExperimentalCheck(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ExperimentalCheck(ctx)
	return err
}

// RawTransactionAsync converts echo context to params.
func (w *ServerInterfaceWrapper) RawTransactionAsync(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RawTransactionAsync(ctx)
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

	router.GET(baseURL+"/v2/accounts/:address/assets", wrapper.AccountAssetsInformation, m...)
	router.GET(baseURL+"/v2/experimental", wrapper.ExperimentalCheck, m...)
	router.POST(baseURL+"/v2/transactions/async", wrapper.RawTransactionAsync, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"FfvsEw0BAA==",
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
