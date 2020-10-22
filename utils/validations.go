package utils

import "golang.org/x/crypto/bcrypt"

// ValidatePassword -
func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
