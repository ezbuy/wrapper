package mssql

import (
	"github.com/ezbuy/wrapper/trace/database"
)

func NewMsSQLTracer(ins string, user string) *database.Tracer {
	return database.NewCustmizedTracer(ins, user, "mssql")
}
