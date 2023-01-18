package enc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

type Header struct {
	Version                   string                 `json:"version"`
	SecretEncryptionAlgorithm string                 `json:"secret_encrypt_algorithm"`
	Secret                    string                 `json:"secret"`
	EncryptionAlgorithm       string                 `json:"encryption_algorithm"`
	Encoding                  string                 `json:"encoding"`
	PublicInfo                map[string]interface{} `json:"public_info"`
}

func pack(header *Header, payload []byte) (packed []byte, err error) {
	headerRaw, err := json.Marshal(header)
	if err != nil {
		return
	}

	headerLen := make([]byte, 8)
	binary.BigEndian.PutUint64(headerLen, uint64(len(headerRaw)))

	buffer := &bytes.Buffer{}
	_, err = buffer.Write(headerLen)
	if err != nil {
		return
	}

	_, err = buffer.Write(headerRaw)
	if err != nil {
		return
	}

	_, err = buffer.Write(payload)
	if err != nil {
		return
	}

	packed = buffer.Bytes()
	return
}

func unpack(content []byte) (header *Header, payload []byte, err error) {
	if len(content) < 8 {
		err = errors.New("invalid content length")
		return
	}
	headerLen := content[:8]
	headerLenInt64 := binary.BigEndian.Uint64(headerLen)

	if len(content) < int(headerLenInt64) {
		err = errors.New("invalid content length")
		return
	}
	headerRaw := content[8 : 8+headerLenInt64]

	header = &Header{}
	err = json.Unmarshal(headerRaw, header)
	if err != nil {
		return
	}

	payload = content[8+headerLenInt64:]
	return
}
