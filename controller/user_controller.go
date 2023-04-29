package controller

import (
	"fmt"
	"net/http"

	"github.com/EmeraldLS/Library_Management_System/code"
	"github.com/EmeraldLS/Library_Management_System/config"
	"github.com/EmeraldLS/Library_Management_System/email"
	"github.com/EmeraldLS/Library_Management_System/model"
	"github.com/EmeraldLS/Library_Management_System/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(providedPassword string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(providedPassword), 14)
	if err != nil {
		fmt.Println(err)
	}
	return string(byte)
}

func Login(c *gin.Context) {
	var user model.LoginStruct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	if err := config.Login(user.UserID, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	usertoken, refreshToken, expirationTime, err := token.GenerateToken(user.UserID, user.Name, user.Email, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "token error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	if err := token.UpdateAllToken(usertoken, refreshToken, expirationTime, user.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "update token error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response":  "success",
		"message":   "Login successful",
		"new token": usertoken,
	})
}

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "json error",
			"message":  err,
		})
		c.Abort()
		return
	}
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "validation err",
			"message":  fmt.Sprintf("An error occured with the data passed. Error(%v)", err),
		})
		c.Abort()
		return
	}
	newCode := code.GetMaxUserCode() + 1
	user.UserCode = newCode
	user.UserID = code.GenUserCodeID(user.UserCode)
	user.Password = HashPassword(user.Password)
	user.Registered_At = carbon.Now().ToDateTimeString()
	user.Updated_At = carbon.Now().ToDateTimeString()
	token, refreshToken, expirationTime, err := token.GenerateToken(user.UserID, user.Name, user.Email, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "token error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	user.ExpirationTime = expirationTime
	user.Token = token
	user.Refresh_Token = refreshToken
	if err := email.ValidateEmail(user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "email error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	err = config.InsertOneUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": "success",
		"message":  user,
	})

}

func Logout(c *gin.Context) {
	tokenString := c.GetHeader("token")
	config.SetTokenExpired(tokenString)
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  "Logout successful",
	})
}

func GetAllUsers(c *gin.Context) {
	users, err := config.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  users,
	})
}

func UpdatePassword(c *gin.Context) {
	token := c.GetHeader("token")
	var user model.UpdatePassword
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	fmt.Println(user)
	if err := config.UpdatePassword(token, user.OldPassword, user.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  "password updated successfully.",
	})
}
