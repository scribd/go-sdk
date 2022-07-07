package builder

import (
	"os"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

const (
	testAppName = "test"
	testAppEnv  = "test"
	testConfDir = "testdata/config/builder"
)

type viperTestData struct {
	Foo string `mapstructure:"foo"`
}

func TestViperConf(t *testing.T) {
	var (
		appRoot             = os.Getenv("APP_ROOT")
		filename            = "getters"
		expectedBool        = true
		expectedFloat64     = 100.59
		expectedInt         = 42
		expectedString      = "foo"
		expectedStringSlice = strings.Join([]string{"foo", "bar"}, ",")
		expectedTime, _     = time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	)

	b, err := NewViper(testConfDir, testAppName, testAppEnv, appRoot)
	if err != nil {
		t.Fatalf("new viper builder test, %s", err.Error())
	}

	b.SetConfigName(filename)

	config, err := b.Build()
	if err != nil {
		t.Fatalf("new viper builder test, %s", err.Error())
	}

	if actual := config.Bool("bool"); actual != expectedBool {
		t.Errorf("Got: %t, expected: %t", actual, expectedBool)
	}

	if actual := config.Float64("float64"); actual != expectedFloat64 {
		t.Errorf("Got: %f, expected: %f", actual, expectedFloat64)
	}

	if actual := config.Int("int"); actual != expectedInt {
		t.Errorf("Got: %d, expected: %d", actual, expectedInt)
	}

	if actual := config.String("string"); actual != expectedString {
		t.Errorf("Got: %s, expected: %s", actual, expectedString)
	}

	actualStrSlice := config.StringSlice("string_slice")

	if strings.Join(actualStrSlice, ",") != expectedStringSlice {
		t.Errorf("Got slice: %s, expected: %s", actualStrSlice, expectedStringSlice)
	}

	if actual := config.Time("time"); actual != expectedTime {
		t.Errorf("Got time: %s, expected time: %s", actual, expectedTime)
	}
}

func TestViperBuilder(t *testing.T) {
	appRoot := os.Getenv("APP_ROOT")

	t.Run("Build", func(t *testing.T) {
		cases := []struct {
			name       string
			configName string
			wantError  bool

			foo string
		}{
			{
				name:       "ValidFile",
				configName: "valid",
				wantError:  false,
				foo:        "bar",
			},
			{
				name:       "ValidWithoutEnvs",
				configName: "valid-no-envs",
				wantError:  true,
			},
			{
				name:       "InvalidFile",
				configName: "invalid",
				wantError:  true,
			},
			{
				name:       "NonExistentFile",
				configName: "nonexistent",
				wantError:  true,
			},
		}

		for _, ct := range cases {
			t.Run(ct.name, func(t *testing.T) {
				b, err := NewViper(testConfDir, testAppName, testAppEnv, appRoot)
				if err != nil {
					if ct.wantError {
						return
					}

					t.Fatalf("new viper builder test, %s", err.Error())
				}

				b.SetConfigName(ct.configName)

				config, err := b.Build()
				if err != nil {
					if ct.wantError {
						return
					}

					t.Fatalf("new viper builder test, %s", err.Error())
				}

				var test viperTestData
				if err := config.Unmarshal(&test); err != nil {
					if ct.wantError {
						return
					}

					t.Fatalf("unmarshalling viper builder test, %s", err.Error())
				}

				assert.Equal(t, ct.foo, test.Foo)
			})
		}
	})
}
