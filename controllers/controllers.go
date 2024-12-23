package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"githum.com/muhammadAslam/ecommerce/database"
	"githum.com/muhammadAslam/ecommerce/models"
	"githum.com/muhammadAslam/ecommerce/tokens"
	"golang.org/x/crypto/bcrypt"
)

var UserData *database.UserData = database.NewUserData(database.DBSet(), "users")
var ProductData *database.ProductData = database.NewProductData(database.DBSet(), "products")
var validate = validator.New()

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}
func VerifyPassword(userPassword string, givenPassword string) (bool, string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, "Invalid password", nil
	} else if err != nil {
		return false, "", err
	}
	return true, "", nil
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// User struct to hold incoming data
		var user models.User

		// Bind incoming JSON to the user struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the user struct
		if validationerror := validate.Struct(user); validationerror != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationerror.Error()})
			return
		}

		// Check if the user already exists in the database
		var existUser models.User
		db := database.Client // Assuming database.Client is the initialized gorm.DB instance
		if err := db.Where("email = ?", user.Email).First(&existUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		if err := db.Where("phone = ?", user.Phone).First(&existUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exits with this phone"})
			return
		}
		defer cancel()

		// Hash the user's password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		token, refreshToken, _ := tokens.GenerateAllTokens(db, user.Email, user.Name)
		user.Token = token
		user.RefreshToken = refreshToken
		user.Roles = "user"

		// Save the new user to the database
		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Return success response
		c.JSON(http.StatusCreated, gin.H{"data": user})
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginUser models.User
		if err := c.ShouldBindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		db := database.Client
		fmt.Println("emails:", loginUser.Email)
		var storedUser models.User
		if err := db.WithContext(ctx).Where("email = ?", loginUser.Email).First(&storedUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": err.Error()})
			return
		}
		passwordIsValid, msg, err := VerifyPassword(storedUser.Password, loginUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify password"})
			return
		}
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			fmt.Println("password", loginUser.Password)
			return
		}
		token, refreshToken, _ := tokens.GenerateAllTokens(db, storedUser.Email, storedUser.Name)
		storedUser.Token = token
		storedUser.RefreshToken = refreshToken
		tokens.UpdateAllTokens(db, token, refreshToken, storedUser.ID)
		c.JSON(http.StatusOK, gin.H{"data": storedUser})

	}
}

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		db := database.Client
		if err := db.WithContext(ctx).Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": product})
	}
}

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var productList []models.Product
		db := database.Client
		db.WithContext(ctx).Find(&productList)
		if len(productList) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"error": "No products found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"products": productList})

	}
}

func GetProductByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		id := c.Param(":id")
		var product models.Product
		db := database.Client
		db.WithContext(ctx).Where("id = ?", id).Find(&product)
		if product.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"product": product})
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		id := c.Param("id")
		var product models.Product
		db := database.Client
		db.WithContext(ctx).Where("id =?", id).Delete(&product)
		if product.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

func GetProds() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var productList []models.Product
		db := database.Client
		db.WithContext(ctx).Find(&productList)
		c.JSON(http.StatusOK, gin.H{"data": productList})
	}
}

func GetProdById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		id := c.Param("id")
		var product models.Product
		db := database.Client
		db.WithContext(ctx).Where("id =?", id).First(&product)
		if product.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": product})

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		product := c.Query("product")
		db := database.Client
		db.Where("name LIKE?", "%"+product+"%").Find(&productList)
		c.JSON(http.StatusOK, gin.H{"data": productList})

	}
}
