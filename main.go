package main

import (
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	// "github.com/heronh/cardapio/controllers/todoControllers"
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

	// Load HTML templates
	r.LoadHTMLGlob("templates/*.html")

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

	// r.GET("/todos", authMiddleware(), todoControllers.GetTodos)

	// Salva nova tarefa no banco de dados
	r.POST("/todos", save_todo)

	// Apaga tarefa do banco de dados
	r.POST("/todos_delete", todos_delete)

	// read port in .env file and starts the server
	port := os.Getenv("HostPort")
	if port == "" {
		port = "8080" // default port if not specified
	}
	r.Run(":" + port)
}

func save_todo(c *gin.Context) {
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

func todos_delete(c *gin.Context) {

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
