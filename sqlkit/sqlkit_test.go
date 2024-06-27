package sqlkit_test

import (
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestFlags(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		f := make(sqlkit.Flags)

		f.Set("foo", "bar")
		f.Set("baz", "qux")

		encoded := f.Encode()
		require.Equal(t, "baz=qux&foo=bar", encoded)
	})

	t.Run("Decode", func(t *testing.T) {
		f := make(sqlkit.Flags)

		err := f.Decode("baz=qux&foo=bar")
		require.NoError(t, err)

		foo := f.Get("foo")
		baz := f.Get("baz")

		require.Equal(t, "bar", foo)
		require.Equal(t, "qux", baz)
	})

	t.Run("Set", func(t *testing.T) {
		f := make(sqlkit.Flags)
		f.Set("foo", "bar")

		val, ok := f["foo"]
		require.True(t, ok)
		require.Equal(t, []string{"bar"}, val)
	})
}
