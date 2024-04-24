package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/aws/aws-sdk-go-v2/aws"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"

	awssasl "github.com/twmb/franz-go/pkg/sasl/aws"
	"github.com/twmb/franz-go/pkg/sasl/plain"

	sdklogger "github.com/scribd/go-sdk/pkg/logger"
	"github.com/scribd/go-sdk/pkg/pubsub"
)

// Config provides a common configuration for Kafka PubSub clients.
type Config struct {
	// Application name that will be used in a serviceName provided to tracer spans
	ApplicationName string
	// Kafka configuration provided by go-sdk
	KafkaConfig pubsub.Kafka
	// AWS session reference, it will be used in case AWS MSK IAM authentication mechanism is used
	//
	// Deprecated: Use AwsConfig instead
	AwsSession *session.Session
	// MsgHandler is a function that will be called when a message is received
	MsgHandler MsgHandler
	// AWS configuration reference, it will be used in case AWS MSK IAM authentication mechanism is used
	AwsConfig *aws.Config
	Logger    sdklogger.Logger
}

const tlsConnectionTimeout = 10 * time.Second

func newConfig(c Config, opts ...kgo.Opt) ([]kgo.Opt, error) {
	options := []kgo.Opt{
		kgo.SeedBrokers(c.KafkaConfig.BrokerUrls...),
		kgo.ClientID(c.KafkaConfig.ClientId),
	}

	if c.KafkaConfig.SASL.Enabled {
		switch c.KafkaConfig.SASLMechanism() {
		case pubsub.Plain:
			options = append(options, getPlainSaslOption(c.KafkaConfig.SASL))
		case pubsub.AWSMskIam:
			options = append(options, getAwsMskIamSaslOption(c.KafkaConfig.SASL.AWSMskIam, c.AwsSession, c.AwsConfig))
		}
	}

	if c.KafkaConfig.TLS.Enabled || c.KafkaConfig.SecurityProtocol == "ssl" {
		var caCertPool *x509.CertPool

		if c.KafkaConfig.TLS.Ca != "" {
			caCertPool = x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(c.KafkaConfig.TLS.Ca))
		}

		var certificates []tls.Certificate
		if c.KafkaConfig.TLS.Cert != "" && c.KafkaConfig.TLS.CertKey != "" {
			cert, err := tls.X509KeyPair([]byte(c.KafkaConfig.TLS.Cert), []byte(c.KafkaConfig.TLS.CertKey))
			if err != nil {
				return nil, err
			}
			certificates = []tls.Certificate{cert}
		}

		if c.KafkaConfig.Cert != "" && c.KafkaConfig.CertKey != "" {
			cert, err := tls.X509KeyPair([]byte(c.KafkaConfig.Cert), []byte(c.KafkaConfig.CertKey))
			if err != nil {
				return nil, err
			}
			certificates = []tls.Certificate{cert}
		}

		var skipTLSVerify bool
		if c.KafkaConfig.TLS.InsecureSkipVerify || !c.KafkaConfig.SSLVerificationEnabled {
			skipTLSVerify = true
		}

		tlsDialer := &tls.Dialer{
			NetDialer: &net.Dialer{Timeout: tlsConnectionTimeout},
			Config: &tls.Config{
				InsecureSkipVerify: skipTLSVerify,
				Certificates:       certificates,
				RootCAs:            caCertPool,
			},
		}

		options = append(options, kgo.Dialer(tlsDialer.DialContext))
	}

	options = append(options, opts...)

	return options, nil
}

func getPlainSaslOption(saslConf pubsub.SASL) kgo.Opt {
	return kgo.SASL(plain.Auth{
		User: saslConf.Username,
		Pass: saslConf.Password,
	}.AsMechanism())
}

func getAwsMskIamSaslOption(iamConf pubsub.SASLAwsMskIam, s *session.Session, awsCfg *aws.Config) kgo.Opt {
	var opt kgo.Opt

	// no AWS session and AWS config provided
	if s == nil && awsCfg == nil {
		opt = kgo.SASL(awssasl.Auth{
			AccessKey:    iamConf.AccessKey,
			SecretKey:    iamConf.SecretKey,
			SessionToken: iamConf.SessionToken,
			UserAgent:    iamConf.UserAgent,
		}.AsManagedStreamingIAMMechanism())
	} else {
		opt = kgo.SASL(
			awssasl.ManagedStreamingIAM(func(ctx context.Context) (awssasl.Auth, error) {
				if s != nil {
					return getAwsSaslAuthFromSession(iamConf, s)
				}

				return getAwsSaslAuthFromConfig(ctx, iamConf, awsCfg)
			}),
		)
	}

	return opt
}

func getAwsSaslAuthFromSession(iamConf pubsub.SASLAwsMskIam, s *session.Session) (awssasl.Auth, error) {
	// If assumable role is not provided, we try to get credentials from the provided AWS session
	if iamConf.AssumableRole == "" {
		val, err := s.Config.Credentials.Get()
		if err != nil {
			return awssasl.Auth{}, err
		}

		return awssasl.Auth{
			AccessKey:    val.AccessKeyID,
			SecretKey:    val.SecretAccessKey,
			SessionToken: val.SessionToken,
			UserAgent:    iamConf.UserAgent,
		}, nil
	}

	svc := sts.New(s)

	res, stsErr := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &iamConf.AssumableRole,
		RoleSessionName: &iamConf.SessionName,
	})
	if stsErr != nil {
		return awssasl.Auth{}, stsErr
	}

	return awssasl.Auth{
		AccessKey:    *res.Credentials.AccessKeyId,
		SecretKey:    *res.Credentials.SecretAccessKey,
		SessionToken: *res.Credentials.SessionToken,
		UserAgent:    iamConf.UserAgent,
	}, nil
}

func getAwsSaslAuthFromConfig(
	ctx context.Context,
	iamConf pubsub.SASLAwsMskIam,
	awsCfg *aws.Config) (awssasl.Auth, error) {
	// If assumable role is not provided, we try to get credentials from the provided AWS config
	if iamConf.AssumableRole == "" {
		val, err := awsCfg.Credentials.Retrieve(ctx)
		if err != nil {
			return awssasl.Auth{}, err
		}

		return awssasl.Auth{
			AccessKey:    val.AccessKeyID,
			SecretKey:    val.SecretAccessKey,
			SessionToken: val.SessionToken,
			UserAgent:    iamConf.UserAgent,
		}, nil
	}

	client := stsv2.NewFromConfig(*awsCfg)

	res, stsErr := client.AssumeRole(ctx, &stsv2.AssumeRoleInput{
		RoleArn:         &iamConf.AssumableRole,
		RoleSessionName: &iamConf.SessionName,
	})
	if stsErr != nil {
		return awssasl.Auth{}, stsErr
	}

	return awssasl.Auth{
		AccessKey:    *res.Credentials.AccessKeyId,
		SecretKey:    *res.Credentials.SecretAccessKey,
		SessionToken: *res.Credentials.SessionToken,
		UserAgent:    iamConf.UserAgent,
	}, nil
}
