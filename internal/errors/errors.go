package errors

import (
	"net/http"
)

var (
	RespJSONEncodeErr     = []byte(`{"error": "json encode error"}`)
	RespRepoDataAccessErr = []byte(`{"error": "repo data access error"}`)
)

func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(error)
}
