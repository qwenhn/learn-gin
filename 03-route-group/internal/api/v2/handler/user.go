package v2Handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all users (v2)",
	})
}

func (uh *UserHandler) GetUsersById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot parse ID",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get users by ID (v2)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Create user (v2)",
	})
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot parse ID",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update user (v2)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot parse ID",
		})
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "Delete user (v2)",
		"data": gin.H{
			"id": id,
		},
	})
}
