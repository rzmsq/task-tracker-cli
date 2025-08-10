package cli

import (
	"bufio"
	"fmt"
	"os"
	s "task-tracker/storage"
	u "task-tracker/user"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func WelcomePrint() {
	fmt.Println("---Welcome to the Task Management System ---")
	fmt.Println("1. Create User")
	fmt.Println("2. Log in")
	fmt.Println("3. Delete User")
	fmt.Println("4. Exit")
}

func ProgOptionPrint(currentUser *u.User) {
	_ = fmt.Sprintf("Hello, %s!", currentUser.Name)
	fmt.Println("Choose an option:")
	fmt.Println("1. Create Task")
	fmt.Println("2. Mark Task as Done")
	fmt.Println("3. View Tasks")
	fmt.Println("4. Log out")
	fmt.Print("> ")
}

func UChoice() (int, error) {
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return 0, err
	}

	return choice, nil
}

func GetNamePass() (string, []byte, error) {
	var name string

	fmt.Println("Enter your name:")
	_, err := fmt.Scanln(&name)
	if err != nil {
		return "", nil, fmt.Errorf("error reading name: %v", err)
	}

	fmt.Println("Enter your password:")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", nil, fmt.Errorf("error reading password: %v", err)
	}

	if len(bytePassword) < 3 {
		fmt.Println("Password must be at least 3 characters long.")
		return "", nil, fmt.Errorf("password too short")
	}

	return name, bytePassword, nil
}

func GetTaskNameDesc() (string, string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var title, description string

	fmt.Println("Enter task title:")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Enter task description:")
	scanner.Scan()
	description = scanner.Text()
	return title, description, nil
}

func CreateUser() (*u.User, error) {

	name, password, err := GetNamePass()
	if err != nil {
		return nil, fmt.Errorf("error getting name and password: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	var newUser u.User
	newUser = newUser.CreateUser(name, hashedPassword)

	// Check if the user already exists
	for _, user := range s.GUsers {
		if user.Name == newUser.Name {
			return nil, fmt.Errorf("user with name %s already exists", newUser.Name)
		}
	}

	s.GUsers = append(s.GUsers, newUser)

	return &s.GUsers[len(s.GUsers)-1], nil
}

func LogInUser() (*u.User, error) {
	name, password, err := GetNamePass()
	if err != nil {
		return nil, fmt.Errorf("error getting name and password: %v", err)
	}

	for _, user := range s.GUsers {
		err = bcrypt.CompareHashAndPassword(user.Password, password)
		if user.Name == name && err == nil {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found or incorrect password")
}

func DeleteUser() (*u.User, error) {
	name, password, err := GetNamePass()
	if err != nil {
		return nil, fmt.Errorf("error getting name and password: %v", err)
	}

	for i, user := range s.GUsers {
		err = bcrypt.CompareHashAndPassword(user.Password, password)
		if user.Name == name && err == nil {
			s.GUsers = append(s.GUsers[:i], s.GUsers[i+1:]...)
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found or incorrect password")
}

func PrintAllTasks(currentUser u.User) {
	fmt.Println("All tasks:")
	for i, task := range currentUser.Tasks {
		var status string
		if task.IsDone {
			status = "Done"
		} else {
			status = "Not Done"
		}
		fmt.Printf("%d. %s Status: %s \n\tDescription: %s\n\tDate: %s\n", i+1, task.Title, status,
			task.Description, task.Date)
	}
}
