package resource

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsSpaceNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		error
		want bool
	}{
		{"std/errors error", errors.New("is not"), false},
		{"unwrapped CloudFoundry error", CloudFoundryError{
			Code: 40004,
		}, true},
		{"std wrapped CloudFoundry error", fmt.Errorf("%w", CloudFoundryError{
			Code: 40004,
		}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSpaceNotFoundError(tt.error)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
