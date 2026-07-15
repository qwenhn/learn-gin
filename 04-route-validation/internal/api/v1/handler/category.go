package v1Handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/04-route-validation/utils"
)

type CategoryHandler struct {
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

var Categories = map[string]bool{
	"php": true,
	"go":  true,
	"js":  true,
}

type GetByCategoryParam struct {
	Category string `uri:"category" binding:"oneof=php go js"`
}

func (ch *CategoryHandler) GetByCategory(ctx *gin.Context) {
	var params GetByCategoryParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Get by category (v1)",
		"category": params.Category,
	})
}
