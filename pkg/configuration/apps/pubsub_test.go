package apps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubsubFetchConfig(t *testing.T) {
	testCases := []struct {
		name      string
		env       string
		filename  string
		wantError bool

		kafka        Kafka
		envOverrides [][]string
	}{
		{
			name:      "NewWithoutConfigFileFails",
			env:       "test",
			filename:  "imaginarypubsub",
			wantError: true,
		},
		{
			name:     "NewWithConfigFileWorks",
			env:      "test",
			filename: "pubsub",
			kafka: Kafka{
				BrokerUrls:             []string{"localhost:9092"},
				ClientId:               "test-app",
				Cert:                   "pem string",
				CertKey:                "pem key",
				SecurityProtocol:       "ssl",
				SSLVerificationEnabled: true,
				Publisher: Publisher{
					MaxAttempts:  3,
					WriteTimeout: 10 * time.Second,
					Topic:        "test-topic",
				},
				Subscriber: Subscriber{
					Topic:   "test-topic",
					GroupId: "sub-group-id",
				},
			},
		},
		{
			name:     "NewWithConfigFileWorks (broker URLs override)",
			env:      "test",
			filename: "pubsub",
			kafka: Kafka{
				BrokerUrls:       []string{"localhost:9092", "localhost:9093"},
				ClientId:         "test-app",
				Cert:             "pem string",
				CertKey:          "pem key",
				SecurityProtocol: "ssl",
				Publisher: Publisher{
					MaxAttempts:  3,
					WriteTimeout: 10 * time.Second,
					Topic:        "test-topic",
				},
				Subscriber: Subscriber{
					Topic:   "test-topic",
					GroupId: "sub-group-id",
				},
				SSLVerificationEnabled: true,
			},

			envOverrides: [][]string{{"APP_PUBSUB_KAFKA_BROKER_URLS", "localhost:9092 localhost:9093"}},
		},
		{
			name:     "NewWithConfigFileWorks (TLS config override)",
			env:      "test",
			filename: "pubsub",
			kafka: Kafka{
				BrokerUrls:       []string{"localhost:9092"},
				ClientId:         "test-app",
				Cert:             "pem string",
				CertKey:          "pem key",
				SecurityProtocol: "ssl",
				Publisher: Publisher{
					MaxAttempts:  3,
					WriteTimeout: 10 * time.Second,
					Topic:        "test-topic",
				},
				Subscriber: Subscriber{
					Topic:   "test-topic",
					GroupId: "sub-group-id",
				},
				SSLVerificationEnabled: true,
				TLS: TLS{
					Enabled: true,
					Cert:    "pem string",
					CertKey: "pem key",
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_KAFKA_TLS_ENABLED", "true"},
				{"APP_PUBSUB_KAFKA_TLS_CERT_PEM", "pem string"},
				{"APP_PUBSUB_KAFKA_TLS_CERT_PEM_KEY", "pem key"},
			},
		},
		{
			name:     "NewWithConfigFileWorks (SASL config override, error)",
			env:      "test",
			filename: "pubsub",
			kafka: Kafka{
				BrokerUrls:       []string{"localhost:9092"},
				ClientId:         "test-app",
				Cert:             "pem string",
				CertKey:          "pem key",
				SecurityProtocol: "ssl",
				Publisher: Publisher{
					MaxAttempts:  3,
					WriteTimeout: 10 * time.Second,
					Topic:        "test-topic",
				},
				Subscriber: Subscriber{
					Topic:   "test-topic",
					GroupId: "sub-group-id",
				},
				SASL: SASL{
					Enabled:   true,
					Mechanism: "test",
				},
				SSLVerificationEnabled: true,
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_KAFKA_SASL_ENABLED", "true"},
				{"APP_PUBSUB_KAFKA_SASL_MECHANISM", "test"},
			},
			wantError: true,
		},
		{
			name:     "NewWithConfigFileWorks (SASL config override, error)",
			env:      "test",
			filename: "pubsub",
			kafka: Kafka{
				BrokerUrls:       []string{"localhost:9092"},
				ClientId:         "test-app",
				Cert:             "pem string",
				CertKey:          "pem key",
				SecurityProtocol: "ssl",
				Publisher: Publisher{
					MaxAttempts:  3,
					WriteTimeout: 10 * time.Second,
					Topic:        "test-topic",
				},
				Subscriber: Subscriber{
					Topic:   "test-topic",
					GroupId: "sub-group-id",
				},
				SASL: SASL{
					Enabled:   true,
					Mechanism: "aws_msk_iam",
					AWSMskIam: SASLAwsMskIam{
						AccessKey: "access key",
						SecretKey: "secret key",
					},
				},
				SSLVerificationEnabled: true,
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_KAFKA_SASL_ENABLED", "true"},
				{"APP_PUBSUB_KAFKA_SASL_MECHANISM", "aws_msk_iam"},
				{"APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_ACCESS_KEY", "access key"},
				{"APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_SECRET_KEY", "secret key"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.envOverrides) > 0 {
				for _, o := range tc.envOverrides {
					t.Setenv(o[0], o[1])
				}
			}

			b := newTestBuilder(tc.filename, "pubsub", t)

			pb := &PubSub{}
			if err := pb.FetchConfig(b); err != nil {
				if tc.wantError {
					return
				}

				t.Fatalf("unmarshalling config, %s", err.Error())
			}

			assert.Equal(t, tc.kafka, pb.Kafka)
		})
	}
}
