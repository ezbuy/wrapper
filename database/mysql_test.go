package database

import (
	"testing"
)

func TestNewMySQLTracer(t *testing.T) {
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
				options: []TracerOption{IgnoreSelectColumnsOption},
			},
			"mysql",
			false,
			true,
		},
		{
			"TestNewMySQLTracerWithEnableRawQuery",
			args{
				options: []TracerOption{RawQueryOption},
			},
			"mysql",
			true,
			false,
		},
		{
			"TestNewMySQLTracerWithAllOptionsEnable",
			args{
				options: []TracerOption{IgnoreSelectColumnsOption, RawQueryOption},
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
		})
	}
}
