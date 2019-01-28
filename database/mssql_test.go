package database

import (
	"testing"
)

func TestNewMsSQLTracer(t *testing.T) {
	type args struct {
		options []WrapperOption
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
				options: []WrapperOption{IgnoreSelectColumnsOption},
			},
			"mssql",
			false,
			true,
		},
		{
			"TestNewMsSQLTracerWithEnableRawQuery",
			args{
				options: []WrapperOption{RawQueryOption},
			},
			"mssql",
			true,
			false,
		},
		{
			"TestNewMsSQLTracerWithAllOptionsEnable",
			args{
				options: []WrapperOption{IgnoreSelectColumnsOption, RawQueryOption},
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
