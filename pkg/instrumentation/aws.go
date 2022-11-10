package instrumentation

import (
	"fmt"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go/aws"
)

// Settings stores DataDog instrumentation settings.
type Settings struct {
	AppName       string
	Analytics     bool
	AnalyticsRate float64
}

// InstrumentAWSSession configures DD tracing mode.
func InstrumentAWSSession(session *awssession.Session, settings Settings) *awssession.Session {
	return awstrace.WrapSession(
		session,
		awstrace.WithServiceName(fmt.Sprintf("%s-app", settings.AppName)),
		awstrace.WithAnalytics(settings.Analytics),
		awstrace.WithAnalyticsRate(settings.AnalyticsRate),
	)
}
