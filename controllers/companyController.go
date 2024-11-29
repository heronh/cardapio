package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heronh/cardapio/models"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func CompanySave(c *gin.Context) {

	fmt.Println("CompanySave")

	/*
		var companyData map[string]interface{}
		if err := c.ShouldBindJSON(&companyData); err != nil {
			fmt.Println("error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for key, value := range companyData {
			fmt.Printf("%s: %v (type: %T)\n", key, value, value)
		}
		fmt.Printf("name: %v\n", companyData["name"])
	*/

	var company = models.Company{}
	if err := c.ShouldBindJSON(&company); err != nil {
		fmt.Println("error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send a response
	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully"})
}

func Company(c *gin.Context) {
	c.HTML(http.StatusOK, "company.html", gin.H{
		"Title":          "Empresa",
		"Heading":        "Empresa!",
		"Message":        "Cadastro de empresa",
		"company_active": "h5",
	})
}
