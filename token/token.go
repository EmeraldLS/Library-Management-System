package token

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/EmeraldLS/Library_Management_System/db"
	"github.com/EmeraldLS/Library_Management_System/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang-module/carbon"
	"go.mongodb.org/mongo-driver/bson"
)

var jwt_key = os.Getenv("jwt_key")

type JWTClaim struct {
	UserID   string
	Name     string
	Email    string
	UserType string
	*jwt.StandardClaims
}

func GenerateToken(userID string, name string, email string, userType string) (string, string, int64, error) {
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	claim := JWTClaim{
		UserID:   userID,
		Name:     name,
		Email:    email,
		UserType: userType,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	refreshClaim := JWTClaim{
		UserID:   userID,
		Name:     name,
		Email:    email,
		UserType: userType,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(jwt_key))
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim).SignedString([]byte(jwt_key))
	if err != nil {
		return "", "", 0, err
	}
	return token, refreshToken, expirationTime, nil
}

func ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwt_key), nil
	})

	_, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("invalid token")
		return err
	}

	var user model.User
	filter := bson.M{"token": signedToken}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = db.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if user.ExpirationTime < time.Now().Unix() {
		err = errors.New("token has expired")
		return err
	}

	if err != nil {
		err = errors.New("invalid token")
		return err
	}

	return nil
}

func UpdateAllToken(token string, refreshToken string, expirationTime int64, userID string) error {
	expirationTime = time.Now().Add(1 * time.Hour).Unix()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	filter := bson.M{"user_id": userID}
	updateObj := bson.M{}
	updateObj["token"] = token
	updateObj["refresh_token"] = refreshToken
	updateObj["updated_at"] = carbon.Now().ToDateTimeString()
	updateObj["expiration_time"] = expirationTime
	updateDetail := bson.M{"$set": updateObj}
	result, err := db.UsersCollection.UpdateOne(ctx, filter, updateDetail)
	if err != nil {
		return err
	}
	fmt.Println("Modified Count", result.ModifiedCount)
	return nil

}
