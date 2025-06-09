package pkg

import "golang.org/x/crypto/bcrypt"

func NewHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func CompareHash(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
