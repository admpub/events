package meta

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	meta := Map{}

	meta.Set("test", "value")

	if value, exists := meta["test"]; exists {
		if s, ok := value.(string); ok {
			assert.Equal(t, s, "value", "they should be equal")
		} else {
			t.Log("test value not a string")
			t.Fail()
		}
	} else {
		t.Log("test value did not set")
		t.Fail()
	}
}

func TestFetch(t *testing.T) {
	meta := Map{}

	assert := assert.New(t)

	result, err := meta.Fetch("test")
	assert.NotNil(err)
	assert.Nil(result)
}

func TestFetchWithCallback(t *testing.T) {
	meta := Map{}

	assert := assert.New(t)

	var called bool
	result, err := meta.Fetch("test", func() (interface{}, error) {
		called = true
		return "value", nil
	})

	assert.Nil(err)
	assert.Equal(called, true, "callback has not been called")
	assert.Equal(result.(string), "value", "they should be equal")

	called = false
	result, err = meta.Fetch("test", func() (interface{}, error) {
		called = true
		return "value2", nil
	})

	assert.Nil(err)
	assert.NotEqual(called, true, "callback has been called")
	assert.NotEqual(result.(string), "value2", "they should not be equal")

	called = false
	result, err = meta.Fetch("test2", func(key string) (interface{}, error) {
		called = true
		return "value2", nil
	})

	assert.Nil(err)
	assert.Equal(called, true, "callback has not been called")
	assert.Equal(result.(string), "value2", "they should be equal")

	called = false
	result, err = meta.Fetch("test3", func(key string) (interface{}, error) {
		called = true
		return nil, errors.New("")
	})

	assert.Nil(result)
	assert.NotNil(err)
	assert.Equal(called, true, "callback has not been called")

	result, err = meta.Fetch("test3")
	assert.Nil(result)
	assert.NotNil(err)
}

func TestFetchWithOptionParameter(t *testing.T) {
	meta := Map{}

	assert := assert.New(t)

	result, err := meta.Fetch("test", "value")

	assert.Nil(err)
	assert.Equal(result.(string), "value", "they should be equal")

	result, err = meta.Fetch("test", "value2")

	assert.Nil(err)
	assert.NotEqual(result.(string), "value2", "they should be equal")
}

func TestCopy(t *testing.T) {
	assert := assert.New(t)

	original := Map{}
	original.Set("test", "value")

	copied := original.Copy()
	result, err := copied.Fetch("test")

	assert.Nil(err)
	assert.Equal(result.(string), "value", "they should be equal")

	original.Set("test", "value2")
	result, err = original.Fetch("test")

	assert.Nil(err)
	assert.Equal(result.(string), "value2", "they should be equal")
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)

	original := Map{}
	original.Set("first", "value")
	original.Set("second", "value2")

	other := Map{}
	other.Set("second", "value3")
	other.Set("third", "value4")

	merged := original.Merge(other)

	first, err := merged.Fetch("first")
	assert.Nil(err)
	assert.Equal(first.(string), "value", "they should be equal")
	second, err := merged.Fetch("second")
	assert.Nil(err)
	assert.Equal(second.(string), "value3", "they should be equal")
	third, err := merged.Fetch("third")
	assert.Nil(err)
	assert.Equal(third.(string), "value4", "they should be equal")
}
