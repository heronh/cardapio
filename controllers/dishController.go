package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
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
