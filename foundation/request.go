package foundation

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func RequestBody(r *http.Request, resource interface{}) error {
	if r.Header.Get("content-type") != "application/json" {
		err := errors.New("invalid content type")
		logrus.Errorf("content type is not json :: %v", err)
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
