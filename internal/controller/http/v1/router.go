package v1

import (
	"log"
	"net/http"

	"github.com/AnNosov/communications_info/internal/usecase"

	"github.com/gorilla/mux"
)

func NewRouter(r *mux.Router, c usecase.CommunicationAction) {
	r.Use(loggingMiddleware)
	r.Use(contentTypeApplicationTextMiddleware)
	r.Use(accessControlAllowOriginMiddleware)
	r.Use(accessControlAllowHeadersMiddleware)
	r.Use(charsetMiddleware)

	NewCommunicationActionRoutes(r, c)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func contentTypeApplicationTextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		next.ServeHTTP(w, r)
	})
}

func accessControlAllowOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func accessControlAllowHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func charsetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("charset", "utf-8")

		next.ServeHTTP(w, r)
	})
}
