package instrumentation

import (
	"fmt"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go/aws"
)

func InstrumentAWSSession(s *awssession.Session) *awssession.Session {
	return awstrace.WrapSession(
		s,
	)
}
