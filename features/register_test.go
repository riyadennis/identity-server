package features

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/riyadennis/identity-server/internal/store"
	"io"
	"net/http"
)
var resisterResponse *response

func RegisterFeatureContext(s *godog.Suite) {
	s.BeforeScenario(beforeScenario)
	s.Step(`^firstName "([^"]*)"$`, firstName)
	s.Step(`^lastName "([^"]*)"$`, lastName)
	s.Step(`^that user register$`, thatUserRegister)
	s.Step(`^errorCode "([^"]*)"$`, errorCode)
	s.AfterScenario(afterScenario)
}

func firstName(firstName string) error {
	user = &store.User{
		FirstName:        firstName,
		Email:           userEmail,
		Terms: false,
	}
	return nil
}

func lastName(lastName string) error {
	user.LastName = lastName
	return nil
}

func thatUserRegister() error {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/register", HOST),
		registerInput(user))
	req.Header.Set("Content-Type", "application/json")
	if err != nil{
		return err
	}
	resisterResponse, err = httpResponse(req)
	if err != nil{
		return err
	}
	fmt.Println(resisterResponse)
	return nil
}

func registerInput(u *store.User) io.Reader {
 return bytes.NewReader( []byte(`{
	"first_name": "` + u.FirstName + `",
	"last_name": "` + u.LastName + `",
	"email": "` + u.Email + `",
	"terms": true
}`))
}

func errorCode(errorCode string) error {
	if resisterResponse.ErrorCode != errorCode{
		return errors.New("invalid response")
	}
	return nil
}



