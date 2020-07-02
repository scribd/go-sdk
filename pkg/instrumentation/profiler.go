package instrumentation

import "gopkg.in/DataDog/dd-trace-go.v1/profiler"

func WithProfiler(appName, appEnv string) (func(), error) {
	err := profiler.Start(
		profiler.WithService(appName),
		profiler.WithEnv(appEnv),
	)

	return profiler.Stop, err
}
