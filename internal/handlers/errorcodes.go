package handlers

const (
	InvalidUserData    = "invalid-user-data"
	EmailMissing       = "email-missing"
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
	return e.Err.Error()
}
