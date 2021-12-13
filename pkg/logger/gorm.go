package logger

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

func NewGormLogger(l Logger) gormLogger {
	return gormLogger{
		logger: l,
	}
}

type gormLogger struct {
	logger Logger
}

// Print formats & prints the log
func (gl gormLogger) Print(values ...interface{}) {
	if values[0] == "sql" {
		location := values[1].(string)
		duration := values[2].(time.Duration)
		sql := values[3].(string)
		formattedValues := formatValues(values[4].([]interface{}))
		affectedRows := values[5].(int64)

		logger := gl.logger.WithFields(Fields{
			"sql": Fields{
				"duration":      duration,
				"affected_rows": affectedRows,
				"file_location": location,
				"values":        formattedValues,
			},
		})

		logger.Debugf(sql)
	} else {
		gl.logger.Tracef("%v", values[2:]...)
	}
}

// Code blatantly stolen and modified for our needs from Gorm:
// https://github.com/jinzhu/gorm/blob/2a3ab99/logger.go#L30-L109
func formatValues(values []interface{}) string {
	var formattedValues []string
	for _, value := range values {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() {
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok {
				if t.IsZero() {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
				}
			} else if b, ok := value.([]byte); ok {
				if str := string(b); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			} else if r, ok := value.(driver.Valuer); ok {
				if value, err := r.Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			} else {
				switch value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
					formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
				default:
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
		} else {
			formattedValues = append(formattedValues, "NULL")
		}
	}

	return strings.Join(formattedValues, " ")
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
