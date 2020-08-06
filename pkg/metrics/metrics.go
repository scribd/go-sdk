/*
Package metrics provides a configured Datadog Statsd client.

`datadog-go` is the official library that provides a [DogStatsD](ddsd) client.

The following documentation is available:

- [GoDoc documentation for Datadog Go](ddgo).
- [Official Datadog DogStatsD documentation](ddsd).

[ddgo]: <http://godoc.org/github.com/DataDog/datadog-go/statsd>
[ddsd]: <https://docs.datadoghq.com/developers/dogstatsd/?tab=go>
*/
package metrics

import "time"

// Metrics is an interface that exposes the common client functions for sending metrics.
type Metrics interface {
	// Gauge measures the value of a metric at a particular time.
	Gauge(name string, value float64, tags []string, rate float64) error

	// Count tracks how many times something happened per second.
	Count(name string, value int64, tags []string, rate float64) error

	// Histogram tracks the statistical distribution of a set of values on each host.
	Histogram(name string, value float64, tags []string, rate float64) error

	// Distribution tracks the statistical distribution of a set of values across your infrastructure.
	Distribution(name string, value float64, tags []string, rate float64) error

	// Decr is just Count of -1
	Decr(name string, tags []string, rate float64) error

	// Incr is just Count of 1
	Incr(name string, tags []string, rate float64) error

	// Set counts the number of unique elements in a group.
	Set(name string, value string, tags []string, rate float64) error

	// Timing sends timing information, it is an alias for TimeInMilliseconds
	Timing(name string, value time.Duration, tags []string, rate float64) error

	// TimeInMilliseconds sends timing information in milliseconds.
	// It is flushed by statsd with percentiles, mean and other info
	// (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
	TimeInMilliseconds(name string, value float64, tags []string, rate float64) error

	// SimpleEvent sends an event with the provided title and text.
	SimpleEvent(title, text string) error

	// Close the client connection.
	Close() error

	// Flush forces a flush of all the queued payloads.
	Flush() error

	// SetWriteTimeout allows the user to set a custom write timeout.
	SetWriteTimeout(d time.Duration) error
}
