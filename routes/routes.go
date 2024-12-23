package routes

import (
	"github.com/gin-gonic/gin"
	"githum.com/muhammadAslam/ecommerce/controllers"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/admin/add-products", controllers.AddProduct())
	incomingRoutes.GET("/admin/get-products", controllers.GetProducts())
	incomingRoutes.GET("/admin/get-product/:id", controllers.GetProductByID())
	incomingRoutes.PUT("/admin/update-product/:id", controllers.UpdateProduct())
	incomingRoutes.DELETE("/admin/delete-product/:id", controllers.DeleteProduct())
}
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User Deatil Api",
		})

	})
	incomingRoutes.PATCH("/user/update", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User Updated Successfully",
		})
	})
}
