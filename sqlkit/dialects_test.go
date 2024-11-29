package sqlkit_test

import (
	"errors"
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
)

func TestGetConnector(t *testing.T) {
	tests := []struct {
		name    string
		dialect sqlkit.Dialect
		want    sqlkit.Connector
		err     error
	}{
		{
			name:    "unknown dialect",
			dialect: sqlkit.Unspecified,
			want:    nil,
			err:     sqlkit.ErrUnsupportedDialect,
		},
		{
			name:    "valid dialect",
			dialect: sqlkit.Postgres,
			want:    new(sqlkit.PostgresConnector),
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sqlkit.GetConnector(tt.dialect)
			if !errors.Is(err, tt.err) {
				t.Errorf("sqlkit.GetConnector(%q) error = %v, want %v", tt.dialect, err, tt.err)
			}

			if got != tt.want {
				t.Errorf("sqlkit.GetConnector(%q) = %v, want %v", tt.dialect, got, tt.want)
			}
		})
	}
}

func TestGetDialectFromString(t *testing.T) {
	tests := []struct {
		name    string
		dialect string
		want    sqlkit.Dialect
		valid   bool
	}{
		{
			name:    "empty string",
			dialect: "",
			want:    sqlkit.Unspecified,
			valid:   false,
		},
		{
			name:    "valid dialect",
			dialect: "postgres",
			want:    sqlkit.Postgres,
			valid:   true,
		},
		{
			name:    "valid dialect - uppercase",
			dialect: "POSTGRES",
			want:    sqlkit.Postgres,
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sqlkit.GetDialectFromString(tt.dialect)
			if got != tt.want {
				t.Errorf("sqlkit.GetDialectFromString(%q) = %v, want %v", tt.dialect, got, tt.want)
			}
		})
	}
}
