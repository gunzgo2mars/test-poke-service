package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHash(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

func Compare(hashed []byte, p []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, p)
}
