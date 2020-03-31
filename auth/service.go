package auth

type LoginService struct {
}

type UserDTO struct {
	Email          string
	HashedPassword string
}

type GetUserByEmail interface {
	GetUserByEmail(email string) (UserDTO, error)
}

func CreateBasicLogin(getter GetUserByEmail) func(email, password string) error {
	return func(email, password string) error {
		return nil
	}
}
