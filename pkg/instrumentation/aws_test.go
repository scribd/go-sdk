package instrumentation

import (
	"context"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awss3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type (
	testCustomResolver struct{}
)

func (t testCustomResolver) ResolveEndpoint(
	ctx context.Context, params awss3v2.EndpointParameters) (smithyendpoints.Endpoint, error) {
	uri, err := url.Parse("http://localhost:4566")
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}

	return smithyendpoints.Endpoint{
		URI: *uri,
	}, nil

}

func TestInstrumentAWSSession(t *testing.T) {
	cfg := aws.NewConfig().
		WithRegion("us-west-2").
		WithDisableSSL(true).
		WithCredentials(credentials.AnonymousCredentials)

	sess := session.Must(session.NewSession(cfg))
	sess = InstrumentAWSSession(sess, Settings{AppName: "testApp"})

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

func TestInstrumentAWSClient(t *testing.T) {
	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion("us-west-2"),
		awscfg.WithCredentialsProvider(awsv2.AnonymousCredentials{}),
	)
	require.NoError(t, err)

	InstrumentAWSClient(&cfg, Settings{AppName: "testApp"})

	var (
		tagComponent    = "component"
		tagAWSOperation = "aws.operation"
		tagAWSRegion    = "aws.region"
	)

	t.Run("s3", func(t *testing.T) {
		mt := mocktracer.Start()
		defer mt.Stop()

		client := awss3v2.NewFromConfig(cfg, awss3v2.WithEndpointResolverV2(&testCustomResolver{}))
		root, ctx := tracer.StartSpanFromContext(context.Background(), "test")

		_, err := client.GetObject(ctx, &awss3v2.GetObjectInput{
			Bucket: aws.String("test-bucket-name"),
			Key:    aws.String("//test//file//name"),
		})

		require.NotNil(t, err)
		root.Finish()

		spans := mt.FinishedSpans()
		assert.Len(t, spans, 2)
		assert.Equal(t, spans[1].TraceID(), spans[0].TraceID())

		s := spans[0]
		assert.Equal(t, "S3.request", s.OperationName())
		assert.Contains(t, s.Tag(tagComponent), "aws/aws-sdk-go-v2/aws")
		assert.Equal(t, "GetObject", s.Tag(tagAWSOperation))
		assert.Equal(t, "us-west-2", s.Tag(tagAWSRegion))
		assert.Equal(t, "S3.GetObject", s.Tag(ext.ResourceName))
		assert.Equal(t, "testApp-aws", s.Tag(ext.ServiceName))
		assert.Equal(t, "GET", s.Tag(ext.HTTPMethod))
		assert.Equal(
			t,
			"http://localhost:4566///test//file//name?x-id=GetObject",
			s.Tag(ext.HTTPURL),
		)
	})
}
