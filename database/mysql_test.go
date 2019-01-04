package database

import (
	"testing"
)

func TestNewMySQLTracer(t *testing.T) {
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
			"TestNewMySQLTracerWithNoOption",
			args{
				options: nil,
			},
			"mysql",
			false,
			false,
		},
		{
			"TestNewMySQLTracerWithEnableIgnoreSelectColumns",
			args{
				options: []func(*Tracer){
					func(t *Tracer) {
						t.EnableIgnoreSelectColumns()
					},
				},
			},
			"mysql",
			false,
			true,
		},
		{
			"TestNewMySQLTracerWithEnableRawQuery",
			args{
				options: []func(*Tracer){
					func(t *Tracer) {
						t.EnableRawQuery()
					},
				},
			},
			"mysql",
			true,
			false,
		},
		{
			"TestNewMySQLTracerWithAllOptionsEnable",
			args{
				options: []func(*Tracer){
					func(t *Tracer) {
						t.EnableRawQuery()
						t.EnableIgnoreSelectColumns()
					},
				},
			},
			"mysql",
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newMySQLTracer(tt.args.options...)
			if got.GetDBType() != tt.dbType {
				t.Errorf("NewMySQLTracer() DBType = %v, want %v", got.GetDBType(), tt.dbType)
			}
			if got.IsIgnoreSelectColumnsEnable() != tt.isIgnoreSelectedColumnsEnable {
				t.Errorf("NewMySQLTracer() option ignoreSelectColumn = %v, want %v",
					got.IsIgnoreSelectColumnsEnable(), tt.isIgnoreSelectedColumnsEnable)
			}
			if got.IsRawQueryEnable() != tt.isRawQueryEnable {
				t.Errorf("NewMySQLTracer() option rawQuery = %v, want %v",
					got.IsRawQueryEnable(), tt.isRawQueryEnable)
			}
		})
	}
}
