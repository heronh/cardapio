package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

func UncheckTodo(c *gin.Context) {
	type RequestData struct {
		Id int `json:"Id"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Unchecking todo with id:", requestData.Id)
	var todo models.Todo
	if err := initializers.DB.Where("id = ?", requestData.Id).First(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find todo"})
		return
	}

	todo.Completed = false
	todo.Updated_at = time.Now()
	if err := initializers.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully unchecked todo"})
}

func CheckTodo(c *gin.Context) {
	type RequestData struct {
		Id int `json:"Id"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Checking todo with id:", requestData.Id)
	var todo models.Todo
	if err := initializers.DB.Where("id = ?", requestData.Id).First(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find todo"})
		return
	}

	todo.Completed = true
	todo.Updated_at = time.Now()
	if err := initializers.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully checked todo"})
}

func DeleteTodo(c *gin.Context) {

	type RequestData struct {
		Id int `json:"Id"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Deleting todo with id:", requestData.Id)
	if err := initializers.DB.Delete(&models.Todo{}, requestData.Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted todo"})
}

func SaveTodo(c *gin.Context) {
	fmt.Println("Creating todo")
	var todo models.Todo
	todo.Created_at = time.Now()
	todo.Updated_at = time.Now()
	todo.Completed = false
	todo.Description = c.PostForm("description")

	fmt.Println(c)
	Id := c.PostForm("Id")
	var userModel models.User
	if err := initializers.DB.Where("id = ?", Id).First(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
		return
	}

	fmt.Println("Todo fields:")
	fmt.Println("Description:", todo.Description)
	fmt.Println("CreatedAt:", todo.Created_at)
	fmt.Println("UpdatedAt:", todo.Updated_at)
	fmt.Println("Completed:", todo.Completed)

	if err := initializers.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}
	c.Redirect(http.StatusFound, "/todos")
}

func GetTodos(c *gin.Context) {

	var todos []models.Todo
	if err := initializers.DB.Order("completed, updated_at desc").Find(&todos).Error; err != nil {
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
