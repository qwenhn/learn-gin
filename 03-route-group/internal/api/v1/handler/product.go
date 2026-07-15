package v1Handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/03-route-group/utils"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

var (
	searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	slugRegex   = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
)

func (ph *ProductHandler) GetAllProducts(ctx *gin.Context) {
	search := ctx.Query("search")

	if err := utils.ValidationRequired("Search", search); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationLength("Search", search, 3, 50); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidationRegex("Search", search, searchRegex, "must contain only letters, numbers and spaces"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be a positive number"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all products (v1)",
		"search":  search,
		"limit":   limit,
	})
}

func (ph *ProductHandler) GetProductsById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot parse ID",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get products by ID (v1)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Create product (v1)",
	})
}

func (ph *ProductHandler) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot parse ID",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update product (v1)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (ph *ProductHandler) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot parse ID",
		})
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "Delete product (v1)",
		"data": gin.H{
			"id": id,
		},
	})
}

func (ph *ProductHandler) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if err := utils.ValidationRegex("Slug", slug, slugRegex, "must contain only lowercase letter, numbers, hyphens and dots"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get by slug (v1)",
		"slug":    slug,
	})
}
