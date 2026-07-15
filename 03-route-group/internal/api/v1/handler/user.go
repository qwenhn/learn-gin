package v1Handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/03-route-group/utils"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all users (v1)",
	})
}

func (uh *UserHandler) GetUsersById(ctx *gin.Context) {
	id, err := utils.ValidationPositiveInt("ID", ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get users by ID (v1)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Create user (v1)",
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
		"message": "Update user (v1)",
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
		"message": "Delete user (v1)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (uh *UserHandler) GetByUuid(ctx *gin.Context) {
	uid, err := utils.ValidationUuid("UUID", ctx.Param("uid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get users by UUID (v1)",
		"data": gin.H{
			"uuid": uid,
		},
	})
}
