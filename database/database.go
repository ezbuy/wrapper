package database

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/ezbuy/wrapper"
	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

// Tracer defines the database tracer
type Tracer struct {
	instance              string
	statement             string
	dbtype                string
	user                  string
	span                  opentracing.Span
	isRawQueryEnable      bool
	isIgnoreSelectColumns bool
	queryBuilders         []func(query string, args ...interface{}) string
}

// do gets the opentracing's global tracer ,and add span tags
// The tags in `Do` will follow the [opentracing spec](https://github.com/opentracing/specification/blob/master/semantic_conventions.md#span-tags-table)
func (t *Tracer) do(ctx context.Context) {
	tracer := opentracing.GlobalTracer()
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = tracer.StartSpan(t.dbtype)
	} else {
		span = tracer.StartSpan(t.dbtype, opentracing.ChildOf(span.Context()))
	}
	tags.DBInstance.Set(span, t.instance)
	tags.DBStatement.Set(span, t.statement)
	tags.DBType.Set(span, t.dbtype)
	tags.DBUser.Set(span, t.user)
	ctx = opentracing.ContextWithSpan(ctx, span)
	t.span = span
}

// close closes the opentracing's span
func (t *Tracer) close() {
	if t.span != nil {
		t.span.Finish()
	}
}

// EnableRawQuery enable the raw query option
// raw query option will convert the ? placeHolders to real data.
// Once enabled, "SELECT a FROM b WHERE c = ?" will be "SELECT a FROM b WHERE c = d"
func (t *Tracer) EnableRawQuery() {
	t.isRawQueryEnable = true
}

// EnableIgnoreSelectColumns enable the ignore select columns option
// ignore select column option will ignore the select columns
// Once enabled, "SELECT a,b FROM c WHERE d = ?" will be "SELECT ... FROM c WHERE d = ?"
func (t *Tracer) EnableIgnoreSelectColumns() {
	t.isIgnoreSelectColumns = true
}

// IsRawQueryEnable checks if raw query option is enable
func (t *Tracer) IsRawQueryEnable() bool {
	return t.isRawQueryEnable
}

// IsIgnoreSelectColumnsEnable checks if ignore select columns option is enable
func (t *Tracer) IsIgnoreSelectColumnsEnable() bool {
	return t.isIgnoreSelectColumns
}

// GetDBType returns the set db type
func (t *Tracer) GetDBType() string {
	return t.dbtype
}

// AddQueryBuilders adds the new builder(fn) to query builders.
// The builders in query builders will be execed in `Do` function
func (t *Tracer) AddQueryBuilders(fn func(query string, args ...interface{}) string) {
	t.queryBuilders = append(t.queryBuilders, fn)
}

// NewCustmizedTracer new a customized tracer with options
func NewCustmizedTracer(dbType string, options ...func(t *Tracer)) *Tracer {
	t := &Tracer{
		dbtype: dbType,
	}
	for _, op := range options {
		op(t)
	}
	return t
}

// newDefaultTracer new a default tracer with
// * ignore select columns option
func newDefaultTracer(dbType string) *Tracer {
	t := &Tracer{
		dbtype:                dbType,
		isRawQueryEnable:      false,
		isIgnoreSelectColumns: true,
	}
	return t
}

// newTracerWrapper new a default tracer wrapper with a tracer
func newTracerWrapper(t *Tracer) *DefaultTracerWrapper {
	return &DefaultTracerWrapper{
		tracer: t,
	}
}

// NewCustmizedTracerWrapper new a customized tracer wrapper with tracer and tracer options
func NewCustmizedTracerWrapper(t *Tracer) *DefaultTracerWrapper {
	return newTracerWrapper(t)
}

// NewDefaultTracerWrapper new a default tracer wrapper with
// * ignore select columns option
func NewDefaultTracerWrapper(dbType string) *DefaultTracerWrapper {
	return newTracerWrapper(newDefaultTracer(dbType))
}

// DefaultTracerWrapper defines a tracer wrapper
// which impls WrapQueryContext and WrapExecContext
type DefaultTracerWrapper struct {
	tracer *Tracer
}

// WrapQueryContext impls wrapper's WrapQueryContext
func (t *DefaultTracerWrapper) WrapQueryContext(ctx context.Context, fn wrapper.QueryContextFunc,
	query string, args ...interface{}) wrapper.QueryContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
		t.tracer.statement = t.hackQueryBuilder(query, args...)
		t.tracer.do(ctx)
		defer t.tracer.close()
		return fn(ctx, query, args...)
	}
	return tracerFn
}

// WrapExecContext impls wrapper's WrapExecContext
func (t *DefaultTracerWrapper) WrapExecContext(ctx context.Context, fn wrapper.ExecContextFunc,
	query string, args ...interface{}) wrapper.ExecContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		t.tracer.statement = t.hackQueryBuilder(query, args...)
		t.tracer.do(ctx)
		defer t.tracer.close()
		return fn(ctx, query, args...)
	}
	return tracerFn
}

// hackQueryBuilder exec all registered query builder
func (t *DefaultTracerWrapper) hackQueryBuilder(query string, args ...interface{}) string {
	if t.tracer.IsRawQueryEnable() {
		t.tracer.AddQueryBuilders(rawQueryBuilder)
	}
	if t.tracer.IsIgnoreSelectColumnsEnable() {
		t.tracer.AddQueryBuilders(ignoreSelectColumnQueryBuilder)
	}
	for _, fn := range t.tracer.queryBuilders {
		query = fn(query, args...)
	}
	return query
}

func rawQueryBuilder(query string, args ...interface{}) string {
	q := strings.Replace(query, "?", "%v", -1)
	return fmt.Sprintf(q, args...)
}

func ignoreSelectColumnQueryBuilder(query string, args ...interface{}) string {
	query = strings.Replace(query, "select", "SELECT", -1)
	query = strings.Replace(query, "from", "FROM", -1)
	r := regexp.MustCompile("SELECT (.*) FROM")
	return r.ReplaceAllString(query, "SELECT ... FROM")
}
