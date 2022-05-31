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
		kafka   Kafka
		wantErr bool

		envOverrides [][]string
	}{
		{
			name: "NewWithConfigFileWorks",
			env:  "test",
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
					Topic: "test-topic",
				},
				SSLVerificationEnabled: true,
			},
		},
		{
			name: "NewWithConfigFileWorks (broker URLs override)",
			env:  "test",
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
					Topic: "test-topic",
				},
				SSLVerificationEnabled: true,
			},

			envOverrides: [][]string{{"APP_PUBSUB_KAFKA_BROKER_URLS", "localhost:9092 localhost:9093"}},
		},
		{
			name: "NewWithConfigFileWorks (TLS config override)",
			env:  "test",
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
					Topic: "test-topic",
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
			name: "NewWithConfigFileWorks (SASL config override, error)",
			env:  "test",
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
					Topic: "test-topic",
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
			wantErr: true,
		},
		{
			name: "NewWithConfigFileWorks (SASL config override, error)",
			env:  "test",
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
					Topic: "test-topic",
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

			assert.Equal(t, tc.kafka, c.Kafka)

			// teardown
			if len(envVariables) > 0 {
				for _, o := range envVariables {
					os.Setenv(o[0], o[1])
				}
			}
		})
	}
}
