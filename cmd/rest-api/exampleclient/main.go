package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type ExampleResponse struct {
	IsExample bool   `json:"isExample"`
	Message   string `json:"message"`
}

func main() {
	addEndpointRequest := AddEndpointRequest{
		Path:     "/example",
		Method:   http.MethodGet,
		Request:  "",
		Response: ExampleResponse{true, "This is an example"},
	}

	request, err := json.Marshal(addEndpointRequest)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:8666/addEndpoint", "application/json", bytes.NewBuffer(request))
	if err != nil {
		panic(err)
	}

	var response AddEndpointResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created ", response.Created)

	getResp, err := http.Get("http://localhost:8666/example")
	if err != nil {
		panic(err)
	}

	var exampleResponse ExampleResponse

	err = json.NewDecoder(getResp.Body).Decode(&exampleResponse)
	if err != nil {
		panic(err)
	}

	fmt.Println("Got ", exampleResponse)
}
