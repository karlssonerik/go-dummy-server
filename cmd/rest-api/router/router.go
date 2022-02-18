package router

import "github.com/gorilla/mux"

var router *mux.Router

func init() {
	router = mux.NewRouter()
}

func Get() *mux.Router {
	return router
}
