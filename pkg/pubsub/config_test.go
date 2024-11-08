package pubsub

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "NewWithoutConfigFileFails",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfig()

			gotError := err != nil
			assert.Equal(t, gotError, tc.wantError)
		})
	}
}

func TestNewConfigWithAppRoot(t *testing.T) {
	testCases := []struct {
		name    string
		env     string
		cfg     Config
		wantErr bool

		envOverrides [][]string
	}{
		{
			name: "NewWithConfigFileWorks",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},
		},
		{
			name: "NewWithConfigFileWorks (broker URLs override)",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},

			envOverrides: [][]string{{"APP_PUBSUB_KAFKA_BROKER_URLS", "localhost:9092 localhost:9093"}},
		},
		{
			name: "NewWithConfigFileWorks (auto_commimt override)",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: false,
						},
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},

			envOverrides: [][]string{{"APP_PUBSUB_KAFKA_SUBSCRIBER_AUTO_COMMIT_ENABLED", "false"}},
		},
		{
			name: "NewWithConfigFileWorks (TLS config override)",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SSLVerificationEnabled: true,
					TLS: TLS{
						Enabled: true,
						Cert:    "pem string",
						CertKey: "pem key",
					},
				},
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_KAFKA_TLS_ENABLED", "true"},
				{"APP_PUBSUB_KAFKA_TLS_CERT_PEM", "pem string"},
				{"APP_PUBSUB_KAFKA_TLS_CERT_PEM_KEY", "pem key"},
			},
		},
		{
			name: "NewWithConfigFileWorks (SASL config override, error)",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SASL: SASL{
						Enabled:   true,
						Mechanism: "test",
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_KAFKA_SASL_ENABLED", "true"},
				{"APP_PUBSUB_KAFKA_SASL_MECHANISM", "test"},
			},
			wantErr: true,
		},
		{
			name: "NewWithConfigFileWorks (SASL config override, error)",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
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
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_KAFKA_SASL_ENABLED", "true"},
				{"APP_PUBSUB_KAFKA_SASL_MECHANISM", "aws_msk_iam"},
				{"APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_ACCESS_KEY", "access key"},
				{"APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_SECRET_KEY", "secret key"},
			},
		},
		{
			name: "NewWithConfigFileWorks (SQS overrides)",
			env:  "test",
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{
						QueueURL: "https://test2.com",
					},
					Subscriber: SQSSubscriber{
						QueueURL:    "https://test3.com",
						MaxMessages: 10,
						Workers:     5,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_SQS_PUBLISHER_QUEUE_URL", "https://test2.com"},
				{"APP_PUBSUB_SQS_SUBSCRIBER_QUEUE_URL", "https://test3.com"},
				{"APP_PUBSUB_SQS_SUBSCRIBER_MAX_MESSAGES", "10"},
				{"APP_PUBSUB_SQS_SUBSCRIBER_WORKERS", "5"},
			},
		},
		{
			name:    "NewWithConfigFileWorks (SQS override, empty publisher QueueURL)",
			env:     "test",
			wantErr: true,
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{
						Enabled: true,
					},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_SQS_PUBLISHER_ENABLED", "true"},
			},
		},
		{
			name:    "NewWithConfigFileWorks (SQS override, empty subscriber QueueURL)",
			env:     "test",
			wantErr: true,
			cfg: Config{
				Kafka: Kafka{
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
						Topic: "test-topic",
						AutoCommit: AutoCommit{
							Enabled: true,
						},
					},
					SSLVerificationEnabled: true,
				},
				SQS: SQS{
					Publisher: SQSPublisher{},
					Subscriber: SQSSubscriber{
						MaxMessages: 1,
						Workers:     1,
						Enabled:     true,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_PUBSUB_SQS_SUBSCRIBER_ENABLED", "true"},
			},
		},
	}

	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var envVariables [][]string

			if len(tc.envOverrides) > 0 {
				for _, o := range tc.envOverrides {
					currentVal := os.Getenv(o[0])
					envVariables = append(envVariables, []string{o[0], currentVal})

					os.Setenv(o[0], o[1])
				}
			}

			_, filename, _, _ := runtime.Caller(0)
			tmpRootParent := filepath.Dir(filename)
			os.Setenv("APP_ROOT", filepath.Join(tmpRootParent, "testdata"))

			c, err := NewConfig()
			if tc.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			assert.Equal(t, tc.cfg, *c)

			// teardown
			if len(envVariables) > 0 {
				for _, o := range envVariables {
					os.Setenv(o[0], o[1])
				}
			}
		})
	}
}
