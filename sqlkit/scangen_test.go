package sqlkit_test

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parser := sqlkit.NewParser()

	t.Run("successful parse, same pkg", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/qux/qux.go", "gen")
		require.NoError(t, err)
		require.NotNil(t, parsed)
	})

	t.Run("successful parse, different pkg", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/foo/bar/bar.go", "bar")
		require.NoError(t, err)
		require.NotNil(t, parsed)
	})

	t.Run("failed parse, non-scanner field", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/corge/garply/garply.go", "garply")
		require.Error(t, err)
		require.Nil(t, parsed)
	})

	t.Run("failed parse, missing embed", func(t *testing.T) {
		parsed, err := parser.Parse("testdata/foo/quux/quux.go", "Quux")
		require.Error(t, err)
		require.Nil(t, parsed)
	})
}

func TestGenerate(t *testing.T) {
	parser := sqlkit.NewParser()

	t.Run("successful generation", func(t *testing.T) {
		const (
			inFile     = "testdata/foo/bar/bar.go"
			outFile    = "testdata/foo/bar/bar_gen.go"
			goldenFile = "testdata/foo/bar/bar_gen.golden.go"
		)

		parsed, err := parser.Parse(inFile, "bar")
		require.NotNil(t, parsed)
		require.NoError(t, err)

		t.Cleanup(func() {
			err = os.Remove(outFile)
			require.NoError(t, err)
		})

		config := sqlkit.GenerateConfig{
			Dialect:    sqlkit.Postgres,
			TableName:  "foos",
			OutputFile: outFile,
			LocalPaths: "github.com/nickcorin/toolkit",
		}

		err = sqlkit.Generate(&config, parsed)
		require.NoError(t, err)

		data, err := os.ReadFile(outFile)
		require.NotNil(t, data)
		require.NoError(t, err)

		assertGolden(t, goldenFile, data)
	})
}

var update = flag.Bool("update", false, "Updates golden files")

func assertGolden(t *testing.T, goldenFile string, actual []byte) {
	t.Helper()

	if *update {
		err := os.WriteFile(goldenFile, actual, 0o644)
		require.NoError(t, err)
		return
	}

	_, err := os.Stat(goldenFile)
	require.NoError(t, err)

	expected, err := os.ReadFile(goldenFile)
	require.NoError(t, err)
	require.NotNil(t, actual)

	require.True(t, bytes.Equal(expected, actual))
}
