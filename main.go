package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/heronh/cardapio/controllers"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	Id    uint   `json:"id"`
	jwt.RegisteredClaims
}

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

	r.POST("/login", login)
	r.GET("/logout", logout)
	r.GET("register", register)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":   "Benvindo",
			"Heading": "Página de acesso!",
			"Message": "",
			"welcome": "h5",
		})
	})

	r.POST("/register", new_user)
	r.GET("/welcome", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/company", func(c *gin.Context) {
		c.HTML(http.StatusOK, "company.html", gin.H{
			"Title":   "Company",
			"Heading": "Company Information",
		})
	})

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

func new_user(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})
		return
	}
	user.Password = string(hash)

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário criado com sucesso"})
}

func register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"Title":           "Cadastro",
		"Heading":         "Cadastro!",
		"Message":         "Cadastro de usuário",
		"register_active": "h5",
	})
}

func logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func login(c *gin.Context) {
	fmt.Println(c)
	creds := Credentials{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})
		return
	}
	creds.Password = string(hash)

	// read data from database
	var user models.User
	if err := initializers.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não cadastrado"})
		return
	}

	expirationTime := time.Now().Add(50 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		Id:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
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
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
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
