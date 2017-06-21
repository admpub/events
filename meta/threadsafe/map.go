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

	if value, exists := source.values[key]; exists {
		source.RUnlock()
		return value, nil
	} else {
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
}
