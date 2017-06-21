package meta

import (
	"fmt"
)

type Map map[string]interface{}

func (source Map) Copy() Map {
	result := Map{}
	for key, value := range source {
		result[key] = value
	}
	return result
}

func (target Map) Merge(source Map) Map {
	result := Map{}
	for key, value := range target {
		result[key] = value
	}
	for key, value := range source {
		result[key] = value
	}
	return result
}

func (source Map) Set(key string, value interface{}) {
	source[key] = value
}

func (source Map) Remove(key string) (value interface{}) {
	value = source[key]
	delete(source, key)
	return
}

func (source Map) Fetch(key string, options ...interface{}) (interface{}, error) {
	if value, exists := source[key]; exists {
		return value, nil
	}
	if len(options) > 0 {
		for _, option := range options {
			switch option := option.(type) {
			case func(key string) (interface{}, error):
				newValue, err := option(key)
				if err == nil {
					source[key] = newValue
					return newValue, nil
				}
			case func() (interface{}, error):
				newValue, err := option()
				if err == nil {
					source[key] = newValue
					return newValue, nil
				}
			default:
				if option != nil {
					source[key] = option
					return option, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Key %s not found in hash", key)
}
