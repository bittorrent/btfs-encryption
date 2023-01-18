package utils

import (
	"crypto/rand"
	"io"
)

func GenerateSecret() (secret []byte, err error) {
	secret = make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, secret)
	return
}
