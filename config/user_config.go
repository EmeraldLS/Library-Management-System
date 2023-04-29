package config

import (
	"context"
	"fmt"
	"time"

	"github.com/EmeraldLS/Library_Management_System/db"
	"github.com/EmeraldLS/Library_Management_System/model"
	"github.com/golang-module/carbon"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection = db.UsersCollection

func InsertOneUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	filter := bson.M{}
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	cursor, err := usersCollection.Find(ctx, filter)
	if err != nil {
		return []model.User{}, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	return users, nil
}

func SetTokenExpired(token string) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	filter := bson.M{"token": token}
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}
	var updateObj = bson.M{}
	user.ExpirationTime = time.Now().Unix()
	updateObj["expiration_time"] = user.ExpirationTime
	updateDetail := bson.M{"$set": updateObj}
	_, err = usersCollection.UpdateOne(ctx, filter, updateDetail)

	if err != nil {
		fmt.Println(err)
	}

}

func VerifyPassword(providedPassword string, hashedPwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(providedPassword)); err != nil {
		return err
	}

	return nil
}

func HashPassword(providedPassword string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(providedPassword), 14)
	if err != nil {
		fmt.Println(err)
	}
	return string(byte)
}

func Login(userId string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var user model.LoginStruct
	filter := bson.M{"user_id": userId}
	if err := usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return fmt.Errorf("user with user id is not found")
	}
	if err := VerifyPassword(password, user.Password); err != nil {
		return fmt.Errorf("incorrect password")
	}
	return nil

}

func UpdatePassword(token string, oldPassword string, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var user model.User
	filter := bson.M{"token": token}
	_ = usersCollection.FindOne(ctx, filter).Decode(&user)
	if err := VerifyPassword(oldPassword, user.Password); err != nil {
		return fmt.Errorf("old password is incorrect")
	}
	hsPwd := HashPassword(newPassword)
	updateObj := bson.M{}
	updateObj["password"] = hsPwd
	updateObj["updated_at"] = carbon.Now().ToDateTimeString()
	update := bson.M{"$set": updateObj}
	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Println(result.ModifiedCount)
	return nil
}
