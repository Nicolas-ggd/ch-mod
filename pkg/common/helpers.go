package common

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
)

var (
	secret = os.Getenv("PRIVATE_KEY")
)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", nil
	}

	hashedPassword := string(hash)

	return hashedPassword, nil
}

func ComparePasswords(password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	return nil
}

func CompareHashAndPasswordBcrypt(hash, rawPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(rawPassword))
	if err != nil {
		return fmt.Errorf("incorrect password: %w", err)
	}

	return nil
}

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
