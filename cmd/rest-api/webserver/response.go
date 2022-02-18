package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-dummy-server/cmd/rest-api/logger"
)

var log = logger.Log()

var ErrResponseInternalServerError = []byte(`{"error": {"message": "internal server error"}}`)

func marshalAndWriteJSONResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, body interface{}) {
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Errorw("Failed to marshal response body",
			"type", fmt.Sprintf("%T", body),
			"error", err)

		code = http.StatusInternalServerError
		bytes = ErrResponseInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(bytes)
	if err != nil {
		log.Errorw("Failed to write response",
			"error", err)
	}
}
