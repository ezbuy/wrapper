package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func init() {
	opentracing.SetGlobalTracer(mocktracer.New())
}
func TestDefaultTracerWrapper_hackQueryBuilder(t *testing.T) {
	type fields struct {
		tracer *Tracer
	}
	type args struct {
		query string
		args  []interface{}
	}
	tracerAllOptions := NewCustmizedTracer("mysql",
		func(t *Tracer) {
			t.UseIgnoreSelectColumnsOption()
			t.UseRawQueryOption()
		})
	tracerWithRawQueryOption := NewCustmizedTracer("mysql",
		func(t *Tracer) {
			t.UseRawQueryOption()
		})
	tracerWithIgnoreSelectOption := NewCustmizedTracer("mysql",
		func(t *Tracer) {
			t.UseIgnoreSelectColumnsOption()
		})
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "TestHackQueryBuilderWithInternalOptions",
			fields: fields{
				tracer: tracerAllOptions,
			},
			args: args{
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			want: `SELECT ... FROM b WHERE c = d`,
		},
		{
			name: "TestHackQueryBuilderWithRawQueryOptions",
			fields: fields{
				tracer: tracerWithRawQueryOption,
			},
			args: args{
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			want: `SELECT a FROM b WHERE c = d`,
		},
		{
			name: "TestHackQueryBuilderWithIgnoreSelectOption",
			fields: fields{
				tracer: tracerWithIgnoreSelectOption,
			},
			args: args{
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			want: `SELECT ... FROM b WHERE c = ?`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := &TracerWrapper{
				tracer: tt.fields.tracer,
			}
			if got := dt.hackQueryBuilder(tt.args.query, tt.args.args...); got != tt.want {
				t.Errorf("DefaultTracerWrapper.hackQueryBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTracer_Do(t *testing.T) {

	type args struct {
		ctx context.Context
	}

	tracerAllOptions := NewCustmizedTracer("mysql",
		func(t *Tracer) {
			t.UseIgnoreSelectColumnsOption()
			t.UseRawQueryOption()
			t.statement = "SELECT ... FROM b WHERE c = d"
		})

	tests := []struct {
		name   string
		fields *Tracer
		args   args
	}{
		{
			name:   "TestDo",
			fields: tracerAllOptions,
			args: args{
				ctx: context.TODO(),
			},
		},
		{
			name:   "TestDo_With_Existed_Context",
			fields: tracerAllOptions,
			args: args{
				ctx: opentracing.ContextWithSpan(context.TODO(),
					opentracing.GlobalTracer().StartSpan("x").SetBaggageItem("a", "b")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := tt.fields
			dt.do(tt.args.ctx)
			if ins := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBInstance)); ins != dt.instance {
				t.Errorf("tags.DBInstance = %v,want %v", ins, dt.instance)
			}
			if st := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBStatement)); st != dt.statement {
				t.Errorf("tags.DBStatement= %v,want %v", st, dt.statement)
			}
			if tp := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBType)); tp != dt.dbtype {
				t.Errorf("tags.DBType= %v,want %v", tp, dt.dbtype)
			}
			if user := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBUser)); user != dt.user {
				t.Errorf("tags.DBUser= %v,want %v", user, dt.user)
			}
			dt.close()
		})
	}
}

func TestDefaultTracerWrapper_WrapQueryContext(t *testing.T) {
	type fields struct {
		tracer *Tracer
	}
	type args struct {
		ctx   context.Context
		fn    QueryContextFunc
		query string
		args  []interface{}
	}

	tests := []struct {
		name          string
		args          args
		wp            *TracerWrapper
		wantStatement string
	}{
		{
			name: "TestDefaultTracerWrapper_WrapQueryContext_MySQL",
			args: args{
				ctx: context.TODO(),
				fn: QueryContextFunc(func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.QueryContext(ctx, query, args...)
				}),
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			wp:            NewDefaultMySQLTracerWrapper(),
			wantStatement: "SELECT ... FROM b WHERE c = ?",
		},
		{
			name: "TestDefaultTracerWrapper_WrapQueryContext_MySQL_Custmized",
			args: args{
				ctx: context.TODO(),
				fn: QueryContextFunc(func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.QueryContext(ctx, query, args...)
				}),
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			wp: NewMySQLTracerWrapper(func(t *Tracer) {
				t.UseIgnoreSelectColumnsOption()
				t.UseRawQueryOption()
			}),
			wantStatement: "SELECT ... FROM b WHERE c = d",
		},
		{
			name: "TestDefaultTracerWrapper_WrapQueryContext_MsSQL",
			args: args{
				ctx: context.TODO(),
				fn: QueryContextFunc(func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.QueryContext(ctx, query, args...)
				}),
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			wp:            NewDefaultMsSQLTracerWrapper(),
			wantStatement: "SELECT ... FROM b WHERE c = ?",
		},
		{
			name: "TestDefaultTracerWrapper_WrapQueryContext_MsSQL_Custmized",
			args: args{
				ctx: context.TODO(),
				fn: QueryContextFunc(func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.QueryContext(ctx, query, args...)
				}),
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			wp: NewMsSQLTracerWrapper(func(t *Tracer) {
				t.UseIgnoreSelectColumnsOption()
				t.UseRawQueryOption()
			}),
			wantStatement: "SELECT ... FROM b WHERE c = d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wp.WrapQueryContext(tt.args.fn, tt.args.query, tt.args.args...)(tt.args.ctx, tt.args.query, tt.args.args...)
			dt := tt.wp.tracer
			if ins := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBInstance)); ins != dt.instance {
				t.Errorf("tags.DBInstance = %v,want %v", ins, dt.instance)
			}
			if st := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBStatement)); st != tt.wantStatement {
				t.Errorf("tags.DBStatement= %v,want %v", st, tt.wantStatement)
			}
			if tp := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBType)); tp != dt.dbtype {
				t.Errorf("tags.DBType= %v,want %v", tp, dt.dbtype)
			}
			if user := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBUser)); user != dt.user {
				t.Errorf("tags.DBUser= %v,want %v", user, dt.user)
			}
			dt.close()
		})
	}
}

