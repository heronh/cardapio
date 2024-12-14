package controllers

import (
	"fmt"
	"net/http"

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

	c.HTML(http.StatusOK, "dish.html", gin.H{
		"Title":     "Novo Prato",
		"CompanyId": CompanyId,
		"UserId":    UserId,
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
