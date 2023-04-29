package config

import (
	"context"
	"time"

	"github.com/EmeraldLS/Library_Management_System/db"
	"github.com/EmeraldLS/Library_Management_System/model"
	"go.mongodb.org/mongo-driver/bson"
)

var booksCollection = db.BooksCollection

func InsertOneBook(book model.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := booksCollection.InsertOne(ctx, book)
	if err != nil {
		return err
	}
	return nil
}

func GetAllBooks() ([]model.Book, error) {
	var books []model.Book
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cursor, err := booksCollection.Find(ctx, bson.D{{}})

	if err != nil {
		return []model.Book{}, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book model.Book
		_ = cursor.Decode(&book)
		books = append(books, book)
	}
	return books, err
}

func GetABook(bookID string) (model.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var book model.Book
	filter := bson.M{"book_id": bookID}
	if err := booksCollection.FindOne(ctx, filter).Decode(&book); err != nil {
		return model.Book{}, err
	}
	return book, nil
}
