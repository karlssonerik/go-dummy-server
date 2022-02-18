package webserver

import (
	"net/http"

	"github.com/SKF/go-utility/v2/log"
	"github.com/gorilla/mux"
)

type RequestKey struct {
	Request string `json:"request"`
	Method  string `json:"method"`
}

type EndpointKey struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type Endpoint struct {
	Name      string
	PathTpl   string
	router    *mux.Router
	responses map[RequestKey]interface{}
	*mux.Route
}

func NewEndpoint(name string, router *mux.Router) *Endpoint {
	return &Endpoint{
		Name:      name,
		Route:     router.NewRoute().Name(name),
		router:    router,
		responses: make(map[RequestKey]interface{}),
	}
}

func (e *Endpoint) Path(tpl string) *Endpoint {
	e.PathTpl = tpl
	e.Route.Path(tpl)

	err := e.Route.GetError()
	if err != nil {
		log.Errorf("Failed to create path", err)
	}

	return e
}

func (e *Endpoint) Methods(methods ...string) *Endpoint {
	e.Route.Methods(methods...)
	return e
}

func (e *Endpoint) AddResponse(method, request string, response interface{}) *Endpoint {
	requestKey := RequestKey{
		Method:  method,
		Request: request,
	}

	e.responses[requestKey] = response
	return e
}

func (e *Endpoint) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Endpoint {
	e.Route.HandlerFunc(f)
	return e
}
