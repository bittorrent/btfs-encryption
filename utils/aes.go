package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func EncryptByAES(content []byte, key []byte) (out []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	out = aesGCM.Seal(nonce, nonce, content, nil)

	return
}

func DecryptByAES(content []byte, key []byte) (out []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := content[:nonceSize], content[nonceSize:]

	out, err = aesGCM.Open(nil, nonce, ciphertext, nil)

	return
}
