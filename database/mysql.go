package database

import (
	"github.com/ezbuy/wrapper"
)

func newMySQLTracer(options ...func(*Tracer)) *Tracer {
	return NewCustmizedTracer("mysql", options...)
}

func NewMySQLTracerWrapper(options ...func(*Tracer)) wrapper.Wrapper {
	return NewCustmizedTracerWrapper(
		newMySQLTracer(options...),
	)
}

func NewDefaultMySQLTracerWrapper() wrapper.Wrapper {
	return NewDefaultTracerWrapper("mysql")
}
