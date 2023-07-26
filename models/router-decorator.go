package models

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RouterDec struct {
	Router *mux.Router
}

// Cannot import glb here, so unfortunately have to define the origins again
var allowedOrigins = []string{"http://localhost:5000", "http://192.168.1.12:5000"}

func checkOrigin(origin string) string {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == origin {
			return origin
		}
	}
	return ""
}

func (routerDec *RouterDec) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	origin = checkOrigin(origin)

	if origin != "" {
		// General backend access
		w.Header().Set(
			"Access-Control-Allow-Origin", origin,
		)
		// Ability to set cookies
		w.Header().Set(
			"Access-Control-Allow-Credentials", "true",
		)
		// Allowed methods
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET",
		)
		// Special allowed fields in headers
		w.Header().Add(
			"Access-Control-Allow-Headers",
			"User-ID",
			// "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With",
		)
	}

	if r.Method == "OPTIONS" {
		return
	}

	routerDec.Router.ServeHTTP(w, r)
}
