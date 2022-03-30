package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// PostData is the POST data from the request
type PostData struct {
	store map[string]interface{}
}

// UnmarshalJSON implements the Unmarshaler interface
func (pd *PostData) UnmarshalJSON(d []byte) error {
	var tmp map[string]interface{}
	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}

	pd.store = tmp

	return nil
}

// Get returns the value for the requested key or nil if
// not found
func (pd *PostData) Get(key string) interface{} {
	if val, ok := pd.store[key]; !ok {
		return nil
	} else {
		return val
	}
}

// GetString returns the string value of the requested key or
// nil if not found
func (pd *PostData) GetString(key, def string) string {
	val := pd.Get(key)
	if v, ok := val.(string); ok {
		return v
	}
	return def
}

func (pd *PostData) GetUUID(key string) (uuid.UUID, error) {
	val := pd.GetString(key, "")
	if val == "" {
		return uuid.Nil, fmt.Errorf("%s field cannot be empty", key)
	}

	r, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return r, nil
}

// Err logs an error message and writes an error to the response
func Err(w http.ResponseWriter, code int, err error) {
	logger.Error(err)
	http.Error(w, http.StatusText(code), code)
}
