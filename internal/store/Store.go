package store

// User hold information needed to complete user registration
type User struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Company          string `json:"company"`
	PostCode         string `json:"post_code"`
	Terms            bool   `json:"terms"`
	RegistrationDate string
}

type Store interface {
	Insert(u *User) error
	Read(email string, password string) (*User, error)
}
