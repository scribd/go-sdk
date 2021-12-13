package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	http2 "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"

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

		_, err = db.First(&TestRecord{}).Rows()
		require.Nil(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(testingHandlerBody))
		require.Nil(t, err)
	})
}

func TestNewDatabaseLoggingMiddleware(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	dbFile := "/tmp/test_db"
	defer os.Remove(dbFile)

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}
	defer db.Close()

	var (
		testRecord    = TestRecord{ID: 1, Name: "test_name"}
		updatedRecord = TestRecord{ID: 1, Name: "new_test_name"}
	)

	errors := db.Begin().
		CreateTable(TestRecord{}).
		Create(testRecord).
		Save(updatedRecord).
		Commit().GetErrors()

	for _, err := range errors {
		t.Fatalf("Errors: %v", err)
	}

	handler := testingDbHandler(t)

	var buffer bytes.Buffer
	config := &sdklogger.Config{
		ConsoleEnabled:    true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      "debug",
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

	var wg sync.WaitGroup

	// test concurrent calls to the handler
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			recorder := httptest.NewRecorder()

			req, reqErr := http.NewRequest("GET", "http://example.com", nil)
			require.Nil(t, reqErr)

			http2.WrapHandler(loggingHandler, "test", "test").ServeHTTP(recorder, req)

		}()
	}
	wg.Wait()

	// read first log entry
	var fields map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = dec.Decode(&fields)
	require.Nil(t, err)

	var sql = (fields["sql"]).(map[string]interface{})
	assert.NotEmpty(t, sql)

	assert.NotEmpty(t, sql["duration"])
	assert.NotEmpty(t, sql["affected_rows"])
	assert.NotEmpty(t, sql["file_location"])
	assert.NotNil(t, sql["values"])
}
