package apps

import (
	"log"
	"os"
	"testing"

	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

const (
	appRootEnvKey   = "APP_ROOT"
	expectedTestEnv = "test"
	testdataDir     = "testdata/config"
)

// newTestBuilder setups a builder for testdata.
func newTestBuilder(confName, appName string, t *testing.T) builder.Builder {
	builder, err := builder.NewViper(testdataDir, appName, "test", getAppRootMust())
	if err != nil {
		t.Fatalf("new builder, err: %s", err.Error())
	}

	builder.SetConfigName(confName)

	return builder
}

func getAppRootMust() string {
	appRoot := os.Getenv(appRootEnvKey)
	if appRoot == "" {
		log.Fatalf("env key %s missing", appRootEnvKey)
	}

	return appRoot
}
