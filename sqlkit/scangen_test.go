package sqlkit_test

import (
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parser := sqlkit.NewParser()

	t.Run("successful parse", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/foo/bar/bar.go", "bar")
		require.NoError(t, err)
		require.NotNil(t, parsed)
	})

	t.Run("failed parse, non-scanner field", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/corge/garply/garply.go", "garply")
		require.Nil(t, parsed)
		require.Error(t, err)
	})

	t.Run("failed parse, missing embed", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/foo/quux/quux.go", "quux")
		require.Nil(t, parsed)
		require.Error(t, err)
	})
}

func TestGenerate(t *testing.T) {
	parser := sqlkit.NewParser()

	t.Run("successful generation", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/foo/bar/bar.go", "bar")
		require.NotNil(t, parsed)
		require.NoError(t, err)

		config := sqlkit.GenerateConfig{
			Dialect:    sqlkit.Postgres,
			TableName:  "foos",
			OutputFile: "testdata/foo/bar/bar_gen.go",
			LocalPaths: "github.com/nickcorin/toolkit",
		}

		err = sqlkit.Generate(&config, parsed)
		require.NoError(t, err)
	})
}
