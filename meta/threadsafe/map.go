package threadsafe

import (
	"fmt"
	"sync"
)

func New() Map {
	return Map{new(sync.RWMutex), make(map[string]interface{})}
}

type Map struct {
	*sync.RWMutex
	values map[string]interface{}
}

func (source Map) Copy() Map {
	source.RLock()
	defer source.RUnlock()

	result := New()
	for key, value := range source.values {
		result.values[key] = value
	}

	return result
}

func (target Map) Merge(source Map) Map {
	result := New()

	target.RLock()
	for key, value := range target.values {
		result.values[key] = value
	}
	target.RUnlock()

	source.RLock()
	for key, value := range source.values {
		result.values[key] = value
	}
	source.RUnlock()

	return result
}

func (source *Map) Set(key string, value interface{}) {
	source.Lock()
	defer source.Unlock()

	source.values[key] = value
}

func (source *Map) Remove(key string) (value interface{}) {
	source.Lock()
	defer source.Unlock()
	value = source.values[key]
	delete(source.values, key)
	return
}

func (source *Map) Fetch(key string, options ...interface{}) (interface{}, error) {
	source.RLock()
	value, exists := source.values[key]
	if exists {
		source.RUnlock()
		return value, nil
	}
	source.RUnlock()
	source.Lock()
	defer source.Unlock()
	// check value after lock again, may be another thread set it while we were waiting for the lock
	if value, exists := source.values[key]; exists {
		return value, nil
	}

	var err error
	if len(options) > 0 {
		for _, option := range options {
			switch option := option.(type) {
			case func(key string) (interface{}, error):
				value, err = option(key)
				if err == nil {
					source.values[key] = value
					return value, nil
				}
			case func() (interface{}, error):
				value, err = option()
				if err == nil {
					source.values[key] = value
					return value, nil
				}
			default:
				if option != nil {
					source.values[key] = option
					return option, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Key %s not found in hash - %s", key, err)
}

func (source Map) Must(key string, options ...interface{}) interface{} {
	val, err := source.Fetch(key, options...)
	if err != nil {
		panic(err)
	}
	return val
}

func (source Map) String(key string, options ...interface{}) (string, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return ``, err
	}
	v, _ := val.(string)
	return v, nil
}

func (source Map) Int(key string, options ...interface{}) (int, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(int)
	return v, nil
}

func (source Map) Uint(key string, options ...interface{}) (uint, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(uint)
	return v, nil
}

func (source Map) Int8(key string, options ...interface{}) (int8, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(int8)
	return v, nil
}

func (source Map) Uint8(key string, options ...interface{}) (uint8, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(uint8)
	return v, nil
}

func (source Map) Int16(key string, options ...interface{}) (int16, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(int16)
	return v, nil
}

func (source Map) Uint16(key string, options ...interface{}) (uint16, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(uint16)
	return v, nil
}

func (source Map) Int32(key string, options ...interface{}) (int32, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(int32)
	return v, nil
}

func (source Map) Uint32(key string, options ...interface{}) (uint32, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(uint32)
	return v, nil
}

func (source Map) Int64(key string, options ...interface{}) (int64, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(int64)
	return v, nil
}

func (source Map) Uint64(key string, options ...interface{}) (uint64, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(uint64)
	return v, nil
}

func (source Map) Float32(key string, options ...interface{}) (float32, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(float32)
	return v, nil
}

func (source Map) Float64(key string, options ...interface{}) (float64, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(float64)
	return v, nil
}

func (source Map) Map(key string, options ...interface{}) (map[string]interface{}, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.(map[string]interface{})
	return v, nil
}

func (source Map) Slice(key string, options ...interface{}) ([]interface{}, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]interface{})
	return v, nil
}

func (source Map) StringSlice(key string, options ...interface{}) ([]string, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]string)
	return v, nil
}

func (source Map) IntSlice(key string, options ...interface{}) ([]int, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]int)
	return v, nil
}

func (source Map) UintSlice(key string, options ...interface{}) ([]uint, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]uint)
	return v, nil
}

func (source Map) Int8Slice(key string, options ...interface{}) ([]int8, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]int8)
	return v, nil
}

func (source Map) Uint8Slice(key string, options ...interface{}) ([]uint8, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]uint8)
	return v, nil
}

func (source Map) Int16Slice(key string, options ...interface{}) ([]int16, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]int16)
	return v, nil
}

func (source Map) Uint16Slice(key string, options ...interface{}) ([]uint16, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]uint16)
	return v, nil
}

func (source Map) Int32Slice(key string, options ...interface{}) ([]int32, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]int32)
	return v, nil
}

func (source Map) Uint32Slice(key string, options ...interface{}) ([]uint32, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]uint32)
	return v, nil
}

func (source Map) Int64Slice(key string, options ...interface{}) ([]int64, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]int64)
	return v, nil
}

func (source Map) Uint64Slice(key string, options ...interface{}) (uint64, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return 0, err
	}
	v, _ := val.(uint64)
	return v, nil
}

func (source Map) Float32Slice(key string, options ...interface{}) ([]float32, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]float32)
	return v, nil
}

func (source Map) Float64Slice(key string, options ...interface{}) ([]float64, error) {
	val, err := source.Fetch(key, options...)
	if err != nil {
		return nil, err
	}
	v, _ := val.([]float64)
	return v, nil
}
