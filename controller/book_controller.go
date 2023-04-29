package controller

import (
	"fmt"
	"net/http"

	"github.com/EmeraldLS/Library_Management_System/code"
	"github.com/EmeraldLS/Library_Management_System/config"
	"github.com/EmeraldLS/Library_Management_System/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon"
)

func InsertBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	validate := validator.New()
	err := validate.Struct(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	book.BookCode = code.GetMaxBookCode() + 1
	book.BookID = code.GenBookCodeID(book.BookCode)
	book.Registered_At = carbon.Now().ToDateTimeString()
	book.Updated_At = carbon.Now().ToDateTimeString()

	if err := config.InsertOneBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": "success",
		"message":  book,
	})
}

func GetAllBooks(c *gin.Context) {
	books, err := config.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetABook(c *gin.Context) {
	bookID := c.Param("book_id")
	book, err := config.GetABook(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  err,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, book)
}
