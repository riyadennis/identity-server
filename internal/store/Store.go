package store

// User hold information needed to complete user registration
type User struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Company          string `json:"company"`
	PostCode         string `json:"post_code"`
	Terms            bool   `json:"terms"`
	RegistrationDate string
}

//TODO to consolidate these functions between two database type
type Store interface {
	Insert(u *User) error
	Read(email string) (*User, error)
	Authenticate(email, password string) (bool, error)
	Delete(email string) (bool, error)
}
