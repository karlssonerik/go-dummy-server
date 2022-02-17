package webserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func Getter(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params := mux.Vars(req)
	id := params["id"]

	response := Response{
		ID: id,
	}

	marshalAndWriteJSONResponse(ctx, w, req, http.StatusOK, response)
}
