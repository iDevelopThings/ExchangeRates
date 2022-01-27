package Middleware

import (
	"net/http"

	"github.com/naoina/denco"
)

func CORS(next denco.HandlerFunc) denco.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params denco.Params) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r, params)
	}
}
