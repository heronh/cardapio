package main

import (
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
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

	// Load HTML templates
	r.LoadHTMLGlob("templates/*.html")

	// Serve static files (CSS) from the 'static' directory
	r.Static("/static", "./static")

	r.POST("/login", login)
	r.GET("/logout", logout)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

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

func logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func login(c *gin.Context) {
	creds := Credentials{}
	fmt.Print(c)
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// creds.Username = c.PostForm("username")
	// creds.Password = c.PostForm("password")
	fmt.Println("Username - ", creds.Username)
	fmt.Println("Password - ", creds.Password)
	if creds.Username != "user" || creds.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
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

		c.Set("username", claims.Username)
		c.Next()
	}
}
