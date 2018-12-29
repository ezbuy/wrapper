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

type Tracer struct {
	instance  string
	statement string
	dbtype    string
	user      string
	span      opentracing.Span
}

func (t *Tracer) Do(ctx context.Context, statement string) {
	tracer := opentracing.GlobalTracer()
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	span = tracer.StartSpan(t.dbtype, opentracing.ChildOf(span.Context()))
	tags.DBInstance.Set(span, t.instance)
	tags.DBStatement.Set(span, statement)
	tags.DBType.Set(span, t.dbtype)
	tags.DBUser.Set(span, t.user)
	ctx = opentracing.ContextWithSpan(ctx, span)
	t.span = span
}

func (t *Tracer) Close() {
	if t.span != nil {
		t.span.Finish()
	}
}

func NewCustmizedTracer(ins, user, dbType string) *Tracer {
	return &Tracer{
		instance: ins,
		user:     user,
		dbtype:   dbType,
	}
}

func NewCustmizedTracerWrapper(t *Tracer, enableRawQuery bool) wrapper.Wrapper {
	return &DefaultTracer{
		isRawQueryEnable:      enableRawQuery,
		isIgnoreSelectColumns: true,
		tracer:                t,
	}
}

func NewDefaultTracerWrapper() wrapper.Wrapper {
	return &DefaultTracer{
		isRawQueryEnable:      false,
		isIgnoreSelectColumns: true,
		tracer:                NewDefaultTracer("", ""),
	}
}

func NewDefaultTracer(ins string, user string) *Tracer {
	return &Tracer{
		dbtype:   "default db",
		instance: ins,
		user:     user,
	}
}

type DefaultTracer struct {
	isRawQueryEnable      bool
	isIgnoreSelectColumns bool
	tracer                *Tracer
}

func (t *DefaultTracer) WrapQueryContext(ctx context.Context, fn wrapper.QueryContextFunc,
	query string, args ...interface{}) wrapper.QueryContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
		t.tracer.Do(ctx, t.hackQueryBuilder(query, args...))
		defer t.tracer.Close()
		return fn(ctx, query, args...)
	}
	return tracerFn
}

func (t *DefaultTracer) WrapExecContext(ctx context.Context, fn wrapper.ExecContextFunc,
	query string, args ...interface{}) wrapper.ExecContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		t.tracer.Do(ctx, t.hackQueryBuilder(query, args...))
		defer t.tracer.Close()
		return fn(ctx, query, args...)
	}
	return tracerFn
}

func (t *DefaultTracer) hackQueryBuilder(query string, args ...interface{}) string {
	if t.isRawQueryEnable {
		q := strings.Replace(query, "?", "%v", -1)
		query = fmt.Sprintf(q, args...)
	}
	if t.isIgnoreSelectColumns {
		query = strings.Replace(query, "select", "SELECT", -1)
		query = strings.Replace(query, "from", "FROM", -1)
		r := regexp.MustCompile("SELECT (.*) FROM")
		query = r.ReplaceAllString(query, "SELECT ... FROM")
	}
	return query
}