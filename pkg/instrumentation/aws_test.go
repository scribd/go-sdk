package instrumentation

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestInstrumentAWSSession(t *testing.T) {
	cfg := aws.NewConfig().
		WithRegion("us-west-2").
		WithDisableSSL(true).
		WithCredentials(credentials.AnonymousCredentials)

	sess := session.Must(session.NewSession(cfg))
	sess = InstrumentAWSSession(sess, "testApp")

	var (
		tagAWSAgent     = "aws.agent"
		tagAWSOperation = "aws.operation"
		tagAWSRegion    = "aws.region"
	)

	t.Run("s3", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		root, ctx := tracer.StartSpanFromContext(context.Background(), "test")

		s3api := s3.New(sess)
		_, err := s3api.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String("test-bucket-name"),
			Key:    aws.String("//test//file//name"),
		})

		require.NotNil(t, err)
		root.Finish()

		spans := mt.FinishedSpans()
		assert.Len(t, spans, 2)
		assert.Equal(t, spans[1].TraceID(), spans[0].TraceID())

		s := spans[0]
		assert.Equal(t, "s3.command", s.OperationName())
		assert.Contains(t, s.Tag(tagAWSAgent), "aws-sdk-go")
		assert.Equal(t, "GetObject", s.Tag(tagAWSOperation))
		assert.Equal(t, "us-west-2", s.Tag(tagAWSRegion))
		assert.Equal(t, "s3.GetObject", s.Tag(ext.ResourceName))
		assert.Equal(t, "testApp-aws", s.Tag(ext.ServiceName))
		assert.Equal(t, "GET", s.Tag(ext.HTTPMethod))
		assert.Equal(t, "http://test-bucket-name.s3.us-west-2.amazonaws.com/test/file/name", s.Tag(ext.HTTPURL))
	})
}
