package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

const (
	prefix = "Too Salty"
	surfix = "^[a-zA-Z0-9_]+$"
)

func Encrypt(s string) string {
	cipher := sha1.New()
	io.WriteString(cipher, prefix)
	io.WriteString(cipher, s)
	io.WriteString(cipher, surfix)

	return fmt.Sprintf("%x", cipher.Sum(nil))
}