func TestDefaultTracerWrapper_WrapExecContext(t *testing.T) {
	type fields struct {
		tracer *Tracer
	}
	type args struct {
		ctx   context.Context
		fn    ExecContextFunc
		query string
		args  []interface{}
	}

	tests := []struct {
		name          string
		args          args
		wp            *TracerWrapper
		wantStatement string
	}{
		{
			name: "TestDefaultTracerWrapper_WrapExecContext",
			args: args{
				ctx: context.TODO(),
				fn: ExecContextFunc(func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.ExecContext(ctx, query, args...)
				}),
				query: "UPDATE a SET c = d WHERE c = ?",
				args:  []interface{}{"e"},
			},
			wp:            NewDefaultMySQLTracerWrapper(),
			wantStatement: "UPDATE a SET c = d WHERE c = ?",
		},
		{
			name: "TestDefaultTracerWrapper_WrapExecContext_MySQL_Custmized",
			args: args{
				ctx: context.TODO(),
				fn: ExecContextFunc(func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.ExecContext(ctx, query, args...)
				}),
				query: "UPDATE a SET c = d WHERE c = ?",
				args:  []interface{}{"e"},
			},
			wp: NewMySQLTracerWrapper(func(t *Tracer) {
				t.UseIgnoreSelectColumnsOption()
				t.UseRawQueryOption()
			}),
			wantStatement: "UPDATE a SET c = d WHERE c = e",
		},
		{
			name: "TestDefaultTracerWrapper_WrapExecContext_MsSQL",
			args: args{
				ctx: context.TODO(),
				fn: ExecContextFunc(func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.ExecContext(ctx, query, args...)
				}),
				query: "UPDATE a SET c = d WHERE c = ?",
				args:  []interface{}{"e"},
			},
			wp:            NewDefaultMsSQLTracerWrapper(),
			wantStatement: "UPDATE a SET c = d WHERE c = ?",
		},
		{
			name: "TestDefaultTracerWrapper_WrapExecContext_MsSQL_Custmized",
			args: args{
				ctx: context.TODO(),
				fn: ExecContextFunc(func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.ExecContext(ctx, query, args...)
				}),
				query: "UPDATE a SET c = d WHERE c = ?",
				args:  []interface{}{"e"},
			},
			wp: NewMsSQLTracerWrapper(func(t *Tracer) {
				t.UseIgnoreSelectColumnsOption()
				t.UseRawQueryOption()
			}),
			wantStatement: "SELECT ... FROM b WHERE c = e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wp.WrapExecContext(tt.args.fn, tt.args.query, tt.args.args...)(tt.args.ctx, tt.args.query, tt.args.args...)
			dt := tt.wp.tracer
			if ins := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBInstance)); ins != dt.instance {
				t.Errorf("tags.DBInstance = %v,want %v", ins, dt.instance)
			}
			if st := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBStatement)); st != dt.statement {
				t.Errorf("tags.DBStatement= %v,want %v", st, dt.statement)
			}
			if tp := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBType)); tp != dt.dbtype {
				t.Errorf("tags.DBType= %v,want %v", tp, dt.dbtype)
			}
			if user := dt.span.(*mocktracer.MockSpan).Tag(string(tags.DBUser)); user != dt.user {
				t.Errorf("tags.DBUser= %v,want %v", user, dt.user)
			}
			dt.close()
		})
	}
}
