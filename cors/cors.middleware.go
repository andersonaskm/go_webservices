package cors

import "net/http"

/*
Middleware para permitir CORS na api
*/
func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length")

		handler.ServeHTTP(rw, r)
	})
}
