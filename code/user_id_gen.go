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

func GetMaxUserCode() int {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	findOptions := options.Find().SetSort(bson.M{"user_code": -1}).SetLimit(1)
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

func GenUserCodeID(user_code int) string {
	prefix := "LSB_USER_"
	userID := fmt.Sprintf("%v%d", prefix, user_code)
	return userID
}
