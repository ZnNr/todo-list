package taskerror

import (
	"encoding/json"
	"errors"
)

var (
	ErrRequireTitle = errors.New("require task title")
	ErrNotFoundTask = errors.New("not found task")
)

// MarshalError преобразует ошибку в формат JSON
func MarshalError(err error) []byte {
	type errJson struct {
		Error string `json:"error"`
	}
	res, _ := json.Marshal(errJson{Error: err.Error()})
	return res
}
