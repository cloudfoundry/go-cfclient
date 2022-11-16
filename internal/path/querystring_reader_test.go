package path

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuerystringReader(t *testing.T) {
	u := "https://api.example.org/v3/apps?page=1&per_page=50"
	reader, err := NewQuerystringReader(u)
	require.NoError(t, err)
	require.Equal(t, 1, reader.Int("page"))
	require.Equal(t, 50, reader.Int("per_page"))

	u = "https://api.example.org/v3/apps"
	reader, err = NewQuerystringReader(u)
	require.NoError(t, err)
	require.Equal(t, 0, reader.Int("page"))

	u = "https://api.example.org/v3/apps?order_by=id"
	reader, err = NewQuerystringReader(u)
	require.NoError(t, err)
	require.Equal(t, "id", reader.String("order_by"))

	_, err = NewQuerystringReader("")
	require.Error(t, err)
}
