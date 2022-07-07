package instrumentation

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

type TestRecord struct {
	ID   int
	Name string
}

func TestInstrumentDatabase(t *testing.T) {
	expectedSpans := []struct {
		name string
		tags map[string]interface{}
	}{
		{
			name: "INSERT",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        false,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      `INSERT  INTO "TEST_RECORDS" ("ID","NAME") VALUES (?,?)`,
				"resource.name": `INSERT  INTO "TEST_RECORDS" ("ID","NAME") VALUES (?,?)`,
			},
		},
		{
			name: "UPDATE",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        false,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      `UPDATE "TEST_RECORDS" SET "NAME" = ?  WHERE "TEST_RECORDS"."ID" = ?`,
				"resource.name": `UPDATE "TEST_RECORDS" SET "NAME" = ?  WHERE "TEST_RECORDS"."ID" = ?`,
			},
		},
		{
			name: "SELECT",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        false,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      `SELECT * FROM "TEST_RECORDS"   ORDER BY "TEST_RECORDS"."ID" ASC LIMIT 1`,
				"resource.name": `SELECT * FROM "TEST_RECORDS"   ORDER BY "TEST_RECORDS"."ID" ASC LIMIT 1`,
			},
		},
		{
			name: "DELETE",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        false,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      `DELETE FROM "TEST_RECORDS"  WHERE "TEST_RECORDS"."ID" = ?`,
				"resource.name": `DELETE FROM "TEST_RECORDS"  WHERE "TEST_RECORDS"."ID" = ?`,
			},
		},
	}

	mt := mocktracer.Start()
	defer mt.Stop()

	dbFile := "/var/tmp/test_db"
	defer os.Remove(dbFile)

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}
	defer db.Close()

	InstrumentDatabase(db, "test_app_name")
	db = TraceDatabase(context.Background(), db)

	var (
		testRecord    = TestRecord{ID: 1, Name: "test_name"}
		updatedRecord = TestRecord{ID: 1, Name: "new_test_name"}
		readRecord    = &TestRecord{}
	)

	errors := db.Begin().
		CreateTable(TestRecord{}).
		Create(testRecord).
		Save(updatedRecord).
		First(readRecord).
		Delete(testRecord).
		Commit().GetErrors()
	for _, err := range errors {
		t.Fatalf("Errors: %v", err)
	}

	spans := mt.FinishedSpans()
	if len(spans) != 4 {
		t.Fatalf("Unexpected number of spans: %d", len(spans))
	}

	for i := range spans {
		actualName := spans[i].OperationName()
		actualTags := spans[i].Tags()

		expectedName := expectedSpans[i].name
		expectedTags := expectedSpans[i].tags

		if actualName != expectedName {
			t.Errorf("Got span: %s, expected: %s", actualName, expectedName)
		}

		if diff := cmp.Diff(expectedTags, actualTags); diff != "" {
			t.Error(diff)
		}
	}
}
