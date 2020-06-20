package middleware

import (
	"net/http"
)

//EnableCors enables the cors header on the responses
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
