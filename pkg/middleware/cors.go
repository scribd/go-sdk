package middleware

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/scribd/go-sdk/pkg/server"
)

type (
	// CorsMiddleware wraps *cors.Cors that will act as a middleware implementation
	CorsMiddleware struct {
		cors *cors.Cors
	}
)

// NewCorsMiddleware creates cors.Cors middleware and attaches it to the CorsMiddleware wrapper
func NewCorsMiddleware(setting server.CorsSetting) *CorsMiddleware {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:         setting.AllowedOrigins,
		AllowOriginFunc:        setting.AllowOriginFunc,
		AllowOriginRequestFunc: setting.AllowOriginRequestFunc,

		AllowedHeaders:   setting.AllowedHeaders,
		AllowedMethods:   setting.AllowedMethods,
		ExposedHeaders:   setting.ExposedHeaders,
		AllowCredentials: setting.AllowCredentials,

		MaxAge: setting.MaxAge,

		OptionsPassthrough: setting.OptionsPassthrough,
	})

	return &CorsMiddleware{
		cors: corsHandler,
	}
}

// Handler implements the middlewares.Handlerer interface: it returns a
// http.Handler to be mounted as middleware. The handler calls cors.Cors.Handler
func (cm *CorsMiddleware) Handler(next http.Handler) http.Handler {
	return cm.cors.Handler(next)
}
