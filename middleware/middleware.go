package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"githum.com/muhammadAslam/ecommerce/tokens"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := tokens.ValidateToken(ClientToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		fmt.Printf("Valid token! Email: %s, Name: %s\n", claims.Email, claims.Name)
		c.Next()
	}
}
