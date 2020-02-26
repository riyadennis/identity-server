package handlers

const (
	FirstNameMissing = "first-name-missing"
	LastNameMissing  = "last-name-missing"
	EmailMissing     = "email-missing"
	TermsMissing     = "terms-missing"
	EmailInvalid     = "email-invalid"
	EmailExists      = "email-exists"
	DatabaseError    = "database-error"
	PassWordError    = "password-error"
	InvalidRequest   = "invalid-request"
	TokenError = "token-error"
	UnAuthorised = "unauthorised"
	UserDoNotExist = "user-do-not-exist"
)

type CustomError struct {
	Code string
	Err  error
}

func (e *CustomError) Error() string {
	return e.Err.Error()
}
