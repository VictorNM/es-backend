package auth

type UserRepository interface {
	FindUserByID(id int) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByUsername(username string) (*User, error)

	CreateUser(u *User) (int, error)
}
