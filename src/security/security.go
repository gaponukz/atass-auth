package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
)

func GenerateCode() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

func GetSha256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
