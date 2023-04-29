package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/EmeraldLS/Library_Management_System/db"
	"github.com/EmeraldLS/Library_Management_System/model"
	"github.com/EmeraldLS/Library_Management_System/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckUserType(userType string) bool {
	return userType == "ADMIN"
}

func Auth1(c *gin.Context) {
	tokenString := c.GetHeader("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "error",
			"message":  "No token provided",
		})
		c.Abort()
		return
	}
	err := token.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
}

func Auth2(c *gin.Context) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	users := db.UsersCollection
	tokenString := c.GetHeader("token")
	err := users.FindOne(ctx, bson.M{"token": tokenString}).Decode(&user)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("No user with Token is found. Error(%v)", err),
		})
		c.Abort()
		return
	}
	if !CheckUserType(user.UserType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "error",
			"message":  "User doesn't have access",
		})
		c.Abort()
		return
	}
}
