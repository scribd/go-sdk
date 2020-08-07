package contextkeys

const (
	// Logger is the context key for carrying the Logger
	Logger = "Logger"

	// RequestID is the context key for the carrying the RequestID.
	RequestID = "RequestID"

	// Database is the context key for the carrying the Database connection pool.
	Database = "Database"

	// Metrics is the context key for the carrying the Metrics client.
	Metrics = "Metrics"
)
