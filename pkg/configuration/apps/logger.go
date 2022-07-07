package apps

import (
	"fmt"
	"os"
	"path"

	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

// Logger stores the configuration for the logger.
// For some loggers there can only be one level across writers, for such
// the level of Console is picked by default.
type Logger struct {
	ConsoleEnabled    bool   `mapstructure:"console_enabled"`
	ConsoleJSONFormat bool   `mapstructure:"console_json_format"`
	ConsoleLevel      string `mapstructure:"console_level"`
	FileEnabled       bool   `mapstructure:"file_enabled"`
	FileJSONFormat    bool   `mapstructure:"file_json_format"`
	FileLevel         string `mapstructure:"file_level"`
	FileLocation      string `mapstructure:"file_location"`
	FileName          string `mapstructure:"file_name"`
}

// FetchConfig fetches data from builder into it's caller.
func (l *Logger) FetchConfig(c builder.Builder) error {
	c.SetDefault("file_location", path.Join(os.Getenv("APP_ROOT"), "log"))
	c.SetDefault("file_name", fmt.Sprintf("%s.log", os.Getenv("APP_ENV")))

	config, err := c.Build()
	if err != nil {
		return fmt.Errorf("unable to build for logger. err: %w", err)
	}

	if err = config.Unmarshal(l); err != nil {
		return fmt.Errorf("Unable to decode into logger struct: %w", err)
	}

	return nil
}
