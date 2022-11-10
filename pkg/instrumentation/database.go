package instrumentation

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	ddtrace "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct {
	serviceName string
}

func newCallbacks(appName string) *callbacks {
	return &callbacks{
		serviceName: fmt.Sprintf("%s-%s", appName, "app"),
	}
}

func (c *callbacks) beforeCreate(scope *gorm.Scope)   { c.before(scope, "INSERT", c.serviceName) }
func (c *callbacks) afterCreate(scope *gorm.Scope)    { c.after(scope) }
func (c *callbacks) beforeQuery(scope *gorm.Scope)    { c.before(scope, "SELECT", c.serviceName) }
func (c *callbacks) afterQuery(scope *gorm.Scope)     { c.after(scope) }
func (c *callbacks) beforeUpdate(scope *gorm.Scope)   { c.before(scope, "UPDATE", c.serviceName) }
func (c *callbacks) afterUpdate(scope *gorm.Scope)    { c.after(scope) }
func (c *callbacks) beforeDelete(scope *gorm.Scope)   { c.before(scope, "DELETE", c.serviceName) }
func (c *callbacks) afterDelete(scope *gorm.Scope)    { c.after(scope) }
func (c *callbacks) beforeRowQuery(scope *gorm.Scope) { c.before(scope, "", c.serviceName) }
func (c *callbacks) afterRowQuery(scope *gorm.Scope)  { c.after(scope) }
func (c *callbacks) before(scope *gorm.Scope, operationName string, serviceName string) {
	val, ok := scope.Get(ParentSpanGormKey)
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
		operationName = strings.Split(scope.SQL, " ")[0]
	}
	sp := ddtrace.StartSpan(operationName, spanOpts...)
	scope.Set(SpanGormKey, sp)
}

func (c *callbacks) after(scope *gorm.Scope) {
	val, ok := scope.Get(SpanGormKey)
	if !ok {
		return
	}
	sp := val.(ddtrace.Span)
	sp.SetTag(ext.ResourceName, strings.ToUpper(scope.SQL))
	sp.SetTag("db.table", scope.TableName())
	sp.SetTag("db.query", strings.ToUpper(scope.SQL))
	sp.SetTag("db.err", scope.HasError())
	sp.SetTag("db.count", scope.DB().RowsAffected)
	sp.Finish()
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	// gorm does some magic, if you pass CallbackProcessor here - nothing works
	switch name {
	case "create":
		db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		db.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		db.Callback().RowQuery().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}
