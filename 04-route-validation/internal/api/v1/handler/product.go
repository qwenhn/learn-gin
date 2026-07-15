package v1Handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qwenhn/gin-restful-api/04-route-validation/utils"
)

type ProductHandler struct {
}

type GetAllProductsParam struct {
	Search string `form:"search" binding:"required,min=3,max=50,search"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Email  string `form:"email" binding:"omitempty,email"`
	Date   string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

type GetBySlugParam struct {
	Slug string `uri:"slug" binding:"min=3,max=5,slug"`
}

type ProductImage struct {
	ImageName string `json:"image_name" binding:"required"`
	ImageLink string `json:"image_link" binding:"required,file_ext=jpg png gif"`
}

type ProductAttribute struct {
	AttributeName  string `json:"attribute_name" binding:"required"`
	AttributeValue string `json:"attribute_value" binding:"required"`
}

type ProductInfo struct {
	InfoKey   string `json:"info_key" binding:"required"`
	InfoValue string `json:"info_value" binding:"required"`
}

type CreateProductParam struct {
	Name              string                 `json:"name" binding:"required,min=3,max=100"`
	Price             int                    `json:"price" binding:"required,gte=10000"`
	Display           *bool                  `json:"display" binding:"omitempty"`
	ProductImage      ProductImage           `json:"product_image" binding:"required"`
	Tags              []string               `json:"tags" binding:"required,gt=3,lt=5"`
	ProductAttributes []ProductAttribute     `json:"product_attributes" binding:"required,gt=0,dive"`
	ProductInfo       map[string]ProductInfo `json:"product_info" binding:"required,gt=0,dive"`
	ProductMetadata   map[string]any         `json:"product_metadata" binding:"omitempty"`
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (ph *ProductHandler) GetAllProducts(ctx *gin.Context) {
	var params GetAllProductsParam

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if params.Email == "" {
		params.Email = "No Email"
	}

	if params.Date == "" {
		params.Date = time.Now().Format("2006-01-02")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all products (v1)",
		"search":  params.Search,
		"limit":   params.Limit,
		"email":   params.Email,
		"date":    params.Date,
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
	var params CreateProductParam

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	if params.Display == nil {
		defaultDisplay := true
		params.Display = &defaultDisplay
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":            "Create product (v1)",
		"name":               params.Name,
		"price":              params.Price,
		"display":            params.Display,
		"product_image":      params.ProductImage,
		"tags":               params.Tags,
		"product_attributes": params.ProductAttributes,
		"product_info":       params.ProductInfo,
		"product_meta":       params.ProductMetadata,
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
	var params GetBySlugParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get by slug (v1)",
		"slug":    params.Slug,
	})
}
