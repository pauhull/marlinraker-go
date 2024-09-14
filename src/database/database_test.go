package database

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marlinraker/src/files"
)

func TestInit(t *testing.T) {

	files.Fs = afero.NewMemMapFs()

	err := Init()
	require.NoError(t, err)

	content, err := afero.ReadFile(files.Fs, "db.json")
	require.NoError(t, err)
	assert.Equal(t, "{}", string(content))

	err = afero.WriteFile(files.Fs, "db.json", []byte("{invalid:json}"), 0755)
	require.NoError(t, err)

	err = Init()
	require.Error(t, err, "malformed db.json")
}

func TestPostItem(t *testing.T) {

	files.Fs = afero.NewMemMapFs()

	err := Init()
	require.NoError(t, err)

	val, err := PostItem("test", "key", "value", false)
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	val, err = PostItem("test", "number", 123, false)
	require.NoError(t, err)
	assert.Equal(t, 123, val)

	val, err = PostItem("test", "bool", true, false)
	require.NoError(t, err)
	assert.Equal(t, true, val)

	val, err = PostItem("test", "slice", []string{"a", "b", "c"}, false)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, val)

	val, err = PostItem("marlinraker", "foo", "bar", true)
	require.NoError(t, err)
	assert.Equal(t, "bar", val)

	val, err = PostItem("marlinraker", "boo", "far", false)
	require.Error(t, err, "reserved namespace access not allowed")
	assert.Nil(t, val)

	content, err := afero.ReadFile(files.Fs, "db.json")
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"test": {
			"key": "value",
			"number": 123,
			"bool": true,
			"slice": ["a", "b", "c"]
		},
		"marlinraker": {
			"foo": "bar"
		}
	}`, string(content))
}

func TestGetItem(t *testing.T) {

	files.Fs = afero.NewMemMapFs()

	content := `{"test":{"key":"value","number":123,"bool":true,"slice":["a","b","c"]},"marlinraker":{"foo":"bar"}}`
	err := afero.WriteFile(files.Fs, "db.json", []byte(content), 0755)
	require.NoError(t, err)

	err = Init()
	require.NoError(t, err)

	val, err := GetItem("test", "key", false)
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	val, err = GetItem("test", "number", false)
	require.NoError(t, err)
	assert.InDelta(t, 123., val, 1e-5)

	val, err = GetItem("test", "bool", false)
	require.NoError(t, err)
	assert.Equal(t, true, val)

	val, err = GetItem("test", "slice", false)
	require.NoError(t, err)
	assert.Equal(t, []any{"a", "b", "c"}, val)

	val, err = GetItem("test", "", false)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"bool":   true,
		"key":    "value",
		"number": 123.,
		"slice":  []any{"a", "b", "c"},
	}, val)

	val, err = GetItem("test", "foo", false)
	require.Error(t, err, `failed to get item: key "foo" in namespace "test" not found`)
	assert.Nil(t, val)

	val, err = GetItem("foo", "", false)
	require.Error(t, err, `failed to get item: namespace "foo" not found`)
	assert.Nil(t, val)

	val, err = GetItem("marlinraker", "foo", false)
	require.Error(t, err, "failed to get item: reserved namespace access not allowed")
	assert.Nil(t, val)

	val, err = GetItem("marlinraker", "foo", true)
	require.NoError(t, err)
	assert.Equal(t, "bar", val)
}

func TestDeleteItem(t *testing.T) {

	files.Fs = afero.NewMemMapFs()

	{
		content := `{"test":{"key":"value","number":123,"bool":true,"slice":["a","b","c"]},"foo":{"bar":0}}`
		err := afero.WriteFile(files.Fs, "db.json", []byte(content), 0755)
		require.NoError(t, err)
	}

	err := Init()
	require.NoError(t, err)

	val, err := DeleteItem("test", "key", false)
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	val, err = DeleteItem("test", "number", false)
	require.NoError(t, err)
	assert.InDelta(t, 123., val, 1e-5)

	val, err = DeleteItem("test", "foo", false)
	require.Error(t, err, `key "foo" in namespace "test" not found`)
	assert.Nil(t, val)

	val, err = DeleteItem("foo", "bar", false)
	require.NoError(t, err)
	assert.InDelta(t, 0., val, 1e-5)

	val, err = DeleteItem("marlinraker", "foo", false)
	require.Error(t, err, "reserved namespace access not allowed")
	assert.Nil(t, val)

	{
		content, err := afero.ReadFile(files.Fs, "db.json")
		require.NoError(t, err)
		assert.JSONEq(t, `{"test":{"bool":true,"slice":["a","b","c"]}}`, string(content))
	}
}

func TestListNamespaces(t *testing.T) {

	files.Fs = afero.NewMemMapFs()

	content := `{"test":{"key":"value","number":123,"bool":true,"slice":["a","b","c"]},"foo":{"bar":0}}`
	err := afero.WriteFile(files.Fs, "db.json", []byte(content), 0755)
	require.NoError(t, err)

	err = Init()
	require.NoError(t, err)

	namespaces := ListNamespaces()
	assert.ElementsMatch(t, append(ReservedNamespaces, "test", "foo"), namespaces)
}
