package security

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, password string) bool {
	byteHash := []byte(hashedPwd)
	passwordHash := []byte(password)
	if err := bcrypt.CompareHashAndPassword(byteHash, passwordHash); err != nil {
		return false
	}

	return true
}
