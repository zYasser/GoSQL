package services

import (
	"GoSQL/helpers"
	"GoSQL/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

func GetProfiles() (map[string]config.DatabaseConnectionInput, error) {
	profiles, err := readEncryptedConfigs("profiles.json")
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func DeleteProfile(id string) error {
	profiles, err := GetProfiles()
	if err != nil {
		return err
	}
	delete(profiles, id)

	file, err := os.OpenFile("profiles.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(profiles)

}

func GetProfile(id string) (config.DatabaseConnectionInput, error) {
	profiles, err := GetProfiles()
	if err != nil {
		return config.DatabaseConnectionInput{}, err
	}
	return profiles[id], nil
}

func CreateProfile(dto config.DatabaseConnectionInput, ctx context.Context) error {
	err := config.ConnectToDb(dto, ctx, false)
	if err != nil {
		return fmt.Errorf("Error Connecting to database: %v", err)

	}

	err = writeEncryptedConfig("profiles.json", dto, false, "")
	if err != nil {
		return fmt.Errorf("Error Saving Profile: %v", err)

	}
	return nil

}

func UpdateProfile(dto config.DatabaseConnectionInput, id string, ctx context.Context) error {
	err := config.ConnectToDb(dto, ctx, false)
	if err != nil {
		return fmt.Errorf("Error Connecting to database: %v", err)

	}
	err = writeEncryptedConfig("profiles.json", dto, true, id)
	if err != nil {
		return fmt.Errorf("Error Saving Profile: %v", err)

	}

	return nil

}

func writeEncryptedConfig(filename string, profile config.DatabaseConnectionInput, update bool, id string) error {
	encPass, err := helpers.EncryptAES(profile.Password)
	if err != nil {
		return err
	}
	profile.Password = encPass
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	existingConfigs, err := GetProfiles()
	if err != nil {
		return err
	}
	if existingConfigs == nil {
		existingConfigs = make(map[string]config.DatabaseConnectionInput)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	if update {
		existingConfigs[id] = profile
	} else {
		profile.ID = uuid.New().String()

		existingConfigs[profile.ID] = profile
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(existingConfigs)
}
func readEncryptedConfigs(filename string) (map[string]config.DatabaseConnectionInput, error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var profiles map[string]config.DatabaseConnectionInput
	defer jsonFile.Close()
	json.Unmarshal(byteValue, &profiles)
	return profiles, nil
}
