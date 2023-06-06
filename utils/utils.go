package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"regexp"
)

const (
	prefix = "Too Salty"
	surfix = "^[a-zA-Z0-9_]+$"
)

var zero = regexp.MustCompile(`^[\+\-]?0\.?0*$`)

func Encrypt(s string) string {
	cipher := sha1.New()
	io.WriteString(cipher, prefix)
	io.WriteString(cipher, s)
	io.WriteString(cipher, surfix)

	return fmt.Sprintf("%x", cipher.Sum(nil))
}

func Abs(s string) string {
	if s[0:1] == "-" {
		return s[1:]
	}
	return s
}

func IsZero(s string) bool {
	return zero.MatchString(s)
}
