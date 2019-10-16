package logger

// Fields is the struct that that stores key/value pairs for structured logs.
type Fields map[string]interface{}

// Set sets a key/value pair in the Fields map.
func (f *Fields) Set(key string, value interface{}) {
	map[string]interface{}(*f)[key] = value
}
