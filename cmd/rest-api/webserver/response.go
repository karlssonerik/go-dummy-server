package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	http_model "github.com/SKF/go-utility/v2/http-model"
	http_server "github.com/SKF/go-utility/v2/http-server"
	"github.com/SKF/go-utility/v2/log"
)

func marshalAndWriteJSONResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, body interface{}) {
	bytes, err := json.Marshal(body)
	if err != nil {
		log.WithError(err).
			WithTracing(ctx).
			WithField("type", fmt.Sprintf("%T", body)).
			Error("Failed to marshal response body")

		code = http.StatusInternalServerError
		bytes = http_model.ErrResponseInternalServerError
	}

	http_server.WriteJSONResponse(ctx, w, r, code, bytes)
}
