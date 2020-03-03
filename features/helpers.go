package features

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)
// response is the response we get back
// from rest call to login endpoint
type response struct {
	// Status is the http status code
	Status int `json:"status"`
	// Message is the message like
	// `welcome John Doe`
	Message string `json:"message"`
	// ErrorCode helps to debug issues
	// will be empty on success requests.
	ErrorCode string `json:"error-code"`
}


// httpResponse will submit and http request using
// http client and then unmarshal the response into
// a struct.
func httpResponse(req *http.Request) (*response, error) {
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	loginResp = &response{}
	err = json.Unmarshal(body, loginResp)
	if err != nil {
		return nil, err
	}
	return loginResp, nil
}

// loginInput returns json request
// body for login endpoint in bytes.
func loginInput(email, password string) []byte {
	return []byte(`{
	"email": "` + email + `",
	"password": "` + password + `"
}`)
}

func registerInput(email, fname,lname string) []byte{
	return []byte(`{
	"first_name": "`+fname+`",
	"last_name": "`+lname+`",
	"email": "`+email+`",
	"terms": true
}`)
}
