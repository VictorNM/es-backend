package auth

import "log"

func MustHashPassword(password string) string {
	hashed, err := hashPassword(password)
	if err != nil {
		log.Panic(err)
	}

	return hashed
}

func GenerateToken(jwt JWTService, u *User) (string, error) {
	return jwt.generateToken(u)
}
