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
	{Name: "Batata frita", Description: "Deliciosa", Enabled: true, Price: 15.00},
	{Name: "Batata frita com bacon", Description: "Deliciosa e com bacon crocante", Enabled: true, Price: 15.00},
	{Name: "Lasanha", Description: "Lasanha a bolonhesa", Enabled: true, Price: 20.00},
	{Name: "Frango à passarinho", Description: "crocante e com pedaços de alho", Enabled: true, Price: 18.00},
	{Name: "Bife à parmegiana", Description: "Bife à parmegiana com arroz e batata frita", Enabled: true, Price: 22.00},
	{Name: "Creme brulee", Description: "Creme brulee com calda de frutas vermelhas", Enabled: true, Price: 10.00},
	{Name: "Bolo de frutas", Description: "Bolo de frutas com calda de chocolate", Enabled: true, Price: 10.00},
	{Name: "Sorvete de chocolate", Description: "Sorvete de chocolate com calda de chocolate", Enabled: true, Price: 10.00},
	{Name: "Coca-cola", Description: "Coca-cola gelada", Enabled: true, Price: 5.00},
	{Name: "Heineken", Description: "Heineken gelada", Enabled: true, Price: 7.00},
	{Name: "Suco de laranja", Description: "Suco de laranja natural", Enabled: true, Price: 5.00},
	{Name: "Moscow mule", Description: "Moscow mule com limão e gengibre", Enabled: true, Price: 15.00},
	{Name: "Caipirinha", Description: "Caipirinha de limão", Enabled: true, Price: 10.00},
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
