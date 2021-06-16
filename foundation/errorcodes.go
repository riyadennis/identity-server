package foundation

const (
	// ValidationFailed is the error code for
	// validation failures in a POST request
	ValidationFailed = "validation-failure"

	// EmailAlreadyExists is error code for duplicate registration
	EmailAlreadyExists = "email-already-exists"

	// DatabaseError is when there is issues with db connection
	DatabaseError = "database-error"

	// PassWordError is to tell the user that password generation failed
	PassWordError = "password-error"

	// InvalidRequest is returned if request is not a valid one
	InvalidRequest = "invalid-request"

	// TokenError is returned if we are not able to generate a token
	TokenError = "Token-error"

	// KeyNotFound is returned if we are not able to find key
	// that we need to encrypt and decrypt tokens
	KeyNotFound = "Key-not-found"

	// UnAuthorised is when a user have invalid or expired token
	UnAuthorised = "unauthorised"

	// UserDoNotExist is returned when search for an email in db fails
	UserDoNotExist = "user-do-not-exist"
)

// CustomError holds error code and details about the error
type CustomError struct {
	Code string
	Err  error
}

// Error returns just the error message for a custom error
func (e *CustomError) Error() string {
	return e.Err.Error()
}
