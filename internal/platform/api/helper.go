package api

import (
	"encoding/json"
	"log"
	"net/http"
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

func Err(w http.ResponseWriter, code int, err error) {
	log.Printf("error: %v", err)
	http.Error(w, http.StatusText(code), code)
}
