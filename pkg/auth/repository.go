package auth

type ReadUserRepository interface {
	FindUserByID(id int) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByUsername(username string) (*User, error)
}

type WriteUserRepository interface {
	CreateUser(u *User) (int, error)
}

type UserRepository interface {
	ReadUserRepository
	WriteUserRepository
}
