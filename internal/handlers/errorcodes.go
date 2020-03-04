package handlers

const (
	InvalidUserData    = "invalid-user-data"
	EmailAlreadyExists = "email-already-exists"
	DatabaseError      = "database-error"
	PassWordError      = "password-error"
	InvalidRequest     = "invalid-request"
	TokenError         = "token-error"
	UnAuthorised       = "unauthorised"
	UserDoNotExist     = "user-do-not-exist"
)

type CustomError struct {
	Code string
	Err  error
}

func (e *CustomError) Error() string {
	if e.Err == nil {
		return "invalid error"
	}
	return e.Err.Error()
}
