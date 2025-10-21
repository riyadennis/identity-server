package foundation

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errEmptyRequest       = errors.New("empty request")
	errMissingContentType = errors.New("missing content type")
)

// RequestBody is used by POST end points to convert the request body to a struct
func RequestBody(r *http.Request, resource interface{}) error {
	if r == nil {
		return errEmptyRequest
	}

	if r.Header.Get("content-type") == "" {
		return errMissingContentType
	}

	err := json.NewDecoder(r.Body).Decode(resource)
	if err != nil {
		return err
	}

	return nil
}
