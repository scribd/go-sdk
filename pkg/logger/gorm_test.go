package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewGormLogger(t *testing.T) {
	sampleQuery := "SELECT id FROM users"

	testDBErrorMsg := "test error"
	testDBError := errors.New(testDBErrorMsg)

	var expectedEffectedRows int64 = 1
	var expectedLastInsertID int64 = 1
	// expectedEffectedRowsString := strconv.Itoa(int(expectedEffectedRows))

	tests := []struct {
		name        string
		isLogged    bool
		resultError bool
		cfg         Config
	}{
		{
			name:        "Empty on fatal log level with error",
			isLogged:    false,
			resultError: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "fatal",
			},
		},
		{
			name:        "Empty on fatal log level with error",
			isLogged:    false,
			resultError: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "fatal",
			},
		},
		{
			name:        "Empty on error log level without error",
			isLogged:    false,
			resultError: false,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "error",
			},
		},
		{
			name:        "Print on error log level with error",
			isLogged:    true,
			resultError: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "error",
			},
		},
		{
			name:        "Empty on debug log level without error",
			isLogged:    false,
			resultError: false,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "debug",
			},
		},
		{
			name:        "Print on debug log level with error",
			isLogged:    true,
			resultError: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "debug",
			},
		},
		{
			name:        "Print on trace log level without error",
			isLogged:    true,
			resultError: false,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "trace",
			},
		},
		{
			name:        "Print on trace log level with error",
			isLogged:    true,
			resultError: true,
			cfg: Config{
				ConsoleEnabled:    true,
				ConsoleJSONFormat: true,
				ConsoleLevel:      "trace",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			b := NewBuilder(&tt.cfg)
			l, err := b.BuildTestLogger(&buffer)
			require.Nil(t, err)

			gormDB, mock := mockGormConnectionWithLogger(t, l)
			ee := mock.ExpectExec(sampleQuery)

			if tt.resultError {
				ee.WillReturnError(testDBError)
			} else {
				ee.WillReturnResult(sqlmock.NewResult(expectedLastInsertID, expectedEffectedRows))
			}

			gormDB.Exec(sampleQuery)

			if tt.isLogged {
				var fields map[string]interface{}
				err := json.Unmarshal(buffer.Bytes(), &fields)
				assert.Nil(t, err)

				assert.Contains(t, buffer.String(), "level")

				assert.Contains(t, buffer.String(), "message")
				assert.Equal(t, fields["message"], gormLoggerMsg)

				assert.Contains(t, buffer.String(), "sql")
				assert.Contains(t, fields["sql"], "sql")
				assert.Contains(t, fields["sql"], "duration")
				assert.Contains(t, fields["sql"], "affected_rows")

				if tt.resultError {
					assert.Equal(t, fields["level"], "error")
					assert.Contains(t, buffer.String(), "error")
					assert.Equal(t, fields["error"], testDBErrorMsg)
				} else {
					assert.Equal(t, fields["level"], "trace")
				}
			} else {
				assert.Empty(t, buffer.Bytes())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

	}
}

func mockGormConnectionWithLogger(t *testing.T, l Logger) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	mockedDB, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{Logger: NewGormLogger(l)},
	)
	if err != nil {
		t.Error(err)
	}

	return mockedDB, mock
}
