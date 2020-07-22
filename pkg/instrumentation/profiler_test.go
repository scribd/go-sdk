package instrumentation

import (
	"testing"
)

// TestProfilerStop checks that Stop called on a disabled profiler
// does not cause any problems.
func TestProfilerStop(t *testing.T) {
	p := NewProfiler(&Config{Enabled: false})
	p.Stop()
}
