package sync

import (
	"github.com/algorand/go-algorand/daemon/algod/api/server/lib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"net/http"
	"runtime"
	"time"
)

// Container for the echo.HandlerFunc which is composed of the lib.ReqContext and echo.Context
type Container func(ctx lib.ReqContext, c echo.Context) error

// Route type description
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc Container
}

type Routes []Route

var routes = Routes{
	//Route{
	//	Name:        "block_doc",
	//	Method:      "GET",
	//	Path:        "/blocks/:roundNumber",
	//	HandlerFunc: GetBlock,
	//},
	Route{
		Name:        "document",
		Method:      "GET",
		Path:        "/:db/:docId",
		HandlerFunc: GetDocument,
	},
	Route{
		Name:        "_all_docs",
		Method:      "GET",
		Path:        "/:db/_all_docs",
		HandlerFunc: GetAllDocs,
	},
	Route{
		Name:        "_all_docs",
		Method:      "POST",
		Path:        "/:db/_all_docs",
		HandlerFunc: PostAllDocs,
	},
	Route{
		Name:   "_local",
		Method: "GET",
		Path:   "/:db/_local",
		HandlerFunc: func(ctx lib.ReqContext, c echo.Context) error {
			return c.JSON(http.StatusBadRequest, DatabaseError{
				Error:  "bad_request",
				Reason: "Invalid _local document id.",
			})
		},
	},
	Route{
		Name:   "_local",
		Method: "GET",
		Path:   "/:db/_local/:id",
		HandlerFunc: func(ctx lib.ReqContext, c echo.Context) error {
			//println("running not found")
			return c.JSON(http.StatusNotFound, DatabaseError{
				Error:  "not_found",
				Reason: "missing",
			})
		},
	},
	Route{
		Name:        "changes",
		Method:      "GET",
		Path:        "/:db/_changes",
		HandlerFunc: GetDatabaseChanges,
	},
	Route{
		Name:        "database",
		Method:      "GET",
		Path:        "/:db",
		HandlerFunc: GetDatabaseInfo,
	},
	Route{
		Name:        "database",
		Method:      "HEAD",
		Path:        "/:db",
		HandlerFunc: GetDatabaseInfo,
	},
	Route{
		Name:        "root",
		Method:      "HEAD",
		Path:        "",
		HandlerFunc: GetNodeInfo,
	},

	Route{
		Name:        "root",
		Method:      "GET",
		Path:        "",
		HandlerFunc: GetNodeInfo,
	},
}

// injectCtx adds the request context to the handler
func injectCtx(ctx lib.ReqContext, handler Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Use default Echo error pattern
		return handler(ctx, c)
	}
}

// RegisterHandlers holds all the routes for sync
func RegisterHandlers(e *echo.Echo, ctx lib.ReqContext) {
	// Inject context into the header middleware for the routes
	defaultHeadersMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timestamp := time.Now()
			c.Response().Header().Set("Cache-Control", "must-revalidate")
			c.Response().Header().Set("Date", timestamp.Format(http.TimeFormat))
			c.Response().Header().Set("Server", "algod/"+ctx.Node.GenesisID()+" ("+runtime.Version()+")")
			c.Response().Header().Set("X-Algod-Request-ID", random.String(10))
			c.Response().Header().Del("Vary")

			c.Response().Before(func() {
				c.Response().Header().Set("X-Algod-Body-Time", time.Since(timestamp).String())
			})

			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}

	for _, route := range routes {
		r := e.Add(route.Method, route.Path, injectCtx(ctx, route.HandlerFunc), defaultHeadersMiddleware)
		r.Name = route.Name
	}
}
