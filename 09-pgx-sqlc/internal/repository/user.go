package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/db/sqlc"
)

type UserRepository interface {
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	FindByUuid(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
}

type SQLUserRepository struct {
	db sqlc.Querier
}

func (ur *SQLUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (ur *SQLUserRepository) FindByUuid(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.GetUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func NewSQLUserRepository(db sqlc.Querier) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}
