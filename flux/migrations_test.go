package flux_test

import "embed"

//go:embed testdata/migrations/postgres/*.sql
var pgMigrations embed.FS
