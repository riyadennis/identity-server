package handlers

const (
	// RegisterEndpoint is to create a new user
	RegisterEndpoint = "/register"

	// DeleteEndpoint is to delete a user
	DeleteEndpoint = "/delete/:id"

	// LoginEndPoint creates a token for the  user of credentials are valid
	LoginEndPoint = "/login"

	// HomeEndPoint is the details end point that a
	// logged in user with valid token can access
	HomeEndPoint = "/home"
)
