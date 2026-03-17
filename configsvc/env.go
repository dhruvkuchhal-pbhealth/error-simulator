package configsvc

import (
	"encoding/json"
	"errors"
)

// Expand returns a string representation of the provided value.
// It accepts either a string or a map[string]string. For a map,
// it returns the JSON encoding. For unsupported types, it returns an error.
func Expand(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case map[string]string:
		b, err := json.Marshal(val)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case nil:
		return "", errors.New("nil value provided to Expand")
	default:
		return "", errors.New("unsupported type provided to Expand")
	}
}
