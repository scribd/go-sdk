package builder

import "time"

// Builder type is used to manuplate Conf settings before building it.
type Builder interface {
	// SetConfigName is set when a file should be fetched.
	SetConfigName(string)
	SetDefault(string, interface{})

	Build() (Conf, error)
}

// Conf is a store to retrieve basic data.
type Conf interface {
	Unmarshal(interface{}) error

	Bool(key string) bool
	Float64(key string) float64
	Int(key string) int
	String(key string) string
	StringMap(key string) map[string]interface{}
	StringMapString(key string) map[string]string
	StringSlice(key string) []string
	Time(key string) time.Time
	Duration(key string) time.Duration
	Set(key string, value interface{})
	IsSet(key string) bool
	AllSettings() map[string]interface{}
}
