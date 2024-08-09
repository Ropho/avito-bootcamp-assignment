package hash

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

// SaltLen defines the length of "Salt" for password hashing
const SaltLen uint32 = 32

// PassLen defines the length of a random password for collector
const PassLen uint32 = 32

// GenerateRandomString creates random crypto string for encoding
func GenerateRandomString(length uint32) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// ComputeHashWithSalt generates hash for password + salt
func ComputeHashWithSalt(pass, salt string) string {

	h := hmac.New(sha512.New, []byte(salt))
	h.Write([]byte(pass))
	hash := h.Sum(nil)

	return hex.EncodeToString(hash)
}

// ComputeHash generates hash for password w/o salt
func ComputeHash(pass string) string {

	h := sha512.New()
	h.Write([]byte(pass))
	hash := h.Sum(nil)

	return hex.EncodeToString(hash)
}

// IsValidPassWithSalt checks whether given password is equal to hash
func IsValidPassWithSalt(hash, pass, salt string) bool {
	newHash := ComputeHashWithSalt(pass, salt)
	return hmac.Equal([]byte(hash), []byte(newHash))
}

// IsValidPass checks whether given password is equal to hash
func IsValidPass(hash, pass string) bool {
	newHash := ComputeHash(pass)
	return hmac.Equal([]byte(hash), []byte(newHash))
}
