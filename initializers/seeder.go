package initializers

import (
	"fmt"
	"time"

	"github.com/heronh/cardapio/models"
	"golang.org/x/crypto/bcrypt"
)

func Seeder() {

	// check if users table is empty, if not, return
	var count int64
	if DB.Model(&models.User{}).Count(&count); count != 0 {
		fmt.Println("Usuários já criados")
		return
	}
	fmt.Println("Criando primeiro usuário e primeira empresa")

	// create a company
	var company = models.Company{}
	company.Name = "Empresa 1"
	company.Description = "Empresa 1 - Descrição"
	company.Category = "Restaurante"
	company.Address = "Rua 1"
	company.Street = "Rua 1"
	company.Number = "1"
	company.Complement = "Complemento 1"
	company.City = "Cidade 1"
	company.State = "Estado 1"
	company.Country = "País 1"
	company.PostalCode = "12345-678"
	company.Phone = "1234-5678"
	company.Website = "www.empresa1.com"
	company.UpdatedAt = time.Now()
	company.CreatedAt = time.Now()

	// save the company to the database and return id
	DB.Create(&company)

	// create a user
	var user = models.User{}
	user.Name = "Heron"
	user.Email = "heron@gmail.com"
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	user.Password = "123456"
	user.Company = company
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	// save the user to the database
	DB.Create(&user)
}
