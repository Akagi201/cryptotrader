package alert_test

import (
	"testing"

	"github.com/Akagi201/cryptotrader/cmd/alert/alert"
	"github.com/Akagi201/cryptotrader/cmd/alert/context"
	"github.com/Akagi201/esalert/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestSearchTPL(t *testing.T) {
	y := []byte(`
interval: "* * * * *"
search_index: foo-{{.Name}}
search_type: bar-{{.Name}}
search: {
	"query": {
		"query_string": {
			"query":"baz-{{.Name}}"
		}
	}
}`)

	var a alert.Alert
	require.Nil(t, yaml.Unmarshal(y, &a))
	require.Nil(t, a.Init())
	require.NotNil(t, a.SearchIndexTPL)
	require.NotNil(t, a.SearchTypeTPL)
	require.NotNil(t, a.SearchTPL)

	c := context.Context{
		Name: "wat",
	}
	searchIndex, searchType, searchQuery, err := a.CreateSearch(c)
	require.Nil(t, err)
	assert.Equal(t, "foo-wat", searchIndex)
	assert.Equal(t, "bar-wat", searchType)
	expectedSearch := search.Dict{
		"query": search.Dict{
			"query_string": search.Dict{
				"query": "baz-wat",
			},
		},
	}
	assert.Equal(t, expectedSearch, searchQuery)
}
