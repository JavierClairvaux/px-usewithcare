package util

import "encoding/json"

type HTTPError struct {
	Error string `json:"error"`
}

type IDs struct {
	IDs []string `json:"IDs"`
}

func GetHTTPError(s string) ([]byte, error) {
	h := &HTTPError{
		Error: s,
	}
	return json.Marshal(h)
}

func GetIDs(s []string) ([]byte, error) {
	i := &IDs{
		IDs: s,
	}
	return json.Marshal(i)
}
