package util

import "encoding/json"

// HTTPError struct for error serializable HTTP response
type HTTPError struct {
	Error string `json:"error"`
}

// IDs  struct with string array for serializable HTTP response
type IDs struct {
	IDs []string `json:"IDs"`
}

// GetHTTPError returns marshaled json with an error string
func GetHTTPError(s string) ([]byte, error) {
	h := &HTTPError{
		Error: s,
	}
	return json.Marshal(h)
}

// GetIDs returns a marshaled json from a string list
func GetIDs(s []string) ([]byte, error) {
	i := &IDs{
		IDs: s,
	}
	return json.Marshal(i)
}
