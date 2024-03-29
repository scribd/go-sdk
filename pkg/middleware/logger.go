package middleware

import (
	"fmt"
	"net/http"
	"time"

	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
	sdkrequestidcontext "github.com/scribd/go-sdk/pkg/context/requestid"
	sdkinstrumentation "github.com/scribd/go-sdk/pkg/instrumentation"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

const (
	// ForwardedForHeader is the 'X-Forwarded-For' (XFF) HTTP header
	// field. It is a common method for identifying the originating
	// IP address of a client connecting to a web server through an
	// HTTP proxy or load balancer.
	ForwardedForHeader = "X-Forwarded-For"
)

// LoggingMiddleware wraps an instantiated sdk.Logger that will be injected
// in the request context.
type LoggingMiddleware struct {
	logger sdklogger.Logger
}

// NewLoggingMiddleware is a wrapper of an SDK Logger. It can be used to
// build a LoggingMiddleware.
func NewLoggingMiddleware(l sdklogger.Logger) LoggingMiddleware {
	return LoggingMiddleware{
		logger: l,
	}
}

// Handler implements the middlewares.Handlerer interface: it returns a
// http.Handler to be mounted as middleware.
// This handler logs every HTTP requests that it receives with its internal
// logger extracting various details from the request itself and calculates
// the total elapsed time per request in milliseconds.
func (lm LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logContext := sdkinstrumentation.TraceLogs(r.Context())

		requestID, err := sdkrequestidcontext.Extract(r.Context())
		if err != nil {
			lm.logger.WithFields(sdklogger.Fields{
				"error": err.Error(),
			}).Tracef("Could not retrieve request id from the context")
		}

		logger := lm.logger.WithFields(sdklogger.Fields{
			"http": sdklogger.Fields{
				"request_id": requestID,
			},
			"dd": sdklogger.Fields{
				"trace_id": logContext.TraceID,
				"span_id":  logContext.SpanID,
			},
		})

		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		ctx := sdkloggercontext.ToContext(r.Context(), logger)

		// Parse the request params/form to populate r.Form for
		// logging. The request form has to be parsed before the
		// request is served.
		if err := r.ParseForm(); err != nil {
			logger.WithFields(sdklogger.Fields{
				"error": err.Error(),
			}).Warnf("Could not parse the request params")
		}

		next.ServeHTTP(lrw, r.WithContext(ctx))

		logger = logger.WithFields(sdklogger.Fields{
			"http": sdklogger.Fields{
				"remote_addr":            r.RemoteAddr,
				"request_id":             requestID,
				"request_ip":             r.Header.Get(ForwardedForHeader),
				"request_method":         r.Method,
				"request_path":           r.URL.EscapedPath(),
				"request_fullpath":       r.URL.RequestURI(),
				"request_params":         r.Form,
				"request_user_agent":     r.UserAgent(),
				"response_status":        lrw.StatusCode,
				"response_time_total_ms": time.Since(start).Milliseconds(),
			},
			"dd": sdklogger.Fields{
				"trace_id": logContext.TraceID,
				"span_id":  logContext.SpanID,
			},
		})

		// Format the message in a similar way to Common Log Format
		message := fmt.Sprintf("%s %s %s %d",
			r.Method,
			r.URL.EscapedPath(),
			r.Proto,
			lrw.StatusCode,
		)

		switch {
		case lrw.StatusCode >= 400 && lrw.StatusCode <= 499:
			logger.Warnf(message)
		case lrw.StatusCode >= 500 && lrw.StatusCode <= 599:
			logger.Errorf(message)
		default:
			logger.Infof(message)
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader uses the wrapped http.ResponseWriter to write the given code.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}
