package services

import (
	"GoSQL/internal/config"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var key, _ = hex.DecodeString("04cc5f1d20f5bd74bfb034ebceac2094134ecfdc7922012a3c99189fb7c00417")

func GetProfiles() ([]config.DatabaseConnectionInput, error) {
	profiles, err := readEncryptedConfigs("profiles.json")
	if err != nil {
		return nil, err
	}
	fmt.Println(profiles)
	return profiles, nil
}

func CreateProfile(dto config.DatabaseConnectionInput, ctx context.Context) error {
	err := config.ConnectToDb(dto, ctx)
	if err != nil {
		return fmt.Errorf("Error Connecting to database: %v", err)

	}

	err = writeEncryptedConfig("profiles.json", dto)
	if err != nil {
		return fmt.Errorf("Error Saving Profile: %v", err)

	}
	return nil

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

func writeEncryptedConfig(filename string, config config.DatabaseConnectionInput) error {
	encPass, err := encryptAES(config.Password)
	if err != nil {
		return err
	}
	config.Password = encPass

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	existingConfigs, err := GetProfiles()
	if err != nil {
		return err
	}
	existingConfigs = append(existingConfigs, config)

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(existingConfigs)
}
func readEncryptedConfigs(filename string) ([]config.DatabaseConnectionInput, error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var profiles []config.DatabaseConnectionInput
	defer jsonFile.Close()
	json.Unmarshal(byteValue, &profiles)
	fmt.Println(profiles)
	return profiles, nil
}
