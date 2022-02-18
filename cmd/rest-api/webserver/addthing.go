package webserver

import (
	"encoding/json"
	"fmt"
	"go-dummy-server/cmd/rest-api/router"
	"net/http"
)

type AddEndpointResponse struct {
	Created string `json:"created"`
}

type AddEndpointRequest struct {
	Path     string      `json:"path"`
	Request  string      `json:"request"`
	Response interface{} `json:"response"`
	Method   string      `json:"method"`
}

var endpoints = make(map[EndpointKey]*Endpoint)

func AddEndpoint(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	request := AddEndpointRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	newEndpoint := NewEndpoint(fmt.Sprintf("%s :: %s", request.Method, request.Path), router.Get()).
		Path(request.Path).
		Methods(request.Method).
		AddResponse(request.Method, request.Request, request.Response).
		HandlerFunc(DoMatchNReturn)

	eKey := EndpointKey{
		Method: request.Method,
		Path:   request.Path,
	}

	endpoints[eKey] = newEndpoint

	marshalAndWriteJSONResponse(ctx, w, req, http.StatusOK, AddEndpointResponse{Created: fmt.Sprintf("%s :: %s", request.Method, request.Path)})
}

func DoMatchNReturn(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	eKey := EndpointKey{
		Method: req.Method,
		Path:   req.RequestURI,
	}

	endpoint := endpoints[eKey]

	var request string
	if req.Method == "POST" || req.Method == "PUT" {
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			panic(err)
		}
	}

	rKey := RequestKey{
		Method:  req.Method,
		Request: request,
	}

	marshalAndWriteJSONResponse(ctx, w, req, http.StatusOK, endpoint.responses[rKey])
}
