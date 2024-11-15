package database

import (
	"errors"
	"fmt"

	"github.com/Garvit-Jethwani/user-management-service/models"
)

var users = []models.User{
	{ID: "1", Name: "John Doe", Email: "john@example.com", Password: "12345"},
}

func CreateUser(user *models.User) error {
	for _, u := range users {
		if u.Email == user.Email {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
	}
	user.ID = fmt.Sprintf("%d", len(users)+1)
	users = append(users, *user)
	return nil
}

func AuthenticateUser(email, password string) (string, error) {
	for _, user := range users {
		if user.Email == email && user.Password == password {
			// Generate a token (dummy token for now)
			return "dummytoken", nil
		}
	}
	return "", errors.New("invalid email or password")
}

func GetUserByID(userID string) (*models.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %s not found", userID)
}
