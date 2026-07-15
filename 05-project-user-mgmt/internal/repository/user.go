package repository

import (
	"fmt"
	"slices"

	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/models"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	Create(user models.User) error
	FindByUUID(uuid string) (models.User, bool)
	Update(uuid string, user models.User) error
	Delete(uuid string) error
	FindByEmail(email string) (models.User, bool)
}

type inMemoryUserRepository struct {
	users []models.User
}

func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make([]models.User, 0),
	}
}

func (ur *inMemoryUserRepository) Create(user models.User) error {
	ur.users = append(ur.users, user)
	return nil
}

func (ur *inMemoryUserRepository) Delete(uuid string) error {
	for i, user := range ur.users {
		if user.UUID == uuid {
			ur.users = slices.Delete(ur.users, i, i+1)
			return nil
		}
	}

	return fmt.Errorf("User not found")
}

func (ur *inMemoryUserRepository) FindAll() ([]models.User, error) {
	return ur.users, nil
}

func (ur *inMemoryUserRepository) FindByEmail(email string) (models.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}

	return models.User{}, false
}

func (ur *inMemoryUserRepository) FindByUUID(uuid string) (models.User, bool) {
	for _, user := range ur.users {
		if user.UUID == uuid {
			return user, true
		}
	}

	return models.User{}, false
}

func (ur *inMemoryUserRepository) Update(uuid string, user models.User) error {
	for i, u := range ur.users {
		if u.UUID == uuid {
			ur.users[i] = user
			return nil
		}
	}

	return fmt.Errorf("User not found")
}
