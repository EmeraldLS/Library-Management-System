package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BooksCollection *mongo.Collection
var UsersCollection *mongo.Collection

func init() {
	_ = godotenv.Load(".env")
	uri := os.Getenv("uri")
	dbname := os.Getenv("dbname")
	booksCol := os.Getenv("booksCol")
	usersCol := os.Getenv("usersCol")
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("An error occured while connecting to mongodb. ", err)
	}
	BooksCollection = client.Database(dbname).Collection(booksCol)
	UsersCollection = client.Database(dbname).Collection(usersCol)
	fmt.Println("All collections are ready for operations.")
	defer cancel()

}
