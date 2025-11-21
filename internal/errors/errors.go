package errors

import (
	"net/http"
	"strings"
)

var (
	RespInvalidCode = []byte(`{"error": "invalid code"}`)

	RespNotFoundErr = []byte(`{"error": "not found"}`)
	RespConflictErr = []byte(`{"error": "already exists"}`)

	RespJSONDecodeErr = []byte(`{"error": "json decode error"}`)
	RespJSONEncodeErr = []byte(`{"error": "json encode error"}`)

	RespRepoDataAccessErr = []byte(`{"error": "repo data access error"}`)
	RespRepoDataInsertErr = []byte(`{"error": "repo data insert error"}`)
)

func BadRequest(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(error)
}

func NotFound(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(error)
}

func Conflict(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusConflict)
	w.Write(error)
}

func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(error)
}

func IsDuplicateDBEntry(errStr string) bool {
	return strings.Contains(errStr, "duplicate key value violates")
}
