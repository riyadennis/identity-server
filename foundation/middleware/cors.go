package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// CORS is a net/http compliant Cross Origin Resource Sharing middleware.
func CORS(next httprouter.Handle, allowedOrigins []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := cors.New(cors.Options{
			AllowedOrigins: allowedOrigins,
			AllowedHeaders: []string{
				"Authorization",
				"Access-Control-Allow-Headers",
				"Origin",
				"Accept",
				"X-Request-ID",
				"X-Requested-With",
				"Content-Type",
				"Access-Control-Request-Method",
				"Access-Control-Request-Headers",
			},
			AllowedMethods: []string{
				http.MethodPost,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodGet,
				http.MethodOptions,
				http.MethodPut,
				http.MethodHead,
			},
			AllowCredentials: true,
		})

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, p)
		})

		c.Handler(handler).ServeHTTP(w, r)
	}
}
