package router

import (
	"os"

	"github.com/EmeraldLS/Library_Management_System/controller"
	"github.com/EmeraldLS/Library_Management_System/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {

	var port = os.Getenv("PORT")
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"response": "success",
			"message":  "homepage",
		})
	})
	api := r.Group("api")
	{
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)
		api.GET("/books", controller.GetAllBooks)
		api.GET("/books/:id", controller.GetABook)
		secured := api.Group("/secured")
		{
			secured.Use(middleware.Auth1)
			adminOnly := secured.Group("/admin_only")
			{
				adminOnly.Use(middleware.Auth2)
				adminOnly.POST("/books", controller.InsertBook)
				adminOnly.GET("/users", controller.GetAllUsers)
			}
			secured.POST("/update_password", controller.UpdatePassword)
			secured.POST("/logout", controller.Logout)
		}

	}
	r.Run("0.0.0.0:" + port)
}
