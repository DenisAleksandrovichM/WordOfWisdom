package pow

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	mrand "math/rand"
	"strings"
)

const HashCashHeader = "X-Hashcash:"

func IsCorrect(hashString string, zerosCount int) bool {
	target := strings.Repeat("0", zerosCount)
	return strings.HasPrefix(hashString, target)
}

func SolvePoW(challenge string, zerosCount int) (string, error) {
	stringLen := zerosCount + mrand.Intn(10) + 5
	for {
		response, err := generateRandomString(stringLen)
		if err != nil {
			return "", fmt.Errorf("generate random string: %w", err)
		}
		hashString := CalculateHash(challenge, response)
		if IsCorrect(hashString, zerosCount) {
			return response, nil
		}
	}
}

func CalculateHash(challenge, response string) string {
	data := []byte(challenge + response)
	hash := sha1.Sum(data)
	return fmt.Sprintf("%x", hash)
}

func generateRandomString(length int) (string, error) {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("rand read: %w", err)
	}

	for i := 0; i < length; i++ {
		b[i] = alphanumeric[b[i]%byte(len(alphanumeric))]
	}

	return string(b), nil
}
