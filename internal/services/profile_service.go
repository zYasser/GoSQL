package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type DatabaseConnectionInput struct {
	ProfileName  string
	DatabaseType string
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

var key = "04cc5f1d20f5bd74bfb034ebceac2094134ecfdc7922012a3c99189fb7c00417"

func CreateProfile(dto DatabaseConnectionInput, ctx context.Context) {
	err := writeEncryptedConfig("profiles.csv", dto)
	if err != nil {
		fmt.Printf("Error Saving Profile: %v", err)
		return
	}
	
}

func encryptAES(plainText string) (string, error) {
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

func decryptAES(cipherText string) (string, error) {
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

func writeEncryptedConfig(filename string, config DatabaseConnectionInput) error {
	encPass, err := encryptAES(config.Password)
	if err != nil {
		return err
	}
	config.Password = encPass

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func readEncryptedConfigs(filename string) ([]DatabaseConnectionInput, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var configs []DatabaseConnectionInput
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configs)
	if err != nil {
		return nil, err
	}

	for i := range configs {
		decPass, err := decryptAES(configs[i].Password)
		if err != nil {
			return nil, err
		}
		configs[i].Password = decPass
	}

	return configs, nil
}
