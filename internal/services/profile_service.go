package services

import (
	"GoSQL/helpers"
	"GoSQL/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetProfiles() ([]config.DatabaseConnectionInput, error) {
	profiles, err := readEncryptedConfigs("profiles.json")
	if err != nil {
		return nil, err
	}
	fmt.Println(profiles)
	return profiles, nil
}

func CreateProfile(dto config.DatabaseConnectionInput, ctx context.Context) error {
	err := config.ConnectToDb(dto, ctx, false)
	if err != nil {
		return fmt.Errorf("Error Connecting to database: %v", err)

	}

	err = writeEncryptedConfig("profiles.json", dto)
	if err != nil {
		return fmt.Errorf("Error Saving Profile: %v", err)

	}
	return nil

}

func writeEncryptedConfig(filename string, config config.DatabaseConnectionInput) error {
	encPass, err := helpers.EncryptAES(config.Password)
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
