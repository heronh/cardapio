package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/initializers"
	"github.com/heronh/cardapio/models"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func CompanySave(c *gin.Context) {

	fmt.Println("CompanySave")

	var company = models.Company{}
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save the company to the database
	fmt.Println(company)
	if err := initializers.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save company"})
		return
	}
	// Send a response
	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully", "CompanyId": company.ID})
}

func Company(c *gin.Context) {
	now := time.Now()
	email := "heron" + fmt.Sprint(now.Unix()) + "@gmail.com"
	name := "Heron" + fmt.Sprint(now.Unix())
	c.HTML(http.StatusOK, "company.html", gin.H{
		"Title":          "Empresa",
		"Heading":        "Empresa!",
		"Message":        "Cadastro de empresa",
		"company_active": "h5",
		"email":          email,
		"name":           name,
	})
}
