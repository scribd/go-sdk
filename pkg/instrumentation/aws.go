package instrumentation

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"

	awstracev2 "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go-v2/aws"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go/aws"
)

// Settings stores DataDog instrumentation settings.
type Settings struct {
	AppName       string
	Analytics     bool
	AnalyticsRate float64
}

// InstrumentAWSSession configures DD tracing mode.
//
// Deprecated: Use InstrumentAWSClient instead.
func InstrumentAWSSession(session *awssession.Session, settings Settings) *awssession.Session {
	return awstrace.WrapSession(
		session,
		awstrace.WithServiceName(fmt.Sprintf("%s-aws", settings.AppName)),
		awstrace.WithAnalytics(settings.Analytics),
		awstrace.WithAnalyticsRate(settings.AnalyticsRate),
	)
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
