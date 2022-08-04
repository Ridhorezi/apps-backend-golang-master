package main

import (
	"fmt"
	"log"
	"startup-backend-api/auth"
	"startup-backend-api/handler"
	"startup-backend-api/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	//==============Connection-Database==============//

	dsn := "root:@tcp(127.0.0.1:3306)/dbstartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//================User-Endpoint==================//

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.9njU5wOeHW1HCFnZsrSbfaWcHQI6Iei36ZrxQSyiXiA")

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("ERROR")
		fmt.Println("ERROR")
	}

	if token.Valid {
		fmt.Println("Valid")
		fmt.Println("Valid")
		fmt.Println("Valid")
	} else {
		fmt.Println("Invalid")
		fmt.Println("Invalid")
		fmt.Println("Invalid")
	}

	userHandler := handler.NewUserHandler(userService, authService)

	//=============Router-And-List-API===============//

	router := gin.Default()

	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run(":8080")
}
