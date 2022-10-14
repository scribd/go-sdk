package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	http2 "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	sdkdatabasecontext "github.com/scribd/go-sdk/pkg/context/database"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

type TestRecord struct {
	ID   int
	Name string
}

func testingDbHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		db, err := sdkdatabasecontext.Extract(req.Context())
		require.Nil(t, err)

		db.Exec("insert asd")
		err = db.First(&TestRecord{}).Error
		require.Nil(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(testingHandlerBody))
		require.Nil(t, err)
	})
}

func TestNewDatabaseLoggingMiddleware(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	dbFile := path.Join(t.TempDir(), "test_db")
	db, err := gorm.Open(sqlite.Open(dbFile))
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}

	testRecord := TestRecord{ID: 1, Name: "test_name"}
	updatedRecord := TestRecord{ID: 1, Name: "new_test_name"}

	if err = db.AutoMigrate(&TestRecord{}); err != nil {
		t.Fatalf("Failed to migrate DB: %s", err)
	}

	if err = db.Begin().
		Create(&testRecord).
		Save(&updatedRecord).
		Commit().Error; err != nil {
		t.Fatalf("Failed to create and update record: %s", err)
	}

	handler := testingDbHandler(t)

	var buffer bytes.Buffer
	config := &sdklogger.Config{
		ConsoleEnabled:    true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      "trace",
		FileEnabled:       false,
	}

	databaseLoggingMiddleware := NewDatabaseLoggingMiddleware()
	databaseLoggingHandler := databaseLoggingMiddleware.Handler(handler)

	databaseMiddleware := NewDatabaseMiddleware(db)
	databaseHandler := databaseMiddleware.Handler(databaseLoggingHandler)

	l, err := sdklogger.NewBuilder(config).BuildTestLogger(&buffer)
	require.Nil(t, err)

	loggingMiddleware := NewLoggingMiddleware(l)
	loggingHandler := loggingMiddleware.Handler(databaseHandler)

	testHandler := NewRequestIDMiddleware().Handler(loggingHandler)

	var wg sync.WaitGroup

	// test concurrent calls to the handler
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			recorder := httptest.NewRecorder()

			req, reqErr := http.NewRequest("GET", "http://example.com", nil)
			require.Nil(t, reqErr)

			http2.WrapHandler(testHandler, "test", "test").ServeHTTP(recorder, req)

		}()
	}
	wg.Wait()

	// read first log entry
	var fields map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = dec.Decode(&fields)
	require.Nil(t, err)

	dbFields, ok := (fields["sql"]).(map[string]interface{})
	assert.True(t, ok, "%s not found in log fields", "trace")
	assert.NotEmpty(t, dbFields)

	assert.Contains(t, dbFields, "duration")
	assert.Contains(t, dbFields, "affected_rows")
	assert.Contains(t, dbFields, "sql")
}
