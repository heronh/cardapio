package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

func Images(c *gin.Context) {
	fmt.Println("\nImages")
	CompanyId := c.MustGet("CompanyId")
	fmt.Println("CompanyId: ", CompanyId)
	UserId := c.MustGet("ID")
	fmt.Println("UserId: ", UserId)

	// carrega as imagens desta empresa
	var images []models.Image
	if err := initializers.DB.Where("company_id = ?", CompanyId).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "images.html", gin.H{
		"Title":     "Imagens",
		"Images":    images,
		"CompanyId": CompanyId,
		"UserId":    UserId,
	})
}

func Upload(c *gin.Context) {

	fmt.Println("\nUpload")
	CompanyIdStr := c.PostForm("CompanyId")
	CompanyId, err := strconv.ParseUint(CompanyIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CompanyId"})
		return
	}
	fmt.Println("CompanyId: ", CompanyId)

	UserIdStr := c.PostForm("UserId")
	UserId, err := strconv.ParseUint(UserIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserId"})
		return
	}
	fmt.Println("UserId: ", UserId)

	// Save the image to the database
	image := models.Image{
		UserID:    uint(UserId),
		CompanyID: uint(CompanyId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsSample:  false,
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	files := form.File["images[]"]

	// Save the files to the server
	for _, file := range files {
		fmt.Println("file: ", file.Filename)
		filename := fmt.Sprintf("_%s_%d_%s", CompanyIdStr, time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, fmt.Sprintf("static/images/%s", filename)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		image_2_be_saved := image
		image_2_be_saved.Name = filename
		if err := initializers.DB.Create(&image_2_be_saved).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.Redirect(http.StatusFound, "/images")
}
