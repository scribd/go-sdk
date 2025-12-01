package instrumentation

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	awstracev2 "github.com/DataDog/dd-trace-go/contrib/aws/aws-sdk-go-v2/v2/aws"
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
		awstracev2.WithService(fmt.Sprintf("%s-aws", settings.AppName)),
		awstracev2.WithAnalytics(settings.Analytics),
		awstracev2.WithAnalyticsRate(settings.AnalyticsRate),
	)
}
