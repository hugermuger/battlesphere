package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	params := argon2id.DefaultParams
	hashedPassword, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
