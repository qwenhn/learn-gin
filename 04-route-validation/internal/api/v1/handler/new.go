package v1Handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/qwenhn/gin-restful-api/04-route-validation/utils"
)

type NewHandler struct {
}

type CreateNewParam struct {
	Title  string `form:"title" binding:"required"`
	Status string `form:"status" binding:"required,oneof=1 2"`
}

func NewNewHandler() *NewHandler {
	return &NewHandler{}
}

func (n *NewHandler) GetAllNews(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Get news (V1)",
			"slug":    "No News",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Get news (V1)",
			"slug":    slug,
		})
	}
}

func (n *NewHandler) CreateNew(ctx *gin.Context) {
	var params CreateNewParam

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	img, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// File less than or equal 5MB
	// 1 << 20 = 1 * 2^20 = 1 * 1048576 = 1MB
	// 5 << 20 = 5 * 2^20 = 5 * 1048576 = 5MB
	if img.Size > 5<<20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File too large (5MB)"})
		return
	}

	// os.ModePerm = 0777 (octal)
	// Grant read, write, and execute permissions to everyone (owner, group, and others)
	path := "./uploads"
	_, err = os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create the upload folder"})
			return
		}
	}

	dst := fmt.Sprintf("%s/%s", path, filepath.Base(img.Filename))
	if err := ctx.SaveUploadedFile(img, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Post news (V1)",
		"title":   params.Title,
		"status":  params.Status,
		"image":   img.Filename,
		"path":    dst,
	})
}
func (n *NewHandler) UploadFileNew(ctx *gin.Context) {
	var params CreateNewParam
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	path := "./uploads"
	filename, err := utils.ValidateAndSaveFile(image, path)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Upload file (V1)",
		"title":   params.Title,
		"status":  params.Status,
		"image":   filename,
		"path":    path + filename,
	})
}
func (n *NewHandler) UploadMultipleFileNew(ctx *gin.Context) {
	const publicURL = "http://localhost:8080/images/"

	var params CreateNewParam
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	images := form.File["images"]
	if len(images) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	var successFiles []string
	var filedFile []map[string]string
	for _, image := range images {
		filename, err := utils.ValidateAndSaveFile(image, "./uploads")
		if err != nil {
			filedFile = append(filedFile, map[string]string{
				"filename": image.Filename,
				"error":    err.Error(),
			})

			continue
		}

		publicImageURL := publicURL + filename
		successFiles = append(successFiles, publicImageURL)
	}

	resp := gin.H{
		"message":       "Upload multiple file (V1)",
		"title":         params.Title,
		"status":        params.Status,
		"success_files": successFiles,
	}

	if len(filedFile) > 0 {
		resp["message"] = "Upload completed with partial errors"
		resp["error_files"] = filedFile
	}

	ctx.JSON(http.StatusOK, resp)
}
