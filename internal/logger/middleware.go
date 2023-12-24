package logger

import (
	"fmt"
	"net/http"
	"time"
)

func MiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Main.Info(fmt.Sprintf("Request: %s %s %s\n", r.Method, r.RequestURI, time.Now().Format(time.RFC822)))
		h.ServeHTTP(w, r)
		Main.Info(fmt.Sprintf("Response: %s %s %s\n", r.Method, r.RequestURI, time.Now().Format(time.RFC822)))
	})
}
