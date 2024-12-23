package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"githum.com/muhammadAslam/ecommerce/controllers"
	"githum.com/muhammadAslam/ecommerce/database"
	"githum.com/muhammadAslam/ecommerce/middleware"
	"githum.com/muhammadAslam/ecommerce/routes"
)

func main() {
	// Print to indicate the application started
	fmt.Println("Hello, World!")

	// Retrieve the port from the environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize the application with products and users data
	productData := database.NewProductData(database.DBSet(), "Products")
	userData := database.NewUserData(database.DBSet(), "Users")

	app := controllers.NewApplication(productData, userData)
	// Initialize Gin router
	router := gin.New()

	// Use built-in Gin logger middleware
	router.Use(gin.Logger())
	router.GET("/get-products", controllers.GetProds())
	router.GET("/get-product/:id", controllers.GetProductByID())
	router.GET("/search-products", controllers.SearchProduct())
	router.POST("/signup", controllers.Signup())
	router.POST("/login", controllers.Login())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the E-commerce API",
		})
	})

	// Define user routes

	// Apply authentication middleware globally for all routes after this point
	router.Use(middleware.Authentication())
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	// Define other routes for the app

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removefromcart", app.RemoveFromCart())
	router.GET("/cartcheckout", app.Checkout()) // Fixed the path typo: "cartcheckput" -> "cartcheckout"
	router.GET("/instantbuy", app.GetInstantBuy())

	// Start the server on the specified port
	log.Fatal(router.Run(":" + port))
}
