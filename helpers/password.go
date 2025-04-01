package helpers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the plain password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // bcrypt.DefaultCost is 10
	if err != nil {
		fmt.Println("❌ Error hashing password:", err)
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares plain password with the hashed one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("❌ Password mismatch!", err)
		return false
	}
	fmt.Println("✅ Password matched!")
	return true
}
