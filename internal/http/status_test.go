package http

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestStatus(t *testing.T) {
	require.True(t, IsStatusIn(http.StatusOK, http.StatusAccepted, http.StatusOK))
	require.True(t, IsStatusIn(http.StatusAccepted, http.StatusAccepted, http.StatusOK))
	require.False(t, IsStatusIn(http.StatusNoContent, http.StatusAccepted, http.StatusOK))
}
