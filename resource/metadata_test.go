package resource_test

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetadata(t *testing.T) {
	type metaTest struct {
		prefix string
		key    string
		value  string
	}
	tests := []metaTest{
		{
			prefix: "pre",
			key:    "key",
			value:  "val",
		},
		{
			prefix: "",
			key:    "empty-pre-key",
			value:  "val",
		},
		{
			prefix: "cf.example.org",
			key:    "key",
			value:  "val",
		},
		{
			prefix: "pre",
			key:    "no-val-key",
			value:  "",
		},
		{
			prefix: "",
			key:    "only-key",
			value:  "",
		},
	}

	for _, tt := range tests {
		k := fmt.Sprintf("%s/%s", tt.prefix, tt.key)
		if tt.prefix == "" {
			k = tt.key
		}

		// add some annotations and labels
		m := resource.Metadata{}
		m.SetAnnotation(tt.prefix, tt.key, tt.value)
		m.SetLabel(tt.prefix, tt.key, tt.value)
		require.Equal(t, tt.value, *m.Annotations[k], "key: %s", k)
		require.Equal(t, tt.value, *m.Labels[k], "key: %s", k)

		// remove them
		m.RemoveAnnotation(tt.prefix, tt.key)
		m.RemoveLabel(tt.prefix, tt.key)
		require.Nil(t, m.Annotations[k], "key: %s", k)
		require.Nil(t, m.Labels[k], "key: %s", k)

		// new annotations and labels
		m = resource.Metadata{}
		m.SetAnnotation(tt.prefix, tt.key, tt.value)
		m.SetLabel(tt.prefix, tt.key, tt.value)
		require.Equal(t, tt.value, *m.Annotations[k], "key: %s", k)
		require.Equal(t, tt.value, *m.Labels[k], "key: %s", k)

		// clear
		m.Clear()
		require.Nil(t, m.Annotations[k], "key: %s", k)
		require.Nil(t, m.Labels[k], "key: %s", k)
	}
}
