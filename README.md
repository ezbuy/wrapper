Wrapper

[![CircleCI](https://circleci.com/gh/ezbuy/wrapper/tree/feature%2Fadd-trace.svg?style=svg)](https://circleci.com/gh/ezbuy/wrapper/tree/feature%2Fadd-trace)
[![codecov](https://codecov.io/gh/ezbuy/wrapper/branch/feature%2Fadd-trace/graph/badge.svg)](https://codecov.io/gh/ezbuy/wrapper)
[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/ezbuy/wrapper)
---

# What is Wrapper ?

Wrapper is a developer friendly toolkit, which will provide a set of wrapper to intergrate 3rd Open Source Project with ezbuy codebase easily and efficiently.

# Goals

* Make middleware intergration internally and standardly
* Build a proxy between 3rd and ezbuy codebase

# Feature Lists

## Database tracer

Now database tracer provides a [jaeger](https://github.com/uber/jaeger-client-go) tracer client to trace SQL query context . Here is some sample usages for different go database package user.

> You can find the full sample usage in test file as well.

### For [sqlx](https://github.com/jmoiron/sqlx) users

```go
wp:= database.NewDefaultMySQLTracerWrapper()
// or wp:= database.NewMySQLTracerWrapper(options...)
originExecContextFunc:= wrapper.ExecContextFunc(func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					db, err:= sqlx.Connect("mysql","test:test@(localhost:3306)/test")
					if err != nil { // handle error
					}
					return db.ExecContext(ctx, query, args...)
				})
res,err:=wp.WrapExecContext(ctx,originExecContextFunc,sql,args...)
if err != nil{
// handle error
}
// handle res
```

### For [redis-orm](https://github.com/ezbuy/redis-orm) users

redis-orm users will do exactly nothing with this incoming changes , we will inject it in the generate code, and also ,we will add a set of db functions with context.

All you need to do is just type `go get -u github.com/ezbuy/redis-orm`

### Other users

To speak more generally, database tracer Wrapper accept a Query/ExecContextFunc and return you the same Query/ExecContextFunc(with tracer internal).

So, all sql packages which provide the Query/ExecContextFunc can add the jaeger tracer within one simple function.

### Tracer Options

* Hide select columns: `t.EnableIgnoreSelectColumns()`
* Show real args instead of `?`: `t.EnableRawQuery()`
* More custmized options are welcome.
```go
// how to use options
options:= []func(t *database.Tracer){
    func(t *database.Tracer){
        t.EnableIgnoreSelectColumns()
    },
    func(t *database.Tracer){
        t.EnableRawQuery()
    },
    func (t *database.Tracer){
        t.AddQueryBuilder(func(query string,args ...interface{})string{
            // handle with query and args
            return query
            })
    }
}
wp:= database.NewMySQLTracerWrapper(options...)
```

# Contribution

Issues and PRs are welcome.

