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
		name  string
		env   string
		kafka Kafka
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
	}

	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, filename, _, _ := runtime.Caller(0)
			tmpRootParent := filepath.Dir(filename)
			os.Setenv("APP_ROOT", filepath.Join(tmpRootParent, "testdata"))

			c, err := NewConfig()
			require.Nil(t, err)

			assert.Equal(t, c.Kafka, tc.kafka)
		})
	}
}
