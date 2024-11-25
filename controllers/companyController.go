package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Company(c *gin.Context) {
	c.HTML(http.StatusOK, "company.html", gin.H{
		"Title":          "Empresa",
		"Heading":        "Empresa!",
		"Message":        "Cadastro de empresa",
		"company_active": "h5",
	})
}
