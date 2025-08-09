package storage

import (
	"encoding/json"
	"fmt"
	"os"
	u "task-tracker/user"
)

const fullFilePath = "./data/data.json"
const directoryPath = "./data"

var GUsers = make([]u.User, 0)

func ReadFromFile() (*[]byte, error) {
	byteValue, err := os.ReadFile(fullFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		fmt.Println("Creating new dir for save file...")
		_, err := os.Lstat(directoryPath)
		if err != nil {
			err = os.Mkdir(directoryPath, 0755)
			if err != nil {
				return nil, err
			}
		}
		fmt.Println("Creating new file for save data...")
		_, err = os.Create(fullFilePath)
		if err != nil {
			return nil, err
		}
		fmt.Println("New data file created successfully.")
	}
	return &byteValue, nil
}

func SaveToFile() error {
	updatedJSON, err := json.MarshalIndent(GUsers, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON or File is empty: %v", err)
	}
	err = os.WriteFile(fullFilePath, updatedJSON, 0666)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
