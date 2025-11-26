package instrumentation

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	awstracev2 "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go-v2/aws"
)

// Settings stores DataDog instrumentation settings.
type Settings struct {
	AppName       string
	Analytics     bool
	AnalyticsRate float64
}

// InstrumentAWSClient configures DD tracing mode.
func InstrumentAWSClient(cfg *aws.Config, settings Settings) {
	awstracev2.AppendMiddleware(
		cfg,
		awstracev2.WithServiceName(fmt.Sprintf("%s-aws", settings.AppName)),
		awstracev2.WithAnalytics(settings.Analytics),
		awstracev2.WithAnalyticsRate(settings.AnalyticsRate),
	)
}
