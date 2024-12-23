package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"githum.com/muhammadAslam/ecommerce/database"
)

type Application struct {
	ProductData *database.ProductData
	UserData    *database.UserData
}

func NewApplication(productData *database.ProductData, userData *database.UserData) *Application {
	return &Application{
		ProductData: productData,
		UserData:    userData,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add product to cart
		productQueryById := c.Query("id")
		if productQueryById == "" {
			log.Println("Product not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product not found"))
			return
		}
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("User not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User not found"))
			return
		}
		productId, err := strconv.Atoi(productQueryById)
		if err != nil {
			log.Println("Invalid product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid product ID"))
			return
		}
		userId, err := strconv.Atoi(userQueryId)
		if err != nil {
			log.Println("Invalid user ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID"))
			return
		}
		if _, err := app.ProductData.GetById(int64(productId)); err != nil {
			log.Println("Product not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product not found"))
			return
		}
		err = database.AddProductToCart(c.Request.Context(), app.ProductData.DB, int64(userId), int64(productId))
		if err != nil {
			log.Println("Failed to add product to cart")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to add product to cart"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})

	}
}

func (app *Application) RemoveFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Remove product from Cart List
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("Product not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product not found"))
			return
		}
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("User not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User not found"))
			return
		}
		productId, err := strconv.Atoi(productQueryId)
		if err != nil {
			log.Println("Invalid product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid product ID"))
			return
		}
		userId, err := strconv.Atoi(userQueryId)
		if err != nil {
			log.Println("Invalid user ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID"))
			return
		}
		err = database.RemoveProductFromCart(c.Request.Context(), app.ProductData.DB, int64(userId), int64(productId))
		if err != nil {
			log.Println("Failed to remove product from cart")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to remove product from cart"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})

	}
}

func (app *Application) GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cart list of products for a user
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("User not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User not found"))
			return
		}
		userId, err := strconv.Atoi(userQueryId)
		if err != nil {
			log.Println("Invalid user ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID"))
			return
		}
		cartItems, err := database.GetCartItems(c.Request.Context(), app.ProductData.DB, int64(userId))
		if err != nil {
			log.Println("Failed to get cart items")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to get cart items"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"cartItems": cartItems})
	}
}

func (app *Application) Checkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Checkout cart items for a user
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("User not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User not found"))
			return
		}
		userId, err := strconv.Atoi(userQueryId)
		if err != nil {
			log.Println("Invalid user ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID"))
			return
		}
		cartItems, err := database.GetCartItems(c.Request.Context(), app.ProductData.DB, int64(userId))
		if err != nil {
			log.Println("Failed to get cart items")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to get cart items"))
			return
		}
		totalPrice := 0.0
		for _, item := range cartItems {
			totalPrice += item.Price * float64(item.Quantity)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart checked out", "price": totalPrice, "data": cartItems})
	}
}
func (app *Application) GetInstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get instant buy product details for a user
		productQueryById := c.Query("id")
		if productQueryById == "" {
			log.Println("Product not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product not found"))
			return
		}
		productId, err := strconv.Atoi(productQueryById)
		if err != nil {
			log.Println("Invalid product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid product ID"))
			return
		}

		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("User not found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User not found"))
			return
		}
		userId, err := strconv.Atoi(userQueryId)
		if err != nil {
			log.Println("Invalid user ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid user ID"))
			return
		}
		product := database.GetInstantBuyProduct(c.Request.Context(), app.ProductData.DB, int64(userId), int64(productId))
		if err != nil {
			log.Println("Failed to get instant buy product")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to get instant buy product"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"product": product})
	}
}
