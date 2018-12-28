package mysql

import (
	"github.com/ezbuy/wrapper/trace/database"
)

func NewMySQLTracer(ins string, user string) *database.Tracer {
	return database.NewCustmizedTracer(ins, user, "mysql")
}
