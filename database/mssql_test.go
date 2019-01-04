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
						t.EnableIgnoreSelectColumns()
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
						t.EnableRawQuery()
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
						t.EnableRawQuery()
						t.EnableIgnoreSelectColumns()
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
			if got.IsIgnoreSelectColumnsEnable() != tt.isIgnoreSelectedColumnsEnable {
				t.Errorf("NewMsSQLTracer() option ignoreSelectColumn = %v, want %v",
					got.IsIgnoreSelectColumnsEnable(), tt.isIgnoreSelectedColumnsEnable)
			}
			if got.IsRawQueryEnable() != tt.isRawQueryEnable {
				t.Errorf("NewMsSQLTracer() option rawQuery = %v, want %v",
					got.IsRawQueryEnable(), tt.isRawQueryEnable)
			}
		})
	}
}
