package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/db"
	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/db/sqlc"
)

type UserRepository interface {
	CountUsers(ctx context.Context, search string, deleted bool) (int64, error)
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error)
	GetAllV2(ctx context.Context, search, orderBy, sort string, limit, offset int32, deleted bool) ([]sqlc.User, error)
	GetByEmail(ctx context.Context, email string) (sqlc.User, error)
	GetByUuid(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	SoftDelete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	UpdatePassword(ctx context.Context, input sqlc.UpdatePasswordParams) (sqlc.User, error)
}

type SQLUserRepository struct {
	db sqlc.Querier
}

func NewSQLUserRepository(db sqlc.Querier) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}

func (ur *SQLUserRepository) CountUsers(ctx context.Context, search string, deleted bool) (int64, error) {
	total, err := ur.db.CountUsers(ctx, sqlc.CountUsersParams{
		Search:  search,
		Deleted: &deleted,
	})

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (ur *SQLUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) Delete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.TrashUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) GetAll(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32) ([]sqlc.User, error) {
	var (
		users []sqlc.User
		err   error
	)

	switch {
	case orderBy == "user_id" && sort == "asc":
		users, err = ur.db.ListUsersUserCreatedAtAsc(ctx, sqlc.ListUsersUserCreatedAtAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_id" && sort == "desc":
		users, err = ur.db.ListUsersUserIdDesc(ctx, sqlc.ListUsersUserIdDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_created_at" && sort == "asc":
		users, err = ur.db.ListUsersUserCreatedAtAsc(ctx, sqlc.ListUsersUserCreatedAtAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_created_at" && sort == "desc":
		users, err = ur.db.ListUsersUserCreatedAtDesc(ctx, sqlc.ListUsersUserCreatedAtDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	}

	if err != nil {
		return []sqlc.User{}, err
	}

	return users, nil
}

func (ur *SQLUserRepository) GetAllV2(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32, deleted bool) ([]sqlc.User, error) {
	query := `SELECT *
		FROM users
		WHERE (
			$1::TEXT IS NULL
			OR $1::TEXT = ''
			OR user_email ILIKE '%' || $1 || '%'
			OR user_fullname ILIKE '%' || $1 || '%'
		)`

	if deleted {
		query += " AND user_deleted_at IS NOT NULL"
	} else {
		query += " AND user_deleted_at IS NULL"
	}

	order := "ASC"
	if sort == "desc" {
		order = "DESC"
	}

	switch orderBy {
	case "user_id", "user_created_at":
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	default:
		query += " ORDER BY user_id ASC"
	}

	query += " LIMIT $2 OFFSET $3 -- name: Get All Version 2"

	rows, err := db.DBPool.Query(ctx, query, search, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []sqlc.User{}
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.UserID,
			&i.UserUuid,
			&i.UserEmail,
			&i.UserPassword,
			&i.UserFullname,
			&i.UserAge,
			&i.UserStatus,
			&i.UserLevel,
			&i.UserDeletedAt,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *SQLUserRepository) GetByEmail(ctx context.Context, email string) (sqlc.User, error) {
	user, err := ur.db.GetUserByEmail(ctx, email)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) GetByUuid(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.GetUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.RestoreUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.SoftDeleteUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.db.UpdateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SQLUserRepository) UpdatePassword(ctx context.Context, input sqlc.UpdatePasswordParams) (sqlc.User, error) {
	user, err := ur.db.UpdatePassword(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}
