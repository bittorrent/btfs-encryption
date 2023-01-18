package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func EncryptByRSA(content []byte, publicKey *rsa.PublicKey) (out []byte, err error) {
	return rsa.EncryptOAEP(
		sha256.New(), rand.Reader,
		publicKey,
		content,
		nil,
	)
}

func DecryptByRSA(content []byte, privateKey *rsa.PrivateKey) (out []byte, err error) {
	return privateKey.Decrypt(
		nil,
		content,
		&rsa.OAEPOptions{
			Hash: crypto.SHA256,
		},
	)
}

func LoadRSAPrivateKeyFromFile(path string) (key *rsa.PrivateKey, err error) {
	bs, err := ioutil.ReadFile("/Users/steve.zhang/.ssh/id_rsa")
	if err != nil {
		return
	}

	sshKey, err := ssh.ParseRawPrivateKey(bs)
	if err != nil {
		return
	}

	key = sshKey.(*rsa.PrivateKey)

	return
}

func LoadRSAPublicKeyFromFile(path string) (key *rsa.PublicKey, err error) {
	bs, err := ioutil.ReadFile("/Users/steve.zhang/.ssh/id_rsa.pub")
	if err != nil {
		return
	}

	parsed, _, _, _, err := ssh.ParseAuthorizedKey(bs)
	if err != nil {
		return
	}

	pub := parsed.(ssh.CryptoPublicKey).CryptoPublicKey()
	key = pub.(*rsa.PublicKey)

	return
}
