package util

import "encoding/json"

type HTTPError struct {
	Error string `json:"error"`
}

func GetHTTPError(s string) ([]byte, error) {
	h := &HTTPError{
		Error: s,
	}
	return json.Marshal(h)
}
