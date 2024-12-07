package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email     string `json:"email"`
	Id        uint   `json:"id"`
	CompanyId uint   `json:"company_id"`
	jwt.RegisteredClaims
}

var JwtKey = []byte("my_secret_key")

func UserSave(c *gin.Context) {

	fmt.Println("UserSave")
	var user = models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	fmt.Println("User: ", user)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})
		return
	}
	user.Password = string(hash)

	// save the user to the database
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save user"})
		return
	}

	// Retrieve the ID of the saved user
	if err := initializers.DB.Last(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user ID"})
		return
	}

	// Send a response
	fmt.Println(user)
	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully"})
}

func New_user(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"Title":           "Cadastro",
		"Heading":         "Cadastro!",
		"Message":         "Cadastro de usuário",
		"register_active": "h5",
	})
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func LoginPage(c *gin.Context) {

	// Parse email from URL
	email := c.Query("email")
	fmt.Println("Email: ", email)

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title":   "Login",
		"Heading": "Login!",
		"Message": "Página de login",
		"login":   "h5",
		"email":   email,
	})
}

func Login(c *gin.Context) {
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
	fmt.Println("User: ", user)
	fmt.Println("Company: ", user.Company)
	fmt.Println("Company ID: ", user.CompanyId)
	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &Claims{
		Email:     creds.Email,
		Id:        user.ID,
		CompanyId: user.CompanyId,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtKey)
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

func CheckEmail(c *gin.Context) {
	email := c.Query("email")
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário encontrado", "user": user})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				redirectUnauthorized(c)
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
			return JwtKey, nil
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
		c.Set("CompanyId", claims.CompanyId)
		c.Next()
	}
}

func redirectUnauthorized(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login?error=unauthorized")
	c.Abort()
}
