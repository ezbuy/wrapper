package database

import (
	"testing"
)

func TestNewMsSQLTracer(t *testing.T) {
	type args struct {
		options []func(*Tracer)
	}
	tests := []struct {
		name                          string
		args                          args
		dbType                        string
		isRawQueryEnable              bool
		isIgnoreSelectedColumnsEnable bool
	}{
		{
			"TestNewMsSQLTracerWithNoOption",
			args{
				options: nil,
			},
			"mssql",
			false,
			false,
		},
		{
			"TestNewMsSQLTracerWithEnableIgnoreSelectColumns",
			args{
				options: []func(*Tracer){
					func(t *Tracer) {
						t.UseIgnoreSelectColumnsOption()
					},
				},
			},
			"mssql",
			false,
			true,
		},
		{
			"TestNewMsSQLTracerWithEnableRawQuery",
			args{
				options: []func(*Tracer){
					func(t *Tracer) {
						t.UseRawQueryOption()
					},
				},
			},
			"mssql",
			true,
			false,
		},
		{
			"TestNewMsSQLTracerWithAllOptionsEnable",
			args{
				options: []func(*Tracer){
					func(t *Tracer) {
						t.UseIgnoreSelectColumnsOption()
						t.UseRawQueryOption()
					},
				},
			},
			"mssql",
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newMsSQLTracer(tt.args.options...)
			if got.GetDBType() != tt.dbType {
				t.Errorf("NewMsSQLTracer() DBType = %v, want %v", got.GetDBType(), tt.dbType)
			}
		})
	}
}
