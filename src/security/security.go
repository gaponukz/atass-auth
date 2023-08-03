package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func GenerateCode() string {
	var n uint32
	_ = binary.Read(rand.Reader, binary.BigEndian, &n)
	code := n%900000 + 100000
	return fmt.Sprintf("%d", code)
}

func Sha256WithSecretFactory(secret string) func(string) string {
	return func(input string) string {
		hash := sha256.Sum256([]byte(input + secret))
		return hex.EncodeToString(hash[:])
	}
}
