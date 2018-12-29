package mysql

import (
	"github.com/ezbuy/wrapper"
	"github.com/ezbuy/wrapper/trace/database"
)

func NewCustomizedMySQLTracer(options ...func(*database.Tracer)) *database.Tracer {
	return database.NewCustmizedTracer("mysql", options...)
}

func NewMySQLTracerWrapper() wrapper.Wrapper {
	return database.NewCustmizedTracerWrapper(
		database.NewCustmizedTracer("mysql"),
	)
}
