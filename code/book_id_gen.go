package code

import (
	"context"
	"fmt"
	"time"

	"github.com/EmeraldLS/Library_Management_System/db"
	"github.com/EmeraldLS/Library_Management_System/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMaxBookCode() int {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	findOptions := options.Find().SetSort(bson.M{"book_code": -1}).SetLimit(1)
	cursor, _ := db.UsersCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)
	var users []model.User
	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	var maxCode int
	for _, user := range users {
		maxCode = user.UserCode
	}
	return maxCode
}

func GenBookCodeID(book_code int) string {
	prefix := "LSB_BOOK_"
	userID := fmt.Sprintf("%v%d", prefix, book_code)
	return userID
}
