package server

import (
	"net/http"
	"rtk/api-mocker/pkg/ctxkeys"
)

type ctxKeyAuth struct{}

func AuthForwardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get("Authorization")

		ctx := ctxkeys.WithAuthKey(r.Context(), headerValue)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
