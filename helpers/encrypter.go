package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

var key, _ = hex.DecodeString("04cc5f1d20f5bd74bfb034ebceac2094134ecfdc7922012a3c99189fb7c00417")

func EncryptAES(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainData := []byte(plainText)
	ciphertext := make([]byte, aes.BlockSize+len(plainData))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainData)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return string(data), nil
}
