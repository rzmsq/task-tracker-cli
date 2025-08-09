package main

import (
	"encoding/json"
	"fmt"
	"task-tracker/cli"
	"task-tracker/consts"
	s "task-tracker/storage"
	u "task-tracker/user"
	"time"
)

func main() {

	var currentUser *u.User // Current authenticated user
	var isLoggedIn bool

	// Read existing users from file or create a new file if it doesn't exist
	data, err := s.ReadFromFile()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Unmarshal the JSON data into the global users slice
	err = json.Unmarshal(*data, &s.GUsers)
	if err != nil {
		fmt.Println("Error unmarshalling JSON or File is empty")
		fmt.Println("Continue...")
	}

	for {
		if !isLoggedIn {
			currentUser, isLoggedIn = WelcomeMenuLoop()
			if !isLoggedIn {
				break
			}
		}
		isLoggedIn = MainMenuLoop(currentUser)
	}

}

func WelcomeMenuLoop() (currentUser *u.User, isLoggedIn bool) {
	isLoggedIn = false

	cli.WelcomePrint()

	for !isLoggedIn {
		// Display the welcome menu and get user choice
		input, err := cli.UChoice()
		if err != nil {
			fmt.Println("Error:", err)
			return nil, false
		}

		// Process user choice
		switch input {
		case consts.CreateUser:
			currentUser, err = cli.CreateUser()
			if err != nil {
				fmt.Println("Error creating user:", err)
				break
			}
			fmt.Printf("User created successfully! ID: %d, Name: %s\n", currentUser.ID, currentUser.Name)
			isLoggedIn = true
		case consts.LogIn:
			currentUser, err = cli.LogInUser()
			if err != nil {
				fmt.Println("Error creating user:", err)
				break
			}
			fmt.Printf("Log in successfully! ID: %d, Name: %s\n", currentUser.ID, currentUser.Name)
			isLoggedIn = true
		case consts.DeleteUser:
			currentUser, err = cli.DeleteUser()
			if err != nil {
				fmt.Println("Error deleting user:", err)
				break
			}
			fmt.Printf("User deleted successfully! ID: %d, Name: %s\n", currentUser.ID, currentUser.Name)
		case consts.Exit:
			err = s.SaveToFile()
			if err != nil {
				fmt.Println("Error saving data:", err)
			}
			fmt.Println("Exiting the program. Goodbye!")
			return nil, false
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-4).")
		}
	}
	return currentUser, isLoggedIn
}

func MainMenuLoop(currentUser *u.User) bool {
	for {
		cli.ProgOptionPrint(currentUser)

		input, err := cli.UChoice()
		if err != nil {
			fmt.Println("Error:", err)
			return true
		}

		switch input {
		case consts.CreateTask:
			title, description, err := cli.GetTaskNameDesc()
			if err != nil {
				fmt.Println("Error getting task details:", err)
				continue
			}
			// Check if the task title already exists for the current user
			for _, task := range currentUser.Tasks {
				if task.Title == title {
					fmt.Println("Task with this title already exists. Please choose a different title.")
					continue
				}
			}
			task := u.Task{
				Title:       title,
				Description: description,
				Date:        time.Now().Format("2006-01-02 15:04:05"),
			}
			currentUser.Tasks = append(currentUser.Tasks, task)
			fmt.Printf("Task created successfully! Title: %s, Description: %s\n", title, description)
		case consts.MarkTaskDone:
			if len(currentUser.Tasks) == 0 {
				fmt.Println("No tasks available to mark as done.")
				continue
			}
			cli.PrintAllTasks(*currentUser)
			fmt.Println("Enter the task number to mark as done:")
			var taskNumber int
			_, err := fmt.Scanln(&taskNumber)
			if err != nil || taskNumber < 1 || taskNumber > len(currentUser.Tasks) {
				fmt.Println("Invalid task number. Please try again.")
				continue
			}
			currentUser.Tasks[taskNumber-1].IsDone = true
			fmt.Printf("Task '%s' marked as done.\n", currentUser.Tasks[taskNumber-1].Title)
		case consts.ViewTasks:
			if len(currentUser.Tasks) == 0 {
				fmt.Println("No tasks available.")
				continue
			}
			cli.PrintAllTasks(*currentUser)
		case consts.LogOut:
			err = s.SaveToFile()
			if err != nil {
				fmt.Println("Error saving data:", err)
				continue
			}
			fmt.Println("You have been logged out.")
			return false
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-4).")
		}
	}
}
