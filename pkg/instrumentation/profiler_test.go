package instrumentation

import (
	"testing"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
)

// TestProfilerStop checks that Stop called on a disabled profiler
// does not cause any problems.
func TestProfilerStop(t *testing.T) {
	p := NewProfiler(apps.Instrumentation{Enabled: false})
	p.Stop()
}
