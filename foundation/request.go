package foundation

import (
	"errors"
	"net/http"

	"github.com/google/jsonapi"
)

var (
	errEmptyRequest = errors.New("empty request")
)

// RequestBody is used by POST end points to convert the request body to a struct
func RequestBody(r *http.Request, resource interface{}) error {
	if r == nil {
		return errEmptyRequest
	}

	err := jsonapi.UnmarshalPayload(r.Body, resource)
	if err != nil {
		return err
	}

	return nil
}
