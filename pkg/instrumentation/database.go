package instrumentation

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	ddtrace "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gorm.io/gorm"
)

const (
	// ParentSpanGormKey is the name of the parent span key
	ParentSpanGormKey = "tracingParentSpan"
	// SpanGormKey is the name of the span key
	SpanGormKey = "tracingSpan"
)

// TraceDatabase sets span to gorm settings, returns cloned DB
func TraceDatabase(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}
	parentSpan, _ := ddtrace.SpanFromContext(ctx)
	return db.Set(ParentSpanGormKey, parentSpan)
}

// InstrumentDatabase adds callbacks for tracing, call TraceDatabase to make it work
func InstrumentDatabase(db *gorm.DB, appName string) {
	callbacks := newCallbacks(appName)

	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row", callbacks)
}

type callbacks struct {
	serviceName string
}

func newCallbacks(appName string) *callbacks {
	return &callbacks{
		serviceName: fmt.Sprintf("%s-%s", appName, "mysql"),
	}
}

func (c *callbacks) beforeCreate(db *gorm.DB) { c.before(db, "INSERT", c.serviceName) }
func (c *callbacks) afterCreate(db *gorm.DB)  { c.after(db) }
func (c *callbacks) beforeQuery(db *gorm.DB)  { c.before(db, "SELECT", c.serviceName) }
func (c *callbacks) afterQuery(db *gorm.DB)   { c.after(db) }
func (c *callbacks) beforeUpdate(db *gorm.DB) { c.before(db, "UPDATE", c.serviceName) }
func (c *callbacks) afterUpdate(db *gorm.DB)  { c.after(db) }
func (c *callbacks) beforeDelete(db *gorm.DB) { c.before(db, "DELETE", c.serviceName) }
func (c *callbacks) afterDelete(db *gorm.DB)  { c.after(db) }
func (c *callbacks) beforeRow(db *gorm.DB)    { c.before(db, "", c.serviceName) }
func (c *callbacks) afterRow(db *gorm.DB)     { c.after(db) }
func (c *callbacks) before(db *gorm.DB, operationName string, serviceName string) {
	val, ok := db.Get(ParentSpanGormKey)
	if !ok {
		return
	}

	parentSpan := val.(ddtrace.Span)
	spanOpts := []ddtrace.StartSpanOption{
		ddtrace.ChildOf(parentSpan.Context()),
		ddtrace.SpanType(ext.SpanTypeSQL),
		ddtrace.ServiceName(serviceName),
	}
	if operationName == "" {
		operationName = strings.Split(db.Statement.SQL.String(), " ")[0]
	}
	sp := ddtrace.StartSpan(operationName, spanOpts...)
	db.Set(SpanGormKey, sp)
}

func (c *callbacks) after(db *gorm.DB) {
	val, ok := db.Get(SpanGormKey)
	if !ok {
		return
	}

	sp := val.(ddtrace.Span)
	sp.SetTag(ext.ResourceName, strings.ToUpper(db.Statement.SQL.String()))
	sp.SetTag("db.table", db.Statement.Table)
	sp.SetTag("db.query", strings.ToUpper(db.Statement.SQL.String()))
	sp.SetTag("db.err", db.Error)
	sp.SetTag("db.count", db.RowsAffected)
	sp.Finish()
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	var err error

	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	// gorm does some magic, if you pass CallbackProcessor here - nothing works
	switch name {
	case "create":
		err = db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		if err != nil {
			return
		}
		err = db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
		if err != nil {
			return
		}
	case "query":
		err = db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		if err != nil {
			return
		}
		err = db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
		if err != nil {
			return
		}
	case "update":
		err = db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		if err != nil {
			return
		}
		err = db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
		if err != nil {
			return
		}
	case "delete":
		err = db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		if err != nil {
			return
		}
		err = db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
		if err != nil {
			return
		}
	case "row":
		err = db.Callback().Row().Before(gormCallbackName).Register(beforeName, c.beforeRow)
		if err != nil {
			return
		}
		err = db.Callback().Row().After(gormCallbackName).Register(afterName, c.afterRow)
		if err != nil {
			return
		}
	}
}
