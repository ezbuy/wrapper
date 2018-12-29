package mssql

import (
	"github.com/ezbuy/wrapper"
	"github.com/ezbuy/wrapper/trace/database"
)

func NewMsSQLTracerWithMoreInfo(ins string, user string) *database.Tracer {
	return database.NewCustmizedTracer(ins, user, "mssql")
}

func NewMsSQLTracerWrapper() wrapper.Wrapper {
	return database.NewCustmizedTracerWrapper(
		database.NewCustmizedTracer("", "", "mssql"), false,
	)
}
