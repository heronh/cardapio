package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
	"gorm.io/gorm/clause"
)

func NewDish(c *gin.Context) {

	fmt.Println("\nNewDish")
	CompanyId := c.MustGet("CompanyId")
	fmt.Println("CompanyId: ", CompanyId)
	UserId := c.MustGet("ID")
	fmt.Println("UserId: ", UserId)

	// Carrega miniaturas de imagens desta empresa
	var images []models.Image
	if err := initializers.DB.Where("company_id = ?", CompanyId).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Carrega as categorias desta empresa
	var categories []models.Category
	if err := initializers.DB.Where("company_id = ?", CompanyId).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Print all category names
	for _, category := range categories {
		fmt.Println("Category Name: ", category.Name)
	}

	// Copia cada um dos nomes retirando a extensão para uma nova variável
	var imageNames string
	for i, image := range images {
		imageNameWithoutExt := strings.TrimSuffix(image.Original, filepath.Ext(image.Original))
		images[i].Name = imageNameWithoutExt
		if imageNames == "" {
			imageNames = imageNameWithoutExt
		} else {
			imageNames = imageNames + ", " + imageNameWithoutExt
		}
	}

	c.HTML(http.StatusOK, "dish.html", gin.H{
		"Title":      "Novo Prato",
		"CompanyId":  CompanyId,
		"UserId":     UserId,
		"Images":     images,
		"ImageNames": imageNames,
		"Categories": categories,
	})
}

func SaveDish(c *gin.Context) {

	fmt.Println("SaveDish")

	type RequestData struct {
		Name        string `json:"Name"`
		Description string `json:"Description"`
		Price       string `json:"Price"`
		CategoryID  uint   `json:"CategoryID"`
		CompanyId   uint   `json:"CompanyId"`
		UserID      uint   `json:"UserID"`
		WeekDays    []int  `gorm:"type:integer[]"`
		ImageIds    []int  `json:"ImageIds"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	fmt.Println("Name: ", requestData.Name)
	fmt.Println("Description: ", requestData.Description)
	fmt.Println("Price: ", requestData.Price)
	fmt.Println("CategoryID: ", requestData.CategoryID)
	fmt.Println("CompanyId: ", requestData.CompanyId)
	fmt.Println("UserId: ", requestData.UserID)
	fmt.Println("WeekDays: ", requestData.WeekDays)
	fmt.Println("Images: ", requestData.ImageIds)

	daysOfWeek := make([]models.DayOfWeek, len(requestData.WeekDays))
	for i, day := range requestData.WeekDays {
		daysOfWeek[i] = models.DayOfWeek(day)
	}

	dish := models.Dish{
		Name:        requestData.Name,
		Description: requestData.Description,
		Price:       parsePrice(requestData.Price),
		CategoryID:  requestData.CategoryID,
		CompanyID:   requestData.CompanyId,
		UserID:      requestData.UserID,
		DaysOfWeek:  daysOfWeek,
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save new dish and return id
	if err := initializers.DB.Create(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Associate dish to images
	for _, imageId := range requestData.ImageIds {
		dishImage := models.DishImage{
			ImageID: uint(imageId),
			DishID:  dish.ID,
		}
		if err := initializers.DB.Clauses(clause.Returning{}).Create(&dishImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Novo prato criado"})
}

func parsePrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0.0
	}
	return price
}

func DeleteDish(c *gin.Context) {
	fmt.Println("\nDeleteDish")

	type RequestData struct {
		DishID uint `json:"DishID"`
		UserID uint `json:"UserID"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("RequestData: ", requestData)

	if err := initializers.DB.Delete(&models.DishImage{}, "dish_id = ?", requestData.DishID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Delete(&models.Dish{}, requestData.DishID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prato deletado"})
}
