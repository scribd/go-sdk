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

// TestNewProfilerEnabled verifies the profiler runs only when instrumentation
// is enabled and profiling is not opted out (ProfilerEnabled defaults to true
// in NewConfig, so profiling is on wherever Enabled is on unless disabled).
func TestNewProfilerEnabled(t *testing.T) {
	testCases := []struct {
		name            string
		enabled         bool
		profilerEnabled bool
		want            bool
	}{
		{name: "on when instrumentation on and not opted out", enabled: true, profilerEnabled: true, want: true},
		{name: "opted out while tracing stays on", enabled: true, profilerEnabled: false, want: false},
		{name: "off when instrumentation off", enabled: false, profilerEnabled: true, want: false},
		{name: "both off", enabled: false, profilerEnabled: false, want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewProfiler(&Config{Enabled: tc.enabled, ProfilerEnabled: tc.profilerEnabled})
			if p.enabled != tc.want {
				t.Errorf("profiler enabled = %v, want %v", p.enabled, tc.want)
			}
		})
	}
}
