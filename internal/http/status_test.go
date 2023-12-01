package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	require.True(t, IsStatusIn(http.StatusOK, http.StatusAccepted, http.StatusOK))
	require.True(t, IsStatusIn(http.StatusAccepted, http.StatusAccepted, http.StatusOK))
	require.False(t, IsStatusIn(http.StatusNoContent, http.StatusAccepted, http.StatusOK))

	require.True(t, IsResponseRedirect(http.StatusMovedPermanently))
	require.False(t, IsResponseRedirect(http.StatusNoContent))

	require.True(t, IsStatusSuccess(http.StatusNoContent))
	require.True(t, IsStatusSuccess(http.StatusTemporaryRedirect))
	require.False(t, IsStatusSuccess(http.StatusInternalServerError))

}
