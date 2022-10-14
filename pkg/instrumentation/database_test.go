package instrumentation

import (
	"context"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
			name: "SELECT",
			tags: map[string]interface{}{
				"db.count":      int64(-1),
				"db.err":        nil,
				"db.table":      "",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      `SELECT COUNT(*) FROM SQLITE_MASTER WHERE TYPE='TABLE' AND NAME=?`,
				"resource.name": `SELECT COUNT(*) FROM SQLITE_MASTER WHERE TYPE='TABLE' AND NAME=?`,
			},
		},
		{
			name: "INSERT",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        nil,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      "INSERT INTO `TEST_RECORDS` (`NAME`,`ID`) VALUES (?,?) RETURNING `ID`",
				"resource.name": "INSERT INTO `TEST_RECORDS` (`NAME`,`ID`) VALUES (?,?) RETURNING `ID`",
			},
		},
		{
			name: "UPDATE",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        nil,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      "UPDATE `TEST_RECORDS` SET `ID`=?,`NAME`=? WHERE `ID` = ? RETURNING `ID`",
				"resource.name": "UPDATE `TEST_RECORDS` SET `ID`=?,`NAME`=? WHERE `ID` = ? RETURNING `ID`",
			},
		},
		{
			name: "SELECT",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        nil,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      "SELECT * FROM `TEST_RECORDS` WHERE `ID` = ? ORDER BY `TEST_RECORDS`.`ID` LIMIT 1",
				"resource.name": "SELECT * FROM `TEST_RECORDS` WHERE `ID` = ? ORDER BY `TEST_RECORDS`.`ID` LIMIT 1",
			},
		},
		{
			name: "DELETE",
			tags: map[string]interface{}{
				"db.count":      int64(1),
				"db.err":        nil,
				"db.table":      "test_records",
				"service.name":  "test_app_name-mysql",
				"span.type":     "sql",
				"db.query":      "DELETE FROM `TEST_RECORDS` WHERE `ID` = ? AND `TEST_RECORDS`.`ID` = ? RETURNING `ID`",
				"resource.name": "DELETE FROM `TEST_RECORDS` WHERE `ID` = ? AND `TEST_RECORDS`.`ID` = ? RETURNING `ID`",
			},
		},
	}

	mt := mocktracer.Start()
	defer mt.Stop()

	dbFile := path.Join(t.TempDir(), "test_db")
	db, err := gorm.Open(sqlite.Open(dbFile))
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}

	InstrumentDatabase(db, "test_app_name")
	db = TraceDatabase(context.Background(), db)

	var (
		testRecord    = TestRecord{ID: 1, Name: "test_name"}
		updatedRecord = TestRecord{ID: 1, Name: "new_test_name"}
		readRecord    = &TestRecord{}
	)

	if err = db.AutoMigrate(TestRecord{}); err != nil {
		t.Fatalf("Failed to migrate DB: %s", err)
	}

	err = db.Begin().Create(&testRecord).
		Save(&updatedRecord).
		First(&readRecord).
		Delete(&testRecord).
		Commit().Error
	if err != nil {
		t.Fatalf("Failed to commit changes on test record: %s", err)
	}

	spans := mt.FinishedSpans()
	if len(spans) != 5 {
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
