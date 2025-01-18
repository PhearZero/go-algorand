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

	"H4sIAAAAAAAC/+x9/XfbtpLov4Kn3XPysaLtfLR7m3d69rlN2+tt0uTEbnfvxnktRI4kXFMALwDKUvPy",
	"v7+DAUCCJCBRtpu0u/0psYiPwWAwmC/MvJ/kYlUJDlyrybP3k4pKugINEv+ieS5qrjNWmL8KULlklWaC",
	"T575b0RpyfhiMp0w82tF9XIynXC6graN6T+dSPhHzSQUk2da1jCdqHwJK2oG1tvKtG5G2mQLkbkhTu0Q",
	"Z88nH3Z8oEUhQakhlK94uSWM52VdANGSckVz80mRa6aXRC+ZIq4zYZwIDkTMiV52GpM5g7JQR36R/6hB",
	"boNVusnTS/rQgphJUcIQzq/FasY4eKigAarZEKIFKWCOjZZUEzODgdU31IIooDJfkrmQe0C1QITwAq9X",
	"k2dvJwp4ARJ3Kwe2xv/OJcCvkGkqF6An76axxc01yEyzVWRpZw77ElRdakWwLa5xwdbAiel1RF7WSpMZ",
	"EMrJm2+/Jk+ePPnCLGRFtYbCEVlyVe3s4Zps98mzSUE1+M9DWqPlQkjKi6xp/+bbr3H+c7fAsa2oUhA/",
	"LKfmCzl7nlqA7xghIcY1LHAfOtRvekQORfvzDOZCwsg9sY3vdFPC+T/pruRU58tKMK4j+0LwK7Gfozws",
	"6L6LhzUAdNpXBlPSDPr2JPvi3ftH00cnH/7p7Wn2X+7Pz558GLn8r5tx92Ag2jCvpQSeb7OFBIqnZUn5",
	"EB9vHD2opajLgizpGjefrpDVu77E9LWsc03L2tAJy6U4LRdCEerIqIA5rUtN/MSk5qVhU2Y0R+2EKVJJ",
	"sWYFFFPDfa+XLF+SnCo7BLYj16wsDQ3WCooUrcVXt+MwfQhRYuC6ET5wQb9fZLTr2oMJ2CA3yPJSKMi0",
	"2HM9+RuH8oKEF0p7V6nDLitysQSCk5sP9rJF3HFD02W5JRr3tSBUEUr81TQlbE62oibXuDklu8L+bjUG",
	"aytikIab07lHzeFNoW+AjAjyZkKUQDkiz5+7Icr4nC1qCYpcL0Ev3Z0nQVWCKyBi9nfItdn2fz9/9QMR",
	"krwEpegCXtP8igDPRQHFETmbEy50QBqOlhCHpmdqHQ6u2CX/dyUMTazUoqL5VfxGL9mKRVb1km7Yql4R",
	"Xq9mIM2W+itECyJB15KnALIj7iHFFd0MJ72QNc9x/9tpO7KcoTamqpJuEWEruvnyZOrAUYSWJamAF4wv",
	"iN7wpBxn5t4PXiZFzYsRYo42expcrKqCnM0ZFKQZZQckbpp98DB+GDyt8BWA4wdJgtPMsgccDpsIzZjT",
	"bb6Qii4gIJkj8qNjbvhViyvgDaGT2RY/VRLWTNSq6ZSAEafeLYFzoSGrJMxZhMbOHToMg7FtHAdeORko",
	"F1xTxqEwzBmBFhoss0rCFEy4W98Z3uIzquDzp6k7vv06cvfnor/rO3d81G5jo8weycjVab66AxuXrDr9",
	"R+iH4dyKLTL782Aj2eLC3DZzVuJN9Hezfx4NtUIm0EGEv5sUW3CqawnPLvlD8xfJyLmmvKCyML+s7E8v",
	"61Kzc7YwP5X2pxdiwfJztkggs4E1qnBht5X9x4wXZ8d6E9UrXghxVVfhgvKO4jrbkrPnqU22Yx5KmKeN",
	"thsqHhcbr4wc2kNvmo1MAJnEXUVNwyvYSjDQ0nyO/2zmSE90Ln81/1RVaXrrah5DraFjdyWj+cCZFU6r",
	"qmQ5NUh84z6br4YJgFUkaNviGC/UZ+8DECspKpCa2UFpVWWlyGmZKU01jvTPEuaTZ5N/Om7tL8e2uzoO",
	"Jn9hep1jJyOyWjEoo1V1wBivjeijdjALw6DxE7IJy/ZQaGLcbqIhJWZYcAlryvVRq7J0+EFzgN+6mVp8",
	"W2nH4rungiURTmzDGSgrAduG9xQJUE8QrQTRigLpohSz5of7p1XVYhC/n1aVxQdKj8BQMIMNU1o9wOXT",
	"9iSF85w9PyLfhWOjKC54uTWXgxU1zN0wd7eWu8Ua25JbQzviPUVwO4U8Mlvj0WDE/LugOFQrlqI0Us9e",
	"WjGN/+rahmRmfh/V+Y9BYiFu08SFipbDnNVx8JdAubnfo5wh4ThzzxE57fe9GdmYUXYQjDprsXjXxIO/",
	"MA0rtZcSAogCanLbQ6Wk24kTEjMU9oZk8qMCSyEVXTCO0E6N+sTJil7Z/RCId0MIoBq9yNKSlSAbE6qT",
	"OR3qjwZ2lj8AtcY21kuiRlItmdKoV2NjsoQSBWfKPUGHpHIjyhix4TsW0cB8LWlladl9sWIX46jP20YW",
	"1ltevCPvxCjMAbsPNhqhujFb3ss6o5Ag1+jB8FUp8qu/UrW8gxM+82MNaR+nIUugBUiypGoZOTg92m5H",
	"G0PfpiHSLJkFUx21S8S/72yRONqeZRZU02CZDva4NBvAmECE/TYGFV9FEfBCLNQdLL8Uh/DuqvqalqWZ",
	"esize6vEgUdxsrIkpjGBFUOPgdOcrYvBKqDkG5ovjVxEclqW09ZWJqqshDWUREjCOAc5JXpJdcv9cGSv",
	"2CEjUWC4vQYSrMbZ2dDGKBtjjASyongFr4w6V5XdPs0VougKemIgigSiRjNKoGmdPfergzVwZMrN0Ah+",
	"s0Y0V4WDH5m53SecmQu7OGsC1d5/2eCvYZgdoE3rVqDg7RRCFtZor81vTJJcSDuEFXHc5OY/QGXb2R7P",
	"+5WEzA0h6RqkoqVZXW9RDxryvauT+1ud2ekkBxkxU73C/9CSmM9GjDOU1FIPQ2lMBP7kwkomBlV2JtMA",
	"Dc6CrKwtl1Q0vzoIyq/byePsZdTJ+8aaj90WukU0O3SxYYW6q23CwVJ71T0h1njn2dFAGNvJdIK5xiDg",
	"QlTEso8eCJZT4GgWIWJz5/f6V2IT5fZiM7jTxQbuZCfMOKOZ/Vdi89xBJuR+zOPYo64zsSGcrkDh9c5D",
	"xmlmaR2TpzMhbyZO9S4YTlp3K6Fm1ECanPaQhE3rKnNnM+KysQ16A7URLruloP7wMYx1sHCu6W+ABWVG",
	"vQssdAe6ayyIVcVKuAPSX0al2BlV8OQxOf/r6WePHv/8+LPPDUlWUiwkXZHZVoMi951dkii9LeFBVD1E",
	"6SI++udPvZOuO25sHCVqmcOKVsOhrPPPqv+2GTHthljrohlX3QA4iiOCudos2on1axvQnsOsXpyD1kbV",
	"fy3F/M654WCGGHTY6HUljWChuo5SJy0dF6bJMWy0pMcVtgRe2EALsw6mjBK8mt0JUaU2vmhnKYjDaAF7",
	"D8Wh29ROsw23Sm5lfRf2HZBSyOgVXEmhRS7KzMh5TEQsNK9dC+Ja+O2q+r9baMk1VcTMje7bmhcJQ4ze",
	"8PH3lx36YsNb3Oy8wex6I6tz847Zly7yWy2kApnpDSdInR370FyKFaGkwI4oa3wH2spfbAXnmq6qV/P5",
	"3Zh7BQ4UMWSxFSgzE7EtjPSjIBfcRjPusVm5Ucegp48Y72bTaQAcRs63PEdf4V0c27Q5b8U4Bi6oLc8D",
	"256BsYRi0SHL29vwUuiwU91TEXAMOl7gZ3RWPIdS02+FvGjF1++kqKs7Z8/9Occuh7rFOHdIYfp6Ozjj",
	"i7IbQbswsB/F1vhJFvR1Y0Swa0DokSJfsMVSB/riayl+gzsxOksMUPxgrWWl6TO0mf0gCsNMdK3uQJRs",
	"B2s5nKHbkK/Rmag1oYSLAnDzaxUXMhMxlxjshTFqOpRb0T7BFJmBoa6c1ma1dUUwAmtwX7QdM5rbE5oh",
	"alQi/qQJHLKt7HQ2nq+UQIstmQFwImYuyMOFn+AiKYaPaS+mORE3wi86cFVS5KAUFJmzxe8FzbezV4fe",
	"gScEHAFuZiFKkDmVtwb2ar0XzivYZhjsqMj9739SDz4BvFpoWu5BLLaJobdvTxtCPW76XQTXnzwkO2up",
	"s1RrxFvDIErQkELhQThJ7l8fosEu3h4ta5AYU/ObUryf5HYE1ID6G9P7baGtq0QIv1PTjYRnNoxTLrxg",
	"FRuspEpn+9iyadSxJZgVBJwwxolx4ITg9YIqbePAGC/QpmmvE5zHCmFmijTASTXEjPyT10CGY+fmHuSq",
	"Vo06ouqqElJDEVsDuqSTc/0Am2YuMQ/GbnQeLUitYN/IKSwF4ztkOQ0Y/6C6cUA7l/ZwcRhUYO75bRSV",
	"HSBaROwC5Ny3CrAbhjEnAGGqRbQlHKZ6lNPETk8nSouqMtxCZzVv+qXQdG5bn+of27ZD4rJODntvFwIU",
	"OlBcewf5tcWsDWBfUkUcHD7GAM05NmBtCLM5jJliPIdsF+WjimdahUdg7yGtq4WkBWQFlHQbiY6wn4n9",
	"vGsA3PFW3RUaMhuJHN/0lpJ94OeOoQWOp2LCI8EvJDdH0KgCLYG43ntGLgDHjjEnR0f3mqFwrugW+fFw",
	"2XarIyPibbgW2uy4owcE2XH0MQAn8NAMfXNUYOes1T37U/wNlJugkSMOn2QLKrWEdvyDFpCwBbtHXsF5",
	"6bH3HgeOss0kG9vDR1JHNmGYfk2lZjmrUNf5HrZ3rvr1J4g6zkkBmrISChJ8sGpgFfYnNoa2P+bNVMFR",
	"trch+APjW2Q5Pk6pC/wVbFHnfm0fZwSmjrvQZSOjmvuJcoKA+pBvI4KHTWBDc11ujaCml7Al1yCBqHpm",
	"QxiG/hQtqiwcIOqf2TGj885GfaM73cXnOFSwvFiwndUJdsN30VMMOuhwukAlRDnCQjZARhSCUbEjpBJm",
	"15l7/+VfAHlK6gDpmDa65pvr/57qoBlXQP4mapJTjipXraGRaYREQQEFSDODEcGaOV10ZoshKGEFVpPE",
	"Lw8f9hf+8KHbc6bIHK79o0nTsI+Ohw/RjvNaKN05XHdgDzXH7SxyfaDjylx8Tgvp85T9IV9u5DE7+bo3",
	"eOPtMmdKKUe4Zvm3ZgC9k7kZs/aQRsaFu+G4o3w53figwbpx38/Zqi6pvguvFaxpmYk1SMkK2MvJ3cRM",
	"8G/WtHzVdMMHoZAbGs0hy/EZ48ix4ML0sS8fzTiMM3OA7auHsQDBme11bjvtUTHbUF22WkHBqIZySyoJ",
	"OdgHf0ZyVM1Sj4h9CpAvKV+gwiBFvXDRvXYcZPi1sqYZWfPBEFGhSm94hkbu2AXgwtT8m08jTgE1Kl3f",
	"Qm4VmGvazOee+Y65mYM96HsMok6y6SSp8RqkrluN1yKn+3B1xGXQkfcC/LQTj3SlIOqM7DPEV7gt5jCZ",
	"zf1tTPbt0DEohxMHIc/tx1TUs1G3y+0dCD12ICKhkqDwigrNVMp+FfPwkboPFdwqDauhJd92/Tlx/N4k",
	"9UXBS8YhWwkO22heFsbhJX6MHie8JhOdUWBJ9e3rIB34e2B15xlDjbfFL+52/4T2PVbqWyHvyiVqBxwt",
	"3o/wQO51t7spb+onpWUZcS26J6x9BqCmTbAuk4QqJXKGMttZoaYuKth6I9171y76XzcPc+7g7PXH7fnQ",
	"wuwIaCOGsiKU5CVDC7LgSss615ecoo0qWGokiMsr42mr5de+SdxMGrFiuqEuOcUAvsZyFQ3YmEPETPMt",
	"gDdeqnqxAKV7us4c4JK7VoyTmjONc63MccnsealAYiTVkW25olsyNzShBfkVpCCzWnelf3yhrTQrS+fQ",
	"M9MQMb/kVJMSqNLkJeMXGxzOO/39keWgr4W8arAQv90XwEExlcWDzb6zX/Fhg1v+0j1ywHB3+9kHnbYp",
	"IyZmmZ0sMf/3/r89e3ua/RfNfj3JvviX43fvn3548HDw4+MPX375/7o/Pfnw5YN/++fYTnnYY++HHeRn",
	"z51mfPYc1Z8gVL8P+0ez/68Yz6JEFkZz9GiL3MdcGY6AHnSNY3oJl1xvuCGkNS1ZYXjLTcihf8MMzqI9",
	"HT2q6WxEzxjm13qgUnELLkMiTKbHGm8sRQ3jM+Mv9dEp6R7f43mZ19xupZe+7UNUH18m5tMmG4NN1PaM",
	"4FP9JfVBnu7Px599Ppm2T+yb75PpxH19F6FkVmxiiRQK2MR0xfCRxD1FKrpVoOPcA2GPhtLZ2I5w2BWs",
	"ZiDVklUfn1MozWZxDuffbDmb04afcRvgb84Puji3znMi5h8fbi0BCqj0MpbAqSOoYat2NwF6YSeVFGvg",
	"U8KO4Khv8ymMvuiC+kqgcx+YKoUYow0158ASmqeKAOvhQkYZVmL003ve4C5/defqkBs4Bld/zlhE773v",
	"vrkgx45hqns2p4cdOsjCEFGl3evRTkCS4Wbhm7JLfsmfwxytD4I/u+QF1fR4RhXL1XGtQH5FS8pzOFoI",
	"8sw/SH1ONb3kA0krmVkyeDVOqnpWspxchQpJS542W9hwhMvLt7RciMvLd4PYjKH64KaK8hc7QWYEYVHr",
	"zOU6yiRcUxnzfakm1w2ObJOZ7ZrVCtmitgZSn0vJjR/nebSqVD/nxXD5VVWa5QdkqFxGB7NlRGnRvEcz",
	"Aop702z29wfhLgZJr71dpVagyC8rWr1lXL8j2WV9cvIEX/a1SSB+cVe+ocltBaOtK8mcHH2jCi7cqpUY",
	"q55VdBFzsV1evtVAK9x9lJdXaOMoS4LdOq8O/QMDHKpdQPPGO7kBFo6DX0fj4s5tL5/XMr4E/IRb2H2B",
	"fqv9ChII3Hi79iQhoLVeZuZsR1elDIn7nWnS3S2MkOWjMRRboLbqMgPOgORLyK9cyjZYVXo77XT3AT9O",
	"0PSsgymbzM++MMR0UuigmAGpq4I6UZzybT+vj7IvKnDQN3AF2wvRZqM6JJFPN6+MSh1UpNRAujTEGh5b",
	"N0Z/811UmX9o6tKz4ONNTxbPGrrwfdIH2Yq8d3CIY0TRyXuSQgSVEURY4k+g4AYLNePdivRjy2M8B67Z",
	"GjIo2YLNYnmI/2PoD/OwGqp0qRddFHIzoCJsTowqP7MXq1PvJeULMNezuVKFoqVNKxsN2kB9aAlU6hlQ",
	"vdPOz8OMHB46VCmv8eU1WvimZgmwMfvNNFrsOFwbrQINRbaNi14+SsefWcChuCE8vnurKRwldV2HukjK",
	"RX8rN9ht1FoXmhfSGcJlv68Ac7aKa7MvBgrh0o3arDbB/VIruoCE7hJ670YmBOl4/HCQfRJJVAYR876o",
	"MZAEoiDbxplZc/QMg/liDjGqmb2ATD+TdRA7nxFmEXcIm5UowDaRq3bvqex4UW1a5BRocdYCkreioAej",
	"i5HwOC6p8scRE8Z6LjtKOvsN897sys13FsQSBllhm8x7/jbsc9CB3u8y9Pm0fD4XX6j0j8irZ3QvfL4Q",
	"2w7BUTQtoISFXbht7AmlzRjVbpCB49V8jrwli4UlBgbqQABwc4DRXB4SYn0jZPQIMTIOwMbABxyY/CDC",
	"s8kXhwDJXcYr6sfGKyL4G+IP+2ygvhFGRWUuV5bwN+aeA7hUFK1k0YuoxmEI41Ni2NyalobNOV28HWSQ",
	"Ig4Vil5COBd68yClaOxwTdkr/6A1WSHhJqsJpVkPdFzU3gHxTGwy+0I5qovMNjND79G3C/heOnYwbTK+",
	"e4rMxAbDufBqsbHye2BJw+HBCGwvG6aQXrFfSs6ywOyadrecG6NChSTjDK0NuaQEvTFTJ2TLFLncD/Lr",
	"3QiAnhmqLVbhzBJ7zQdd8WR4mbe32rTNG+ufhcWOf+oIRXcpgb+hfaybEe+vbebDdHY1f6I+SirAoWXp",
	"NikabefKpl08JENjnxw6QOzA6uu+HBhFazfWq4vXAGsxVmKY79ApOUSbghJQCc46oml2FYsUMLo84D1+",
	"7rsFxjrcPcq3D4IAQgkLpjS0TiMfF/QpzPEU80cLMU+vTldybtb3Rojm8rduc+zYWeZHXwFG4M+ZVDpD",
	"j1t0CabRtwqNSN+apnEJtBuiaKstsCLOcXHaK9hmBSvrOL26eb9/bqb9obloVD3DW4xxG6A1w+og0cDl",
	"HVPb2PadC35hF/yC3tl6x50G09RMLA25dOf4g5yLHgPbxQ4iBBgjjuGuJVG6g0EGD86H3DGQRoOYlqNd",
	"3obBYSr82Huj1Pyz99TNb0eKriVIAxh/ISgWCyh8ejPvD+NBErlS8EVQxqqqduXMOyI2dR1mntuRtM6F",
	"4UMqCD8Q9zPGC9jEoQ+1AoS8fVmHCfdwkgVwm64kbhaKoiYM8ccWga3uI/tC+w8AokHQFz1ndhudbHep",
	"2U7cgBJo4XQSBX59u4/lcEMc6qap8OlO6tfdRwgHRJpiOqjsMkxDkGDAtKpYsek5nuyoSSMYPci6nJC2",
	"kLW4wfZgoBsEHSW4Ti5xF2rtDOzHqPMeG63Mxl67wGJD3zR3D/CLWqIHoxPZPExc3+hqI9f+/U/nWki6",
	"AOeFyixItxoCl3MIGoK08IpoZsNJCjafQ+h9UTfxHHSAG9jYixGkGyGyuIumZlx//jRGRnuop4VxP8ri",
	"FBOhhZRP/mLo5fIyfWBKaq6EYGtu4KqKPtf/HrbZT7SsjZLBpGrDc53bqXv5HrDr69X3sMWR90a9GsD2",
	"7Apant4A0mDM0t98UkEG73uqU+MA1cvOFh6wU6fxXbqjrXFVKdLE394ynaoN3aXc5mC0QRIGljG7cR6P",
	"TTCnB7qI75Pyvk1gxX4ZJJD3w6mY8jU8h1dRk4tiH+1eAC098eJyJh+mk9tFAsRuMzfiHly/bi7QKJ4x",
	"0tR6hjuBPQeinFaVFGtaZi5eInX5S7F2lz829+EVH1mTiVP2xTenL1478D9MJ3kJVGaNJSC5KmxX/WFW",
	"ZetY7L5KbLZvZ+i0lqJg85uMzGGMxTVm9u4ZmwZVYdr4meAoupiLeTzgfS/vc6E+dok7Qn6gaiJ+Wp+n",
	"DfjpBvnQNWWldzZ6aBPB6bi4caWFolwhHODWwUJBzFd2p+xmcLrjp6Olrj08Ced6hakp4xoHd4krkRW5",
	"4B9659LTt0J2mL97mRgNHvrtxCojZFs8JmK1fQHPvjB1RKzg9cviF3MaHz4Mj9rDh1PyS+k+BADi7zP3",
	"O+oXDx9GvYdRM5ZhEmil4nQFD5pXFsmN+LgKOIfrcRf06XrVSJYiTYYNhdooII/ua4e9a8kcPgv3SwEl",
	"mJ+Oxijp4aZbdIfAjDlB56mXiE2Q6crWDFVE8H5MNT6CNaSFzN6VZLDO2OER4vUKHZiZKlkeD+3gM2XY",
	"K7fBlKYxwcYJa60ZsWaJ2Fxes2As02xMztQekMEcUWSqaNrWFncz4Y53zdk/aiCsMFrNnIHEe6131Xnl",
	"AEcdCKRxu5gb2Pqp2uFvYwfZ4W/ytqBdRpCd/rvnjU/JLzRW9ejACPBwxgHj3hG97ejDUbN9zbbshmCO",
	"02PG1I73jM456xJzRGvBM5XNpfgV4o4Q9B9FEmF4xydDM++vwGORe32W0jiV25L27ez7tnu8bpza+Fvr",
	"wn7RTdm1m1ym8VN92EbeROlV8XTNDskpJSyMMOg+DUiwFjxeQTAslkHx0UeU2/Nks0B0XpjFT2X4lvPY",
	"jt+eSgfz4P1rSa9nNFYjxuhCBqZgeztxUloQ39lvgGpyHNjZSRDB3bRlNpNcBbL1QQyz0t5Qr7HTjtZo",
	"WgUGKSpUXaY2TKFUIjJMza8pt2XUTT/Lr1xvBdYFb3pdC4l5IFU8pKuAnK2i5tjLy7dFPgzfKdiC2Qrh",
	"tYKgBLUbiNhkk0hFrox3k7nDoeZsTk6mQR18txsFWzPFZiVgi0e2xYwqvC4bd3jTxSwPuF4qbP54RPNl",
	"zQsJhV4qi1glSKN7opDXBCbOQF8DcHKC7R59Qe5jSKZia3hgsOiEoMmzR19gQI394yR2y7oK77tYdoE8",
	"2wdrx+kYY1LtGIZJulHj0ddzCfArpG+HHafJdh1zlrClu1D2n6UV5XQB8fcZqz0w2b64m+jO7+GFW28A",
	"KC3FljAdnx80Nfwp8ebbsD8LBsnFasX0ygXuKbEy9NTWl7aT+uGwEJmvF+Xh8h8x/rXy4X89W9dHVmPo",
	"KvFmC6OUf0AfbYjWKaE2+WfJ2sh0X7CUnPncwlhAq6mbZXFj5jJLR1kSA9XnpJKMa7R/1Hqe/cWoxZLm",
	"hv0dpcDNZp8/jRSi6tZq4YcB/tHxLkGBXMdRLxNk72UW15fc54JnK8NRigdtjoXgVCYDdeMhmam40N1D",
	"j5V8zShZktzqDrnRgFPfivD4jgFvSYrNeg6ix4NX9tEps5Zx8qC12aEf37xwUsZKyFjBgPa4O4lDgpYM",
	"1vhiLr5JZsxb7oUsR+3CbaD/tPFPXuQMxDJ/lqOKQODR3PVY3kjxP71sM5+jY9W+ROzZAIWMWDud3e4j",
	"RxseZnXr+29twBh+S2BuNNpwlCFWEtH3Nry+6fMp4oX6INk97xgcH/1CpNHBUY5/+BCBfvhw6sTgXx53",
	"P1v2/vBhPAFx1ORmfm2xcBuNGPvG9vArETGA+aqFTUCRy48QMUCmLinzwTDBmRtqSroV4j6+FHE377vi",
	"0abxU3B5+Ra/eDzgH31EfGJmiRvYvlJIH/ZuhcwoyRTN9yDOnZKvxGYs4fTuIE88vwMUJVAy0jyHKxlU",
	"AI266/fGiwQ0akadQSmMkhkWBQrt+X8cPJvFT3dgu2Zl8VOb2613kUjK82U0SnhmOv5sZfTOFWxZZbTO",
	"yJJyDmV0OKvb/ux14IiW/ncxdp4V4yPb9ivQ2uX2FtcC3gXTA+UnNOhlujQThFjtps1q0jKUC1EQnKct",
	"atEyx2Ep51gJzcj7Zhx2VWsXt4pvwV3CoTkrMQwz7jfGlpmkOpFAC+ud+/pCZhwsP66smcGODpJQtsKL",
	"WdFVVQKezDVIusCugkOvO6ZQw5GDihVEVeYTtsSEFYLoWnIi5vNgGcA1k1Bup6SiStlBTsyyYINzT549",
	"OjmJmr0QOyNWarHol/mqXcqjY2xiv7giS7YUwEHA7of1Q0tRh2zskHBcTcl/1KB0jKfiB/tyFb2k5ta2",
	"9SSb2qdH5DvMfGSIuJPqHs2VPolwN6FmXZWCFlNMbnzxzekLYme1fWwJeVvPcoHWui75R90r4xOM+sxO",
	"icw548fZncrDrFrprCk/GctNaFq0BTJZL+YG7Xghdo7Ic2tCbQr420kIpsiWKyiCapdWiUfiMP/RmuZL",
	"tE12JKA0rxxfiNWzs9ZzE7w+bKofIcM2cLtarLYU65QIvQR5zRTgi3xYQzcdYpMb1NnGfXrE7vJkzbml",
	"lKMDhNGm1tGhaPfAWUnWBxVEIesh/kDLlK3HfGhd2nPsFX+L0Sty2/P6++R6PsU2eemcCznlgrMcSyHE",
	"JGlM3TbOTTmiakTcv6gm7oRGDle0tG7zFthhMVls1zNCh7ihyz/4ajbVUof9U8PGlVxbgFaOs0Ex9ZWu",
	"nUOMcQWumpUhopBPChkJaoo+hGgCKA4kI8zKlLBwfmu+/eDs35gU44pxtHQ5tDn9zLqsSsXQM80J02Qh",
	"QLn1dF/zqLemzxFmaSxg8+7ohViw/JwtcAwbRmeWbWNGh0Od+ghSF7Fp2n5t2rrc+c3PnXAwO+lpVblJ",
	"03XQo4Kk3vAkgmNxSz6QJEBuM3442g5y2xn6jfepITRYY9QaVHgPDwijqaXdHeUbo1taisIWxL6ojCbQ",
	"ZTwCxgvGvQs1fkHk0SsBNwbPa6KfyiXVVncYxdMugJaJBxD4Qtn64G87VL9ygEEJrtHPkd7Gtgx4gnE0",
	"DVqJn/It8YfCUHcgTHxNyyZ0OlLUG6UqJ0QV+LioV+Y7xjgM4878k8kOuvY+32u6YzWOQ2+iVI7CWV0s",
	"QGe0KGKprb7CrwS/+kdisIG8bopQNa8DuznKh9TmJsoFV/Vqx1y+wS2nC+rmR6ghrN3vdxgz7cy2+G+s",
	"AlN6Z1zQ9MGvcn2EdHFYYv7hK+OY1GtoOlNskY3HBN4pt0dHO/XNCL3tf6eU7p/r/i5e4/a4XLhHMf72",
	"jbk4wsS9g/h0e7U0eXUxFlzgd5/wqMkI2eVKeJUN6oxh1ANuXmTLesD7hlHA17RMvIQPfSX2frX+g9R7",
	"+DyZvoFql55LU7KTBSVTHtlY4Z73ZehCTMUH2/Dgu/NauLXuRGjad/d9x1NnY8RaZpH00N3MidZu8KFe",
	"NFdwYGjSpGUp8tGn3g1zajql03mK1crluY7EsK1XogjpPIyGAogzLRueGwn5R90z+g0Vo+gXeR0frWOz",
	"cF+dIDggGUSaA3hqH+R5YPzUdqJw2MAg6vBIvmUlliX69/NXP0zS2xbge7iBLplu1ISc2obmzVKfGBYi",
	"snqs3hP9XSXM15geJ07nrs5s9MO3SkensSljkp9eRLsN9mwhYundh8k6Ji3uPKaCrWv3wh72cCtjW/j9",
	"OpWcxFfIwe9hJR4XPzd1BRhgzUTt4x796wNvjLG/uuRXnYo7Cc4TfdPzqf2FSe/mhascbZfprGHf/2Tj",
	"HwhwLbe/A1/nYNP75ZwieqY1DLdNSFN0dFQR0o48OqZ6VKxQkdPKvJXaXuodWhoUfhqQ1fMxgvgAHx+m",
	"k7PiIFE1VuxqYkeJXXgv2GKpsVbGX4EWIF/vqQXS1v/AI1YJxdrav6UZzCVfXuJwR2Of+RgCZmEtk+FY",
	"Pvx7DbnGgs9tWKsEOKSyiZnMu1v/rAmSNmQ1r6FcKZBd9T+GVZ73SNeDlGVB2j1bIfdofLWL0+bxgn17",
	"eU1Vmyipl61g9Jvp+RxyzEe+M0XcfyyBB+nHpt4iirDMg4xxrHlBiBn1D7f3twDtyuC2E56gstWtwUll",
	"kLiC7T1FOtQQLdnbPJ+9ScpuxIB1Pvvs7SkXjovXZKqhDMSCD8Z3SdDbsjTJbOtBwsMbzuVJ0lwcbRLE",
	"HVN6MfAGc5muByVcRfE8lUVuWK08rfk/x+LwyoWm0ibld2gfI2fDklXXLmU4JvRrvJY+eTgo/5vP3mln",
	"KdmVq9yBWLE+4msqC9/iTtKx2buJxYGeNzOz9unUMLwoUgQFXyHmpTBiRJZ6ytl9rdSE+t5TNia7TZ2F",
	"cM1BSigaZ2QpFGRa+KdWu+DYhQobeH4jJKhk4TELXDLp/Js2qz4WYKSYZJ66ePNwgUTCihroZJD7Pj3n",
	"LmR/bb/79Be+AN9e225Dr/srQftHc0wNkBhS/Zy423J/Wo2bmHkZ5yAz7/PtJ8Ln3VyImPG2qHN7QYcH",
	"ozGFj85atYOVRC2k+XCVPR0hSE9xBdtjqwT5Etp+B0OgreRkQQ9S/fY2+U4N3yoG9+JOwPu0GRwrIcos",
	"4WY8G2bv71P8FcuvALNvNo9LjOx3r3s2zCTkPnq3mjiS6+XWZ6uvKuBQPDgi5JTb53w+pKRb2LM3Ob+n",
	"d82/wVmL2hbUcObso0sefxeFpS7kLbmZH2Y3D1NgWN0tp7KD7MkNv+GpYLdrLIvRrZ97NFYrHwZ59KSS",
	"gKgsFDGZ5Nz6ir/Ggx4zHGHykSBLDoYQUOJ8zESVIhZFf5MEKWaoOKbCyRAgDXxMno4GCjd4FAEufm5P",
	"Mk732aebFHMioQ3fuGneTZfK0rJmldLo+zM3s3T53VxICGfE8FCbY7d5coYJbPE/M6YlldubZMfsoipm",
	"PUlieW8gZBMD2S6kjYMc4rAsxXWGzCprKszEVFvTTnUvY1/usO1nTvUMgohKqpygtiVLWpBcSAl52CP+",
	"0tpCtRISslJggGUs9mOujdy9wueVnJRiQUSViwJspaY4BaXmqjmnKDZBEM8WRYGlHXynb/sEdDxySnOn",
	"Wg9uhqLW3sIGfvMvTB+bM6LNp2YXndkogsRbAVAuf5rDkG08hBcJxyYc6tsS47x5zjZINyBjR35OtKxh",
	"SlyLfnV6d/CpBLJiSllQGlq6ZmWJKRvYJoh5aEKG4qhNiL1nGNC8Zhj11k3fYaXhytx5TU6TkAechwnH",
	"iF5KUS+WQWr3Bk6v8sraKcThKD+qGgMT8e2mmeIpWQmlnaZpR2qX3AZ73s8F11KUZdcoZUX0hbO0v6Sb",
	"0zzXL4S4mtH86gHqtVzoZqXF1Gc26IfltjPJXlK/7gWcIQ2o/UmybTsMUnVEO5pB9ljcwCi+z8ocgPlu",
	"Pwfdb3M/HS6sv64uM42rMaecUC1WLI+fqT9WnGsyOjXGoqLZAm1VU5vfBZvhYQ8vqyasCVnkEM3AabQs",
	"4ylxjMCFdyC7Mf9FCbw/LpmDYzSJi3LIXJwUleVJWa8HAEJqkw7oWtpSqKEk1nAVsbBJSjA4pQ/oyFsF",
	"YwBvB5sZ4c6B0nAroAZxxw2A963xYWqzOtoY5pnY+O8P2rSPNwL+w24q7zCPVHDleUta0oZX+hRRCY4Q",
	"Ty6/MxLxAhNOzMbGIzZlq0fe8AEA6QjFDgyj4hQPBWNOWQlFFqt6etbYqKaBpu0eRXarvuO9bDl5Tmtf",
	"dNSMXUtwKYusiC+7/q+KGlISTfOhJZkXsAH7oupXkMJWE50G/hcobbHRnjFAVFkJa+gEbro8SjWKmmwN",
	"vq9qOpMCoEJvZN9GFotIDO/ynuHErT0LYtrGYDdqSbGItTtF9phJokadDc/sMVFjj5KBaM2Kmnbwpw4V",
	"ObpmQHOUI6ga6AiZ1yPHTvOjHeGNH+DU94+JMh4T78bxoYNZUBx1uxjQ3gjlWqVOPY8HKIdJwhoHC85W",
	"NI5YS+It31AVveZpg+SQ5Ft1a+Q+McEDxH6zgRylGqfvQOE0noSTwuUbQmrnAIXVCkyXiLV9CZxwERR3",
	"vaaqUVXa7KX+BzsxNmLcadM3cCq3ccS331mCgxHVS2OYVCRkQ6c3N89/kpO48yAmx4vRiAL38HaH/ctT",
	"t1M7sAEW0edmP43sj+VR3S3muPiUzGo/UFmKa1utNdRDn4P3g1rq8y4gJ5az5lr28dJTl1i3b+pgwUuR",
	"Fd0SIfEfo3X+o6Ylm2+Rz1jwfTeiltSQkHO82ogAF39tJt4tXk09YN7aIvxUdt1s7JjBcFszSgC0uch9",
	"WS1BVvQKwm3AYAfLP3NtGKeqZ2i5MFd2bzuHWHCL98mRVrQINX1M0brtcAeftNv0/t/tK9RwKp9ZsSpp",
	"7mvzuuJgXT6D9bc9ceklrHY/Ux7yNU8CTU3vlmilz2tR3MBkeiDrir39SRU+6oA9qHU8qPl0q2WMtPz2",
	"qtvseOA9ail3vQtjo24GQIcVUveBHxaM/Tj4j2ZPTi1jDPi/F7wnSkSH8Npq0B8By53cNxFYrbV6JjaZ",
	"hLnaF2BizdVGnZdt1hxvYmU8l0CVjbg5e+UUzzY5MONGEbYxoY1PsxmlgDnjLbNkvKp1RI/BHMF8GyAs",
	"NPojWhMutJSUYITJNS1frUFKVqQ2zpwOW0w1LM7iHR2ub8SE0dypwwGYanU4fBndmtHDZuYCt+XfbLim",
	"0pQXVBZhc8ZJDtLc++SabtXNPUqNc2CfT4kG0kw3X0fgXULStoCUW+cUvqW/pwGQ3qHjZ4TDBuOCI84a",
	"a9rRIuGfGcLwh3DYrOgmK8UC3+8mDoTLCo0ePqsCCo5mcCufjVu3n0exX2H3NFgQwzEiLXDWMVPsPvev",
	"cCtRjfyRM73z5FsbZf9BtY27tQfTI5Uv2uB/SyzD8xh7A+/SHoXv4L2w6Z+qeNqDYBMh4R/q2sUTu4hh",
	"EC6BQmgEH19osBtpEXtpby0DGVoM1I7wflBtKDvNXXjW0JQ2MDVYpExdnoIDLW3WPu/vpQR4aApR7qx3",
	"p21CZsw4h1Rn3J2ZIKtEleVjYj5tzZzCuQkcpF0YE/QROAES627CY1RTRaqTcaxTTurQApXJclb7vF1V",
	"vkvpT5mJEhy964IQc+RleIStcQxf8jTGlGn/jVnXDNYwCUKJhLyWaCa+ptv9Bf8SudrP/3r62aPHPz/+",
	"7HNiGpCCLUC1+f57BfPauEDG+3afjxsJOFiejm+Cz/thEef9j/5RVbMp7qxZbqvaZL6DcoGH2JcjF0Ds",
	"Ke6wUNuN9grHaUP7f1/bFVvkne9YDAW//Z5JUZbxeiuNXBVxoMR2K3ChGA2kAqmY0oYRdj2gTLcR0WqJ",
	"5kHMur22eZwEz8Hbjx0VMJ0IuYotJBVQi/wMsyo4rxGBTVU6XmU9PbvW5fQ0a6FDoRGjYmZAKlE50Z7N",
	"SQwifEEkg5e1zvCJFvEgRrZhtjZaNkaILvI8Tnphqfrd3L5bRlnHOb3ZxIh44Q/lDUgz5Z9IZwy5CSdp",
	"Tfu/G/4RSYFyZ1yjWe5vwSui+sGON8eng7iHJv3HKNCG6TAi5IEAJF7bdt5JBg/FghTg0noJ0J/gHch9",
	"8eNl61je+ywEIfEd9oAXPp9t2zUvGRw4nziX9ssGKcFS3qUoobP8fS9yPettLpJgi5zRRGtQli2JoVgY",
	"PLdWXzevmBNayeCxsxRCE6OZlmXkkbS14+CZCgnHqARyTcuPzzW+ZVLpU8QHFG/ST6PCl7Ihki0q1c0y",
	"ZL6go+YOXsXe3dT8NT7M/g8wexS959xQzgk/uM3QuENLG149b7zRwMk1jmmDrB59TmauzE0lIWeq79y/",
	"9sJJ8zAUJJu7gFbY6D0vUfet8yehb0HGcx+JQ34I3FuNz95B2B7RT8xUEic3SuUx6huQRQR/MR4VlsXe",
	"c13csiTKzRIuBakTD0y4NCz4PXZ5NrWJuXRqBcN1jr6tO7iNXNTt2sZmCxtdWeXy8q2ejUnyFa+CYrpj",
	"lrE7KYdyUDGU3yC/mMWRG8PNG6OYn1IZp21W5URW/N5+1KzcG7DSqXHwYTpZ2AxGmMX/Z1e16ePepR6C",
	"RJ4vt/TbpIuxiImstTN5MFWQ8WlE4QLXLZJtHl815rVkeosVu70Bjf0czcf0XZPbw+WGaXxp7u7T4gq4",
	"j/doM4HUyt+u3wla4n1kXXzc3EKiPCLf2Nz67qB8eW/2r/DkL0+LkyeP/nX2l5PPTnJ4+tkXJyf0i6f0",
	"0RdPHsHjv3z29AQezT//Yva4ePz08ezp46eff/ZF/uTpo9nTz7/413uGDxmQLaC+qMazyX9mp+VCZKev",
	"z7ILA2yLE1qx78HsDerKc0w1hkjN8STCirJy8sz/9H/8CTvKxaod3v86cZXRJkutK/Xs+Pj6+voo7HK8",
	"wKf/mRZ1vjz282AOuo688vqsidG3cTi4o631GDe1Sf5lvr355vyCnL4+O2oJZvJscnJ0cvTIFZXntGKT",
	"Z5Mn+BOeniXu+zFmtj1WrmjFcfNW68N08K2qbEkL82nRpO8zfy2Blphgx/yxAi1Z7j9JoMXW/V9d08UC",
	"5BG+3rA/rR8fe2nk+L3LnPBh17fjMDLk+H0nwUSxp6ePfNjX5Pi9L1q9e8BOwWIXc2aQGnV5fgfapVuy",
	"todIrg70NLjRp0RhxQrzUyWZMOd1ai7fAjAuAMPbJKbu17LmuXUW2ymA439fnv4nOsxfnv4n+ZKcTN2D",
	"A4UKTWx6++K6IbSzwoI9jFNUX21Pm2wmrXN98uxtzMjkgkWrelaynFg5BQ+qocLgHDUjtnwSLYoTe0+g",
	"o6/h+oaTn2RfvHv/2V8+xKTJgWzcIClI8NHx+gpfcxiRtqKbL1Mo27gIdDPuP2qQ23YRK7qZhAAPPaiR",
	"rGf+gZAvvR7GJgZRi/9+/uoHIiRx2vNrml81j6P8a7j2BWD4GM70TEHsLtYQaOD1ytxR7pXVSi2qburt",
	"Bs3vsE4pAors5PHJieehTkMJDuixO/fBTD2z1pDQMEwnMFQOn8IrAhua63JLqAriJDBq0dcU7j1hE1XW",
	"CaTfaRodzui2JPoK4dDX+JHaEELTcg98F736qx10uJCfylyy+5+/D5ARhSCapjLcWk8jf+7uf4/dHUol",
	"pBLmTDOMy26vHH+ddYB0smi59eAmEo0ckb+JGmVHoxXUGhoWKCSys+bCtD4RN6fLixQE0rVPh/DLw4f9",
	"hT982Ib9zeEamSzl2LCPjocPj8xOPT2Qle20U3cSeI86O4cMN9isl3TTRE1TwgXPOCyoZmsggcL59OTR",
	"H3aFZ9zGqRth2Qr1H6aTz/7AW3bGjWBDS4It7Wqe/GFXcw5yzXIgF7CqhKSSlVvyI28eAlilB+WTIfv7",
	"kV9xcc09Ioy+Wq9WVG6dEE0bnlPzoOLWTv4zyHDUCtrIRelCYSwMiqhWpvVZEPli8u6D1wFG6h67mh3P",
	"sPbs2KYQKixp7QQ9E+r4PdrWk78fOwdp4qPVm1Of0QVi2xz71IyJljYJV/xjR2l6rzdmnbuHM22C8XKq",
	"82VdHb/H/6CaHCzYVtM41ht+jCGjx+87eHKfB3jq/t52D1tgongPnJjPFap5uz4fv7f/BhPBpgLJzG2F",
	"eTTdrzbf8TGWeN8Of97yPPrjcB2dXK+Jn4+9lSamcXdbvu/82SU5tax1Ia6DWdC/YZ1zQ8jMx1r1/z6+",
	"pkwbGcqlGKVzDXLYWQMtj10lr96vbfGMwResCBL82JO6KmFzDHUV3jf0+qLzVFTaXBpfCbRjpPjxJpsx",
	"jkwqZKKt1dJ+HGpQA9Z5sQQbnusdvxERVQsyk4IWOVXa/NGm+u+qzh9uqZ71U3+cRdx6CCZaI4bZKg27",
	"Odrr68Fxx8igwb6Qs+d+wvZ92m8utw0g+ooWxCelyshLWpoNh4KcOu2gg43fWub69ELSJ5ZqPpoY8pU/",
	"fIpQzNDX0R9lPKdOUJxyjMxhlEzDABbAM8eCspkotq5+4ETSa72xKTz6zO2Ydm+Mrp2SSrpSqY93YMT8",
	"fVsu9xks/7QT/mkn/NOS9Ked8M/d/dNOONJO+KcV7U8r2v9IK9ohprOYmOnMP2lpk62B24D2nt5H2/IV",
	"DYvvJhdjupHJOq9MsVIG00eEXGBmGGpuCViDpCXJqbLSlctitMLgT0xRBsWzS551ILEhlmbi++1/bWzr",
	"ZX1y8gTIyYN+H6VZWYa8edgX5V38ZJ+ffEkuJ5eTwUgSVmINhX0rG6ZPt732Dvu/mnFfDeou4CN5TL3j",
	"M5kRVc/nLGcW5aXgC0IXoo3LxnytXOAXkAY4W72KMD1171iYezxtd6WX5b0ruQ8lgLN2C/dGHPTIJR5s",
	"YAjvwEiDfxkTZvA/Wkq/abKr2zLSnWMPuOqfXOVjcJVPzlf+6D7cwLT431LMfHry9A+7oNAQ/YPQ5Ft8",
	"c3A7cczlEc2jRbxuKmj5PDLe3NfGLYdxwHiLNhHAb9+Zi0CBXPsLtg1rfXZ8jInFlkLp44m5/rohr+HH",
	"dw3M7/3tVEm2xvrsaN0Uki0Yp2Xm4kKzNnT18dHJ5MP/DwAA///FDP0mkyYBAA==",
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
