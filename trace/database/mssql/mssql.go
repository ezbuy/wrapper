package mssql

import (
	"github.com/ezbuy/wrapper"
	"github.com/ezbuy/wrapper/trace/database"
)

func NewMsSQLTracerWithOptions(options ...func(*database.Tracer)) *database.Tracer {
	return database.NewCustmizedTracer("mssql", options...)
}

func NewMsSQLTracerWrapper() wrapper.Wrapper {
	return database.NewCustmizedTracerWrapper(
		database.NewCustmizedTracer("mssql"),
	)
}
