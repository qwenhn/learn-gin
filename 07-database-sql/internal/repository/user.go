package repository

import (
	"database/sql"
	"fmt"

	"github.com/qwenhn/gin-restful-api/07-database-sql/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindById(id int, user *models.User) error
}

type SQLUserRepository struct {
	db *sql.DB
}

func (ur *SQLUserRepository) Create(user *models.User) error {
	row := ur.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING user_id", user.Name, user.Email)

	err := row.Scan(&user.Id)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *SQLUserRepository) FindById(id int, user *models.User) error {
	row := ur.db.QueryRow("SELECT * FROM users WHERE user_id = $1", id)
	err := row.Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
		return err
	}

	return nil
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}
