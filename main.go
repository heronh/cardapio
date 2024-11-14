package main

import (
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

func init() {
	// Load the environment variables
	initializers.LoadEnv()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {

	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*.html")

	// Serve static files (CSS) from the 'static' directory
	r.Static("/static", "./static")

	r.GET("/welcome", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"Title":   "Benvindo",
			"Heading": "Página de acesso!",
			"Message": "",
			"welcome": "h5",
		})
	})

	// Lê banco de dados e lista tarefas
	r.GET("/todos", func(c *gin.Context) {
		var todos []models.Todo
		if err := initializers.DB.Find(&todos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve todos"})
			return
		}
		// fmt.Println("Todos:")
		// for i := range todos {
		// 	fmt.Println(todos[i].Description)
		// }
		c.HTML(http.StatusOK, "todo.html", gin.H{"Todos": todos})
	})

	// Salva nova tarefa no banco de dados
	r.POST("/todos", func(c *gin.Context) {
		fmt.Println("Creating todo")
		var todo models.Todo
		todo.Created_at = time.Now()
		todo.Updated_at = time.Now()
		todo.Completed = false
		todo.Description = c.PostForm("description")

		fmt.Println("Todo fields:")
		fmt.Println("Description:", todo.Description)
		fmt.Println("CreatedAt:", todo.Created_at)
		fmt.Println("UpdatedAt:", todo.Updated_at)
		fmt.Println("Completed:", todo.Completed)

		if err := initializers.DB.Create(&todo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
			return
		}
		c.JSON(http.StatusCreated, todo)
	})

	// read port in .env file and starts the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}
	r.Run(":" + port)
}
