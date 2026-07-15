package v1Handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/03-route-group/utils"
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

func (ch *CategoryHandler) GetByCategory(ctx *gin.Context) {
	category := ctx.Param("category")

	if err := utils.ValidationInList("Category", category, Categories); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Get by category (v1)",
		"category": category,
	})
}
