package enc

import (
	"context"
	"encoding/hex"
	"github.com/bittorrent/btfs-encryption/btfs"
	"github.com/bittorrent/btfs-encryption/utils"
	"io/ioutil"
	"path/filepath"
)

func Encrypt(source, publicKeyPath string) (content []byte, err error) {
	publicKey, err := utils.LoadRSAPublicKeyFromFile(publicKeyPath)
	if err != nil {
		return
	}

	payload, err := utils.TarGiz(source)
	if err != nil {
		return
	}

	key, err := utils.GenerateSecret()
	if err != nil {
		return
	}

	encryptedPayload, err := utils.EncryptByAES(payload, key)
	if err != nil {
		return
	}

	encryptKey, err := utils.EncryptByRSA(key, publicKey)
	if err != nil {
		return
	}

	header := &Header{
		Version:                   "1.0.0",
		SecretEncryptionAlgorithm: "RSA",
		Secret:                    hex.EncodeToString(encryptKey),
		EncryptionAlgorithm:       "AES",
		Encoding:                  "tar/giz",
		PublicInfo: map[string]interface{}{
			"author": "demo",
		},
	}

	content, err = pack(header, encryptedPayload)
	return
}

func Decrypt(content []byte, dest, privateKeyPath string) (err error) {
	header, encryptPayload, err := unpack(content)
	if err != nil {
		return
	}

	encryptKey, err := hex.DecodeString(header.Secret)
	if err != nil {
		return
	}

	privateKey, err := utils.LoadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return
	}

	key, err := utils.DecryptByRSA(encryptKey, privateKey)
	if err != nil {
		return
	}

	payload, err := utils.DecryptByAES(encryptPayload, key)
	if err != nil {
		return
	}

	err = utils.UnTarGiz(payload, dest)
	return
}

func EncryptToBTFS(source, publicKeyPath string) (rst *btfs.AddResult, err error) {
	content, err := Encrypt(source, publicKeyPath)
	if err != nil {
		return
	}

	rst, err = btfs.Add(
		context.Background(),
		content,
		filepath.Base(source)+".bte",
		&btfs.AddOptions{
			Pin: false,
		},
	)
	return
}

func DecryptFromBTFS(cid, dest, privateKeyPath string) (err error) {
	content, err := btfs.Cat(context.Background(), cid)
	if err != nil {
		return
	}

	err = Decrypt(content, dest, privateKeyPath)
	return
}

func EncryptToLocal(source, dest, publicKeyPath string) (err error) {
	content, err := Encrypt(source, publicKeyPath)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(
		filepath.Join(dest, filepath.Base(source)+".bte"),
		content,
		0755,
	)
	return
}

func DecryptFromLocal(source, dest, privateKeyPath string) (err error) {
	content, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	err = Decrypt(content, dest, privateKeyPath)
	return
}
