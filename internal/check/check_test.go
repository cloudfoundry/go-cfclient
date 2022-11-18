package check_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/check"
	"github.com/stretchr/testify/require"
	"testing"
)

type Y struct{}

func TestIsPointer(t *testing.T) {
	var s Y
	var p *Y

	require.False(t, check.IsPointer(nil))
	require.False(t, check.IsPointer(s))
	require.True(t, check.IsPointer(&p))
	require.True(t, check.IsPointer(p))
}

func TestIsNil(t *testing.T) {
	var s Y
	var p *Y

	require.False(t, check.IsNil(s))
	require.False(t, check.IsNil(&p))
	require.True(t, check.IsNil(p))
	require.True(t, check.IsNil(nil))

}

func TestIsNilAndIsPointer(t *testing.T) {
	var s Y
	var p *Y

	require.True(t, !check.IsNil(s) && !check.IsPointer(s))
	require.False(t, !check.IsNil(p) && !check.IsPointer(p))
	require.False(t, !check.IsNil(nil) && !check.IsPointer(nil))
}
