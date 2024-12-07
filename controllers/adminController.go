package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

var dishes []models.Dish = []models.Dish{
	{Name: "Feijoada", Description: "Feijoada mineira completa", Enabled: true, Price: 25.00},
	{Name: "Macarronada", Description: "Macarronada com molho de tomate", Enabled: true, Price: 15.00},
	{Name: "Lasanha", Description: "Lasanha de frango", Enabled: true, Price: 20.00},
	{Name: "Frango grelhado", Description: "Frango grelhado com batata frita", Enabled: true, Price: 18.00},
	{Name: "Bife à parmegiana", Description: "Bife à parmegiana com arroz e batata frita", Enabled: true, Price: 22.00},
}

func Admin(c *gin.Context) {
	fmt.Println("\nAdmin")
	CompanyId := c.MustGet("CompanyId")
	fmt.Println("CompanyId: ", CompanyId)
	UserId := c.MustGet("ID")
	fmt.Println("UserId: ", UserId)

	// carrega os pratos desta empresa
	var dishes []models.Dish
	if err := initializers.DB.Where("company_id = ?", CompanyId).Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"Title":     "Administração",
		"Dishes":    dishes,
		"CompanyId": CompanyId,
		"UserId":    UserId,
	})
}

func CreateDishes(c *gin.Context) {
	fmt.Println("\nCreateDishes")

	type RequestData struct {
		UserID    uint `json:"UserId"`
		CompanyID uint `json:"CompanyId"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("RequestData: ", requestData)

	for _, dish := range dishes {
		dish.UserID = requestData.UserID
		dish.CompanyID = requestData.CompanyID
		dish.CreatedAt = time.Time{}
		dish.UpdatedAt = time.Time{}
		if err := initializers.DB.Create(&dish).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pratos criados"})
}

func CheckUncheckDish(c *gin.Context) {
	fmt.Println("\nCheckUncheckDish")

	type RequestData struct {
		DishID  uint `json:"DishID"`
		UserID  uint `json:"UserID"`
		Enabled bool `json:"Enabled"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("RequestData: ", requestData)

	var dish models.Dish
	dish.UserID = requestData.UserID
	dish.UpdatedAt = time.Time{}
	if err := initializers.DB.First(&dish, requestData.DishID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dish.Enabled = !dish.Enabled
	if err := initializers.DB.Save(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prato atualizado"})
}
