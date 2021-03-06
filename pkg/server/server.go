package server

import (
	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/endpoint"

	netHTTP "net/http"
	httpserver "phonebook/pkg/http"

	"phonebook/pkg/middleware"

	"github.com/go-kit/kit/transport/http"
)

type Logger struct {
	Logger    log.Logger
	Namespace string
	Subsystem string
	Action    string
	Domain    string
}

//HTTPOption server info
type HTTPOption struct {
	DecodeModel interface{}
	Logger      *Logger
}

//NewHTTPServer create go kit HTTP server
func NewHTTPServer(e endpoint.Endpoint, httpOpt HTTPOption, options ...http.ServerOption) netHTTP.Handler {
	middlewares := middleware.Nop()

	if httpOpt.DecodeModel != nil {
		mval := middleware.Validator()
		middlewares = endpoint.Chain(mval)
	}

	if httpOpt.Logger != nil {
		mlog := middleware.LogAndInstrumentation(
			httpOpt.Logger.Logger,
			httpOpt.Logger.Namespace,
			httpOpt.Logger.Subsystem,
			httpOpt.Logger.Action,
			httpOpt.Logger.Domain,
		)
		middlewares = endpoint.Chain(mlog, middlewares)
	}

	e = middlewares(e)

	return http.NewServer(e, httpserver.Decode(httpOpt.DecodeModel), httpserver.Encode(), options...)
}

//NewHTTPJSONServer create go kit HTTP server
func NewHTTPJSONServer(e endpoint.Endpoint, httpOpt HTTPOption, options ...http.ServerOption) *http.Server {
	middlewares := middleware.Nop()

	if httpOpt.DecodeModel != nil {
		mval := middleware.Validator()
		middlewares = endpoint.Chain(mval)
	}

	if httpOpt.Logger != nil {
		mlog := middleware.LogAndInstrumentation(
			httpOpt.Logger.Logger,
			httpOpt.Logger.Namespace,
			httpOpt.Logger.Subsystem,
			httpOpt.Logger.Action,
			httpOpt.Logger.Domain,
		)
		middlewares = endpoint.Chain(mlog, middlewares)
	}
	return http.NewServer(e, httpserver.DecodeJSON(httpOpt.DecodeModel), httpserver.Encode(), options...)
}
