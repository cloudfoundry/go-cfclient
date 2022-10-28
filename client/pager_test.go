package client

import (
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/stretchr/testify/require"
)

func TestQuerystringReader(t *testing.T) {
	u := "https://api.example.org/v3/apps?page=1&per_page=50"
	reader, err := newQuerystringReader(u)
	require.NoError(t, err)
	require.Equal(t, 1, reader.Int("page"))
	require.Equal(t, 50, reader.Int("per_page"))

	u = "https://api.example.org/v3/apps"
	reader, err = newQuerystringReader(u)
	require.NoError(t, err)
	require.Equal(t, 0, reader.Int("page"))

	u = "https://api.example.org/v3/apps?order_by=id"
	reader, err = newQuerystringReader(u)
	require.NoError(t, err)
	require.Equal(t, "id", reader.String("order_by"))

	_, err = newQuerystringReader("")
	require.Error(t, err)
}

func TestPager(t *testing.T) {
	paginationPage1 := resource.Pagination{
		TotalResults: 120,
		TotalPages:   3,
		First: resource.Link{
			Href: "https://api.example.org/v3/apps?page=1&per_page=50",
		},
		Last: resource.Link{
			Href: "https://api.example.org/v3/apps?page=3&per_page=50",
		},
		Next: resource.Link{
			Href: "https://api.example.org/v3/apps?page=2&per_page=50",
		},
		Previous: resource.Link{},
	}
	paginationPage2 := resource.Pagination{
		TotalResults: 120,
		TotalPages:   3,
		First: resource.Link{
			Href: "https://api.example.org/v3/apps?page=1&per_page=50",
		},
		Last: resource.Link{
			Href: "https://api.example.org/v3/apps?page=3&per_page=50",
		},
		Next: resource.Link{
			Href: "https://api.example.org/v3/apps?page=3&per_page=50",
		},
		Previous: resource.Link{
			Href: "https://api.example.org/v3/apps?page=1&per_page=50",
		},
	}
	paginationPage3 := resource.Pagination{
		TotalResults: 120,
		TotalPages:   3,
		First: resource.Link{
			Href: "https://api.example.org/v3/apps?page=1&per_page=50",
		},
		Last: resource.Link{
			Href: "https://api.example.org/v3/apps?page=3&per_page=50",
		},
		Next: resource.Link{},
		Previous: resource.Link{
			Href: "https://api.example.org/v3/apps?page=2&per_page=50",
		},
	}
	listOpts := NewAppListOptions()

	// Defaults
	require.Equal(t, 1, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)
	require.Equal(t, "", listOpts.OrderBy)

	// First page
	pager := NewPager(paginationPage1)
	require.True(t, pager.HasNextPage())
	require.False(t, pager.HasPreviousPage())
	pager.NextPage(listOpts)
	require.Equal(t, 2, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)

	// Second page
	pager = NewPager(paginationPage2)
	require.True(t, pager.HasNextPage())
	require.True(t, pager.HasPreviousPage())
	pager.NextPage(listOpts)
	require.Equal(t, 3, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)

	// Third page
	pager = NewPager(paginationPage3)
	require.False(t, pager.HasNextPage())
	require.True(t, pager.HasPreviousPage())
	pager.NextPage(listOpts)
	require.Equal(t, 3, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)

	pager.PreviousPage(listOpts)
	require.Equal(t, 2, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)

	// Second page
	pager = NewPager(paginationPage2)
	require.True(t, pager.HasNextPage())
	require.True(t, pager.HasPreviousPage())
	pager.PreviousPage(listOpts)
	require.Equal(t, 1, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)

	// First page
	pager = NewPager(paginationPage1)
	require.True(t, pager.HasNextPage())
	require.False(t, pager.HasPreviousPage())
	pager.PreviousPage(listOpts)
	require.Equal(t, 1, listOpts.Page)
	require.Equal(t, 50, listOpts.PerPage)
}
