package database

import (
	"testing"
)

func TestNewMsSQLTracer(t *testing.T) {
	type args struct {
		options []TracerOption
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
				options: []TracerOption{IgnoreSelectColumnsOption},
			},
			"mssql",
			false,
			true,
		},
		{
			"TestNewMsSQLTracerWithEnableRawQuery",
			args{
				options: []TracerOption{RawQueryOption},
			},
			"mssql",
			true,
			false,
		},
		{
			"TestNewMsSQLTracerWithAllOptionsEnable",
			args{
				options: []TracerOption{IgnoreSelectColumnsOption, RawQueryOption},
			},
			"mssql",
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newMsSQLTracer(tt.args.options...)
			if got.dbtype != tt.dbType {
				t.Errorf("NewMsSQLTracer() DBType = %v, want %v", got.dbtype, tt.dbType)
			}
		})
	}
}
