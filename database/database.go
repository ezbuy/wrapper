package database

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

var (
	// RawQueryOption enable the raw query option
	// raw query option will convert the ? placeHolders to real data.
	// Once enabled, "SELECT a FROM b WHERE c = ?" will be "SELECT a FROM b WHERE c = d"
	RawQueryOption = rawQueryOption{}
	// IgnoreSelectColumnsOption enable the ignore select columns option
	// ignore select column option will ignore the select columns
	// Once enabled, "SELECT a,b FROM c WHERE d = ?" will be "SELECT ... FROM c WHERE d = ?"
	IgnoreSelectColumnsOption = ignoreSelectColumnsOption{}
)

// tracer defines the database tracer
type tracer struct {
	instance      string
	statement     string
	dbtype        string
	user          string
	span          opentracing.Span
	queryBuilders []func(query string, args ...interface{}) string
}

// TracerOption defines the wrapper's option
type TracerOption interface {
	QueryBuilder() func(query string, args ...interface{}) string
}

// do gets the opentracing's global tracer ,and add span tags
// The tags in `Do` will follow the [opentracing spec](https://github.com/opentracing/specification/blob/master/semantic_conventions.md#span-tags-table)
func (t *tracer) do(ctx context.Context) {
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
func (t *tracer) close() {
	if t.span != nil {
		t.span.Finish()
	}
}

type rawQueryOption struct{}

func (opt rawQueryOption) QueryBuilder() func(query string, args ...interface{}) string {
	return rawQueryBuilder
}

type ignoreSelectColumnsOption struct{}

func (opt ignoreSelectColumnsOption) QueryBuilder() func(query string, args ...interface{}) string {
	return ignoreSelectColumnQueryBuilder
}

// addQueryBuilder adds the new builder(fn) to query builders.
// The builders in query builders will be execed in `Do` function
func (t *tracer) addQueryBuilder(fn func(query string, args ...interface{}) string) {
	t.queryBuilders = append(t.queryBuilders, fn)
}

// newTracer new a customized tracer with options
func newTracer(dbType string, options ...TracerOption) *tracer {
	t := &tracer{
		dbtype: dbType,
	}
	for _, op := range options {
		t.addQueryBuilder(op.QueryBuilder())
	}
	return t
}

// newTracerWithIgnoreColumnsOption new a default tracer with
// * ignore select columns option
func newTracerWithIgnoreColumnsOption(dbType string) *tracer {
	return newTracer(dbType, IgnoreSelectColumnsOption)
}

// newTracerWrapper new a default tracer wrapper with a tracer
func newTracerWrapper(t *tracer) *TracerWrapper {
	return &TracerWrapper{
		tracer: t,
	}
}

// newTracerWrapperWithTracer new a customized tracer wrapper with tracer and tracer options
func newTracerWrapperWithTracer(t *tracer) *TracerWrapper {
	return newTracerWrapper(t)
}

// NewTracerWrapper new a default tracer wrapper with
// * ignore select columns option
func NewTracerWrapper(dbType string) *TracerWrapper {
	return newTracerWrapper(newTracerWithIgnoreColumnsOption(dbType))
}

// TracerWrapper defines a tracer wrapper
// which impls WrapQueryContext and WrapExecContext
type TracerWrapper struct {
	tracer *tracer
}

// WrapQueryContext impls wrapper's WrapQueryContext
func (t *TracerWrapper) WrapQueryContext(fn QueryContextFunc, query string, args ...interface{}) QueryContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
		t.tracer.statement = t.hackQueryBuilder(query, args...)
		t.tracer.do(ctx)
		defer t.tracer.close()
		return fn(ctx, query, args...)
	}
	return tracerFn
}

// WrapExecContext impls wrapper's WrapExecContext
func (t *TracerWrapper) WrapExecContext(fn ExecContextFunc, query string, args ...interface{}) ExecContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		t.tracer.statement = t.hackQueryBuilder(query, args...)
		t.tracer.do(ctx)
		defer t.tracer.close()
		return fn(ctx, query, args...)
	}
	return tracerFn
}

// hackQueryBuilder exec all registered query builder
func (t *TracerWrapper) hackQueryBuilder(query string, args ...interface{}) string {
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
	query = strings.Replace(query, "\n", "", -1)
	query = strings.Replace(query, "select", "SELECT", -1)
	query = strings.Replace(query, "from", "FROM", -1)
	r := regexp.MustCompile("SELECT (.*) FROM")
	return r.ReplaceAllString(query, "SELECT ... FROM")
}
