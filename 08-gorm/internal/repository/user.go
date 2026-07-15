package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/qwenhn/gin-restful-api/08-gorm/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindById(id int, user *models.User) error
}

type SQLUserRepository struct {
	db *gorm.DB
}

func (ur *SQLUserRepository) Create(user *models.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *SQLUserRepository) FindById(id int, user *models.User) error {
	if err := ur.db.First(user, id).Error; err != nil {
		return err
	}

	return nil
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}
