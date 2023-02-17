package middleware

import (
	"net/http"

	sdkdatabasecontext "github.com/scribd/go-sdk/pkg/context/database"
	sdkinstrumentation "github.com/scribd/go-sdk/pkg/instrumentation"

	"gorm.io/gorm"
)

// DatabaseMiddleware wraps an instantiated *gorm.DB that will be
// injected in the request context.
type DatabaseMiddleware struct {
	Database *gorm.DB
}

// NewDatabaseMiddleware is a constructor used to build a DatabaseMiddleware.
func NewDatabaseMiddleware(d *gorm.DB) DatabaseMiddleware {
	return DatabaseMiddleware{
		Database: d,
	}
}

// Handler implements the middlewares.Handlerer interface: it returns a
// http.Handler to be mounted as middleware. The handler injects the database
// connection pool to the request context.
func (dm DatabaseMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := sdkinstrumentation.TraceDatabase(r.Context(), dm.Database)
		ctx := sdkdatabasecontext.ToContext(r.Context(), db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
