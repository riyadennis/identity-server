package foundation

import (
	"errors"
	"net/http"

	"github.com/google/jsonapi"
)

var (
	errEmptyRequest   = errors.New("empty request")
	errInvalidContent = errors.New("invalid content type")
)

// RequestBody is used by POST end points to convert the request body to a struct
func RequestBody(r *http.Request, resource interface{}) error {
	if r == nil {
		return errEmptyRequest
	}

	if r.Header.Get("content-type") != "application/json" {
		return errInvalidContent
	}

	err := jsonapi.UnmarshalPayload(r.Body, resource)
	if err != nil {
		return err
	}

	return nil
}
