package sqlkit_test

import (
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("successful generation", func(t *testing.T) {
		p := sqlkit.NewParser()

		parsed, err := p.Parse("testdata/foo/bar/bar.go", "barType")
		require.NoError(t, err)

		config := sqlkit.GenerateConfig{
			Dialect:    sqlkit.Postgres,
			TableName:  "foos",
			OutputFile: "testdata/foo/bar/bar_gen.go",
			LocalPaths: "github.com/nickcorin/toolkit/sqlkit/testdata/foo",
		}

		err = sqlkit.Generate(&config, parsed)
		require.NoError(t, err)
	})
}
