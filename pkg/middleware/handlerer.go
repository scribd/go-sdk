package middleware

import (
	"net/http"
)

// A Handlerer provides a http.Handler that responds to an HTTP request.
// The returned Handler can be mounted as Middleware.
type Handlerer interface {
	Handler(next http.Handler) http.Handler
}
