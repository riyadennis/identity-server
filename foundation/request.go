package foundation

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestBody is used by POST end points to convert the request body to a struct
func RequestBody(r *http.Request, resource interface{}) error {
	if r.Header.Get("content-type") != "application/json" {
		err := errors.New("invalid content type")
		return err
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if err == io.EOF {
			err := errors.New("empty request body")
			return err
		}
		return err
	}

	err = json.Unmarshal(data, resource)
	if err != nil {
		return err
	}

	return nil
}
