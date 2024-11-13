package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var host string
var port string
var user string
var password string
var dbname string
var dsn string

var DB *gorm.DB

func loadEnv() {
	host = os.Getenv("host")
	port = os.Getenv("port")
	user = os.Getenv("user")
	password = os.Getenv("password")
	dbname = os.Getenv("dbname")
	fmt.Println("host", host)
	fmt.Println(host, port, user, password, dbname)
	dsn = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
}

func ConnectToDb() {
	// Connect to the database
	fmt.Println("Connecting to the database")
	var err error
	loadEnv()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
