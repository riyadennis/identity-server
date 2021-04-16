package foundation

import (
	"errors"
	"net/http"

	"github.com/google/jsonapi"
)

// RequestBody is used by POST end points to convert the request body to a struct
func RequestBody(r *http.Request, resource interface{}) error {
	if r.Header.Get("content-type") != "application/json" {
		err := errors.New("invalid content type")
		return err
	}

	err := jsonapi.UnmarshalPayload(r.Body, resource)
	if err != nil {
		return err
	}

	return nil
}
