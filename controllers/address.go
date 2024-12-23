package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"githum.com/muhammadAslam/ecommerce/database"
	"githum.com/muhammadAslam/ecommerce/models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var userId int64
		if err := c.ShouldBindQuery(&userId); err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to bind query parameter:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind query parameter"})
			return
		}
		address := models.Address{}
		if err := c.ShouldBindJSON(&address); err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to bind JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
			return
		}
		address.UserID = userId
		db := database.Client
		if err := db.WithContext(ctx).Create(&address).Error; err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to create address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusCreated, gin.H{"data": address})
	}
}

func GetAddresses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var userId int64
		if err := c.ShouldBindQuery(&userId); err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to bind query parameter:", err)
			c.JSON(http.StatusNotFound, nil)
			return
		}
		var addresses []models.Address
		db := database.Client
		if err := db.WithContext(ctx).Where("user_id =?", userId).Find(&addresses).Error; err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to fetch addresses:", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"addresses": addresses})
	}
}

func UpdateAddress(id int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Validate and bind userId from query parameters
		var userId int64
		if err := c.ShouldBindQuery(&userId); err != nil || userId == 0 {
			c.Header("Content-Type", "application/json")
			log.Println("Invalid or missing user ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be provided and valid"})
			return
		}

		// Retrieve the existing address based on ID
		var existingAddress models.Address
		db := database.Client
		if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userId).First(&existingAddress).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.Header("Content-Type", "application/json")
				log.Println("Address not found")
				c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			} else {
				c.Header("Content-Type", "application/json")
				log.Println("Failed to retrieve address:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve address"})
			}
			return
		}

		// Bind the incoming JSON payload to the Address struct
		var updatedData models.Address
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to bind JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
			return
		}

		// Update the fields of the existing address
		existingAddress.Street = updatedData.Street
		existingAddress.City = updatedData.City
		existingAddress.State = updatedData.State
		existingAddress.Country = updatedData.Country
		existingAddress.UpdatedAt = time.Now()

		// Save the updated address to the database
		if err := db.WithContext(ctx).Save(&existingAddress).Error; err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to update address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
			return
		}

		// Respond with success
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "data": existingAddress})
	}
}

func DeleteAddress(id int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Validate userId query parameter
		userId := c.Query("user_id")
		if userId == "" {
			c.Header("Content-Type", "application/json")
			log.Println("Please provide a valid user ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid user ID"})
			return
		}

		// Fetch addresses to validate if the user has addresses
		var addresses []models.Address
		db := database.Client
		if err := db.WithContext(ctx).Where("user_id = ?", userId).Find(&addresses).Error; err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to fetch addresses:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
			return
		}

		// Check if any addresses exist for the user
		if len(addresses) == 0 {
			c.Header("Content-Type", "application/json")
			log.Println("No addresses found for the given user ID")
			c.JSON(http.StatusNotFound, gin.H{"error": "No addresses found for the given user ID"})
			return
		}

		// Delete the address by its ID
		if err := db.WithContext(ctx).Where("id = ?", id).Delete(&models.Address{}).Error; err != nil {
			c.Header("Content-Type", "application/json")
			log.Println("Failed to delete address:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
			return
		}

		// Respond with success
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
	}
}
