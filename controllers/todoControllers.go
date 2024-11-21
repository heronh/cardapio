package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo
	if err := initializers.DB.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve todos"})
		return
	}

	// retrieve email and user id from the context
	Email, _ := c.Get("email")
	ID, _ := c.Get("ID")
	c.HTML(http.StatusOK, "todo.html", gin.H{
		"Todos": todos,
		"Email": Email,
		"Id":    ID,
	})
}
