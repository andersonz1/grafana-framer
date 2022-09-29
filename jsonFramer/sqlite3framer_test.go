package jsonFramer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/andersonz1/grafana-framer/jsonFramer"
)

func Test_queryJSONUsingSQLite3(t *testing.T) {
	tests := []struct {
		name         string
		jsonString   string
		query        string
		rootSelector string
		want         string
		wantErr      error
		test         func(t *testing.T, want string)
	}{
		{
			name:       "empty array should not throw error",
			jsonString: `[]`,
			query:      "select * from input",
			test: func(t *testing.T, got string) {
				require.Equal(t, "[]", got)
			},
		},
		{
			name:       "valid array should not throw error",
			jsonString: `[{ "name": "foo" },{ "name": "bar" }]`,
			query:      "select * from input",
			test: func(t *testing.T, got string) {
				require.Equal(t, "[\n  {\n    \"name\": \"foo\"\n  },\n  {\n    \"name\": \"bar\"\n  }\n]\n", got)
			},
		},
		{
			name:       "valid summarize should not throw error",
			jsonString: `[{ "name": "foo" },{ "name": "bar" }]`,
			query:      "select count(*) as 'count' from input",
			test: func(t *testing.T, got string) {
				require.Equal(t, "[\n  {\n    \"count\": 2\n  }\n]\n", got)
			},
		},
		{
			name:         "valid nested json document should not throw error",
			jsonString:   `{ "users" : [{ "name": "foo" },{ "name": "bar" }] }`,
			query:        "select * from input",
			rootSelector: ".users",
			test: func(t *testing.T, got string) {
				require.Equal(t, "[\n  {\n    \"name\": \"foo\"\n  },\n  {\n    \"name\": \"bar\"\n  }\n]\n", got)
			},
		},
		{
			name:         "valid nested with invalid root selector json document should not throw error",
			jsonString:   `{ "users" : [{ "name": "foo" },{ "name": "bar" }] }`,
			query:        "select * from input",
			rootSelector: "users",
			test: func(t *testing.T, got string) {
				require.Equal(t, "[\n  {\n    \"name\": \"foo\"\n  },\n  {\n    \"name\": \"bar\"\n  }\n]\n", got)
			},
		},
		{
			name:         "valid deep nested with invalid root selector json document should not throw error",
			jsonString:   `{ "foo" : { "users" : [{ "name": "foo" },{ "name": "bar" }] }}`,
			query:        "select * from input",
			rootSelector: "foo.users",
			test: func(t *testing.T, got string) {
				require.Equal(t, "[\n  {\n    \"name\": \"foo\"\n  },\n  {\n    \"name\": \"bar\"\n  }\n]\n", got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonFramer.QueryJSONUsingSQLite3(tt.jsonString, tt.query, tt.rootSelector)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
			}
			require.Nil(t, err)
			require.NotNil(t, got)
			if tt.test != nil {
				tt.test(t, got)
			}
		})
	}
}
