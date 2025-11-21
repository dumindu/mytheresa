package middleware

import (
	"net/http"

	"github.com/rs/xid"

	"github.com/mytheresa/go-hiring-challenge/internal/util/ctxutil"
)

const requestIDHeaderKey = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(requestIDHeaderKey)
		if requestID == "" {
			requestID = xid.New().String()
		}

		ctx = ctxutil.SetRequestID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
