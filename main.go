package main

import (
	"log"
	"net/http"
	"path/filepath"

	"os"

	"github.com/gin-gonic/gin"

	"github.com/heronh/cardapio/controllers"
	"github.com/heronh/cardapio/initializers"
)

func init() {
	// Load the environment variables
	initializers.LoadEnv()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.Seeder()
}

func main() {

	r := gin.Default()

	// Serve template files and it's subfolders from the 'templates' directory
	// load templates recursively
	files, err := loadTemplates("templates")
	if err != nil {
		log.Println(err)
	}
	r.LoadHTMLFiles(files...)

	// Serve static files (CSS) from the 'static' directory
	r.Static("/static", "./static")
	// Serve Bootstrap icons from the 'node_modules/bootstrap-icons' directory
	r.Static("/icons", "./static/bootstrap-icons")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":   "Benvindo ao seu próximo cardápio digital!",
			"Heading": "Página de acesso!",
			"Message": "Benvindo ao seu próximo cardápio digital!",
			"welcome": "h5",
		})
	})

	r.GET("/welcome", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	})

	r.GET("/images", controllers.AuthMiddleware(), controllers.Images)
	r.POST("/images/upload", controllers.Upload)

	r.GET("/admin", controllers.AuthMiddleware(), controllers.Admin)
	r.POST("/create-dishes", controllers.CreateDishes)
	r.POST("/admin/check-uncheck-dish", controllers.CheckUncheckDish)
	r.POST("/admin/delete-dish", controllers.DeleteDish)
	r.GET("/dish", controllers.AuthMiddleware(), controllers.NewDish)

	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.GET("register", controllers.Register)
	r.POST("/register", controllers.New_user)
	r.GET("/login", controllers.LoginPage)
	r.POST("/user-save", controllers.UserSave)

	// Funções relativas ao cadastro de empresas
	r.GET("/company", controllers.Company)
	r.POST("/user-check-email", controllers.CheckEmail)
	r.POST("/company-save", controllers.CompanySave)

	// Funções relativas as tarefas
	r.GET("/todos", controllers.AuthMiddleware(), controllers.GetTodos)
	r.POST("/todos", controllers.SaveTodo)
	r.POST("/todos_delete", controllers.DeleteTodo)
	r.POST("/todos_check", controllers.CheckTodo)
	r.POST("/todos_uncheck", controllers.UncheckTodo)

	// read port in .env file and starts the server
	port := os.Getenv("HostPort")
	if port == "" {
		port = "8080" // default port if not specified
	}
	r.Run(":" + port)
}

func loadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				loadTemplates(path)
			}
		} else {
			files = append(files, path)
		}
		return err
	})
	return files, err
}
