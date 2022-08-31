package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"startup-backend-api/auth"
	"startup-backend-api/campaign"
	"startup-backend-api/handler"
	"startup-backend-api/images/helper"
	"startup-backend-api/payment"
	"startup-backend-api/transaction"
	"startup-backend-api/user"
	"strings"

	webHandler "startup-backend-api/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	//===============Connection-Database===============//

	dsn := "root:@tcp(127.0.0.1:3306)/dbstartup?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//=================User-Endpoint===================//

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	//==============Jwt-Token-Validasi=================//

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.9njU5wOeHW1HCFnZsrSbfaWcHQI6Iei36ZrxQSyiXiA")

	if err != nil {
		fmt.Println("ERROR")
	}

	if token.Valid {
		fmt.Println("Valid")
	} else {
		fmt.Println("Invalid")
	}

	//=================Campaign-Endpoint===================//

	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	//================Transaction-Endpoint=================//

	transactionRepository := transaction.NewRepository(db)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	//===================User-Web-Handler==================//

	dashboardWebHandler := webHandler.NewDashboardHandler()
	userWebHandler := webHandler.NewUserHandler(userService)

	//=================Router-And-List-API=================//

	router := gin.Default()

	// router.Use(cors.Default())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "PATCH", "POST", "DELETE", "PUT", "OPTIONS", "HEAD"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// router.LoadHTMLGlob("web/templates/**/*")

	router.HTMLRender = loadTemplates("./web/templates") // router for web cms

	//========Load-Static-Folder=======//

	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/extensions", "./web/assets/extensions")
	router.Static("/fonts", "./web/assets/fonts")
	router.Static("/image", "./web/assets/image")
	router.Static("/js", "./web/assets/js")

	api := router.Group("/api/v1")

	//=============Users===============//

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	//===========Campaign=============//

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	//===========Transaction==========//

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	//=======Router-Web-Dashboard======//

	router.GET("/", dashboardWebHandler.Index)

	//============User-Web============//

	router.GET("/users/", userWebHandler.Index)
	router.GET("/users/add/", userWebHandler.Add)
	router.POST("/users/", userWebHandler.Create)

	//=========Run-Port-8080==========//

	router.Run(":8080")
}

//===================Auth-Middleware=================//

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}

}

//===================Load-Template-Html=================//

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
