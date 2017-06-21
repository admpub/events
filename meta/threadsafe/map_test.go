package threadsafe

import (
	"github.com/stretchr/testify/assert"

	"errors"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	meta := New()

	meta.Set("test", "value")

	if value, exists := meta.values["test"]; exists {
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
	meta := New()

	assert := assert.New(t)

	result, err := meta.Fetch("test")
	assert.NotNil(err)
	assert.Nil(result)

	start := time.Now()
	group := new(sync.WaitGroup)
	group.Add(7)
	for i, key := range []string{"key1", "key1", "key2", "key2", "key3", "key3", "key4"} {
		go func(i int, key string) {
			defer group.Done()
			result, err := meta.Fetch(key, func(key string) (interface{}, error) {
				time.Sleep(1 * time.Second)
				return key, nil
			})

			assert.Nil(err)
			assert.NotNil(result)
		}(i, key)
	}
	group.Wait()
	if time.Since(start) < 4*time.Second {
		t.Fail()
	}
}

func TestFetchWithCallback(t *testing.T) {
	meta := New()

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

	calls := make(map[int]struct{})
	group := new(sync.WaitGroup)
	count := 4
	group.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer group.Done()
			result, err := meta.Fetch("test4", func() (interface{}, error) {
				calls[i] = struct{}{}
				return i, nil
			})

			assert.Nil(err)
			assert.NotNil(result)
		}(i)
	}
	group.Wait()
	assert.Equal(len(calls), 1, "callback has been called multiple times")
}

func TestFetchWithOptionParameter(t *testing.T) {
	meta := New()

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

	original := New()
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

	original := New()
	original.Set("first", "value")
	original.Set("second", "value2")

	other := New()
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
