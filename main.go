package main

import (
	"log"
	"net/http"
	"path/filepath"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/heronh/cardapio/controllers"
	"github.com/heronh/cardapio/initializers"
)

func init() {
	// Load the environment variables
	initializers.LoadEnv()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
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
			"Title":   "Benvindo",
			"Heading": "Página de acesso!",
			"Message": "",
			"welcome": "h5",
		})
	})

	r.GET("/welcome", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	})

	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.GET("register", controllers.Register)
	r.POST("/register", controllers.New_user)
	r.GET("/login", controllers.LoginPage)

	// Funções relativas ao cadastro de empresas
	r.GET("/company", controllers.Company)
	r.POST("/user-check-email", controllers.CheckEmail)
	r.POST("/company-save", controllers.CompanySave)

	// Funções relativas as tarefas
	r.GET("/todos", authMiddleware(), controllers.GetTodos)
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			c.Abort()
			return
		}

		tokenStr := cookie.Value
		claims := &controllers.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return controllers.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Redirect(http.StatusFound, "/login?error=unauthorized")
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.Redirect(http.StatusFound, "/login?error=unauthorized")
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("ID", claims.Id)
		c.Next()
	}
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
