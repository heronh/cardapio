package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

func CreateCategory(c *gin.Context) {

	fmt.Println("\nCreateCategory")

	// Parse incoming JSON data
	/*
		var jsonData map[string]interface{}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Incoming JSON data: ", jsonData)
	*/

	type RequestData struct {
		UserID    uint   `json:"UserId"`
		CompanyID uint   `json:"CompanyId"`
		Name      string `json:"Name"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("RequestData: ", requestData)

	// Create a new category
	category := models.Category{
		UserID:    requestData.UserID,
		CompanyID: requestData.CompanyID,
		Name:      requestData.Name,
	}

	// Check if the category already exists
	var existingCategory models.Category
	result := initializers.DB.Where("company_id = ? AND name = ?", requestData.CompanyID, requestData.Name).First(&existingCategory)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Categoria já existe!"})
		return
	}

	result = initializers.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusNotModified, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nova categoria criada!", "category": category.Name})
}

func GetCategories(c *gin.Context) {
	fmt.Println("\nGetCategories")

	// Parse incoming JSON data
	/*
		var jsonData map[string]interface{}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Incoming JSON data: ", jsonData)
	*/

	type RequestData struct {
		CompanyID uint `json:"CompanyId"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("RequestData: ", requestData)

	// Get all categories
	var categories []models.Category
	result := initializers.DB.Where("company_id = ?", requestData.CompanyID).Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusNotModified, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
