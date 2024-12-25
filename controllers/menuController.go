package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

func Menu(c *gin.Context) {
	fmt.Println("Menu")

	// Parse time stamp from URL and convert to string
	timestamp := c.Query("timestamp")
	var company models.Company
	if err := initializers.DB.Where("Stamp = ?", timestamp).First(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lista os pratos desta empresa
	var dishes []models.Dish
	if err := initializers.DB.Where("company_id = ?", company.ID).Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "menu.html", gin.H{
		"Title":     "Menu",
		"Dishes":    dishes,
		"CompanyId": company.ID,
	})
}
