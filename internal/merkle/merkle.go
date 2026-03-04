package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func Combine(a, b string) string {
	data := []byte(a + b)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}