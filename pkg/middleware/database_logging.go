package middleware

import (
	"net/http"

	"gorm.io/gorm"

	sdkdatabasecontext "github.com/scribd/go-sdk/pkg/context/database"
	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

// DatabaseLoggingMiddleware wraps an instantiated sdk.Logger that will be injected
// in the request context.
type DatabaseLoggingMiddleware struct{}

// NewLoggingMiddleware is a wrapper of an SDK Logger. It can be used to
// build a LoggingMiddleware.
func NewDatabaseLoggingMiddleware() DatabaseLoggingMiddleware {
	return DatabaseLoggingMiddleware{}
}

// Handler implements go-sdk's middleware.Handlerer interface: it returns a
// http.Handler to be mounted as middleware.
// This handler extracts the logger from the request context, attaches it to
// gorm's database connection pool and logs the database queries with their
// meta-information using the logger.
func (dlm DatabaseLoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := sdkdatabasecontext.Extract(r.Context())
		if err != nil {
			http.Error(w, "Unable to get DB connection", http.StatusInternalServerError)
			return
		}

		l, err := sdkloggercontext.Extract(r.Context())
		if err != nil {
			http.Error(w, "Unable to get logger", http.StatusInternalServerError)
			return
		}

		newDB := db.Session(&gorm.Session{
			Logger: sdklogger.NewGormLogger(l),
			NewDB:  true,
		})

		ctx := sdkdatabasecontext.ToContext(r.Context(), newDB)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
