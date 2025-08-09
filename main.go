package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	CreateUser = iota + 1
	LogIn
	Exit
)

const fullFilePath = "./data/data.json"
const directoryPath = "./data"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks,omitempty"`
}

type Task struct {
	ID          int
	UserID      int
	Title       string
	Description string
	Date        string
}

var gUsers []User = make([]User, 0)

func (u *User) CreateUser(name, password string) User {
	return User{
		ID:       time.Now().Nanosecond(),
		Name:     name,
		Password: password,
	}
}

func welcomePrint() {
	fmt.Println("---Welcome to the Task Management System ---")
	fmt.Println("1. Create User")
	fmt.Println("2. Log in")
	fmt.Println("3. Exit")
}

func uChoice() (int, error) {
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return 0, err
	}

	if choice < 1 || choice > 3 {
		return 0, fmt.Errorf("invalid choice. Please select a valid option (1-3)")
	}

	return choice, nil
}

func getNamePass() (string, string, error) {
	var name, password string

	fmt.Println("Enter your name:")
	_, err := fmt.Scanln(&name)
	if err != nil {
		return "", "", fmt.Errorf("error reading name: %v", err)
	}

	fmt.Println("Enter your password:")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", fmt.Errorf("error reading password: %v", err)
	}

	password = string(bytePassword)
	if len(password) < 3 {
		fmt.Println("Password must be at least 3 characters long.")
		return "", "", fmt.Errorf("password too short")
	}

	return name, password, nil
}

func createUser() (*User, error) {

	name, password, err := getNamePass()
	if err != nil {
		return nil, fmt.Errorf("error getting name and password: %v", err)
	}

	var newUser User
	newUser = newUser.CreateUser(name, password)
	gUsers = append(gUsers, newUser)

	return &gUsers[len(gUsers)-1], nil
}

func logInUser() (*User, error) {
	name, password, err := getNamePass()
	if err != nil {
		return nil, fmt.Errorf("error getting name and password: %v", err)
	}

	for _, user := range gUsers {
		if user.Name == name && user.Password == password {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found or incorrect password")
}

func readFromFile() (*[]byte, error) {
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

func saveToFile() error {
	updatedJSON, err := json.MarshalIndent(gUsers, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON or File is empty: %v", err)
	}
	err = os.WriteFile(fullFilePath, updatedJSON, 0666)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func main() {

	var currentUser *User

	welcomePrint()

	// Read existing users from file or create a new file if it doesn't exist
	data, err := readFromFile()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Unmarshal the JSON data into the global users slice
	err = json.Unmarshal(*data, &gUsers)
	if err != nil {
		fmt.Println("Error unmarshalling JSON or File is empty")
		fmt.Println("Continue...")
	}

	for {

		// Display the welcome menu and get user choice
		input, err := uChoice()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch input {
		case CreateUser:
			currentUser, err = createUser()
			if err != nil {
				fmt.Println("Error creating user:", err)
				break
			}
			fmt.Printf("User created successfully! ID: %d, Name: %s\n", currentUser.ID, currentUser.Name)
		case LogIn:
			currentUser, err = logInUser()
			if err != nil {
				fmt.Println("Error creating user:", err)
				break
			}
			fmt.Printf("Log in successfully! ID: %d, Name: %s\n", currentUser.ID, currentUser.Name)
		case Exit:
			err = saveToFile()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Exiting the program. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-3).")
		}
	}
}
