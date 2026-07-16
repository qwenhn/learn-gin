package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/db/sqlc"
	"github.com/qwenhn/gin-restful-api/09-pgx-sqlc/internal/repository"
)

type UserHandler struct {
	repo repository.UserRepository
}

type UserResponse struct {
	UserID    int32     `json:"user_id"`
	Uuid      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (uh *UserHandler) GetUserByUuid(ctx *gin.Context) {
	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uh.repo.FindByUuid(ctx, uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		UserID:    user.UserID,
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input sqlc.CreateUserParams
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.repo.Create(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		UserID:    user.UserID,
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}
