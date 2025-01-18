// Package private provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package private

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
	// Gets the merged config file.
	// (GET /debug/settings/config)
	GetConfig(ctx echo.Context) error

	// (GET /debug/settings/pprof)
	GetDebugSettingsProf(ctx echo.Context) error

	// (PUT /debug/settings/pprof)
	PutDebugSettingsProf(ctx echo.Context) error
	// Aborts a catchpoint catchup.
	// (DELETE /v2/catchup/{catchpoint})
	AbortCatchup(ctx echo.Context, catchpoint string) error
	// Starts a catchpoint catchup.
	// (POST /v2/catchup/{catchpoint})
	StartCatchup(ctx echo.Context, catchpoint string, params StartCatchupParams) error

	// (POST /v2/shutdown)
	ShutdownNode(ctx echo.Context, params ShutdownNodeParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetConfig converts echo context to params.
func (w *ServerInterfaceWrapper) GetConfig(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetConfig(ctx)
	return err
}

// GetDebugSettingsProf converts echo context to params.
func (w *ServerInterfaceWrapper) GetDebugSettingsProf(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDebugSettingsProf(ctx)
	return err
}

// PutDebugSettingsProf converts echo context to params.
func (w *ServerInterfaceWrapper) PutDebugSettingsProf(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutDebugSettingsProf(ctx)
	return err
}

// AbortCatchup converts echo context to params.
func (w *ServerInterfaceWrapper) AbortCatchup(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "catchpoint" -------------
	var catchpoint string

	err = runtime.BindStyledParameterWithLocation("simple", false, "catchpoint", runtime.ParamLocationPath, ctx.Param("catchpoint"), &catchpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter catchpoint: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AbortCatchup(ctx, catchpoint)
	return err
}

// StartCatchup converts echo context to params.
func (w *ServerInterfaceWrapper) StartCatchup(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "catchpoint" -------------
	var catchpoint string

	err = runtime.BindStyledParameterWithLocation("simple", false, "catchpoint", runtime.ParamLocationPath, ctx.Param("catchpoint"), &catchpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter catchpoint: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params StartCatchupParams
	// ------------- Optional query parameter "min" -------------

	err = runtime.BindQueryParameter("form", true, false, "min", ctx.QueryParams(), &params.Min)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter min: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.StartCatchup(ctx, catchpoint, params)
	return err
}

// ShutdownNode converts echo context to params.
func (w *ServerInterfaceWrapper) ShutdownNode(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ShutdownNodeParams
	// ------------- Optional query parameter "timeout" -------------

	err = runtime.BindQueryParameter("form", true, false, "timeout", ctx.QueryParams(), &params.Timeout)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter timeout: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ShutdownNode(ctx, params)
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

	router.GET(baseURL+"/debug/settings/config", wrapper.GetConfig, m...)
	router.GET(baseURL+"/debug/settings/pprof", wrapper.GetDebugSettingsProf, m...)
	router.PUT(baseURL+"/debug/settings/pprof", wrapper.PutDebugSettingsProf, m...)
	router.DELETE(baseURL+"/v2/catchup/:catchpoint", wrapper.AbortCatchup, m...)
	router.POST(baseURL+"/v2/catchup/:catchpoint", wrapper.StartCatchup, m...)
	router.POST(baseURL+"/v2/shutdown", wrapper.ShutdownNode, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
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
