package v1Handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/04-route-validation/utils"
)

type UserHandler struct {
}

type GetUsersByIdParam struct {
	ID int `uri:"id" binding:"gt=0"`
}

type GetByUuidParam struct {
	Uuid string `uri:"uid" binding:"uuid"` // uri ONLY map with params in route
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
	var params GetUsersByIdParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get users by ID (v1)",
		"data": gin.H{
			"id": params.ID,
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
	var params GetByUuidParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get users by UUID (v1)",
		"data": gin.H{
			"uuid": params.Uuid,
		},
	})
}
