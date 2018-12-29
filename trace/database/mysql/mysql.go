package mysql

import (
	"github.com/ezbuy/wrapper"
	"github.com/ezbuy/wrapper/trace/database"
)

func NewCustomizedMySQLTracer(ins string, user string) *database.Tracer {
	return database.NewCustmizedTracer(ins, user, "mysql")
}

func NewMySQLTracerWrapper() wrapper.Wrapper {
	return database.NewCustmizedTracerWrapper(
		database.NewCustmizedTracer("", "", "mysql"), false,
	)
}
