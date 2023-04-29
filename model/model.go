package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BookCode      int                `json:"book_code,omitempty" bson:"book_code,omitempty"`
	BookID        string             `json:"book_id,omitempty" bson:"book_id,omitempty"`
	BookName      string             `json:"book_name,omitempty" bson:"book_name,omitempty" validate:"required"`
	ISBN          int                `json:"isbn,omitempty" bson:"isbn,omitempty" validate:"required"`
	Author        string             `json:"author,omitempty" bson:"author,omitempty" validate:"required"`
	Edition       int                `json:"edition,omitempty" bson:"edition,omitempty" validate:"required"`
	Publication   string             `json:"publication,omitempty" bson:"publication,omitempty" validate:"required"`
	Registered_At string             `json:"registered_at,omitempty" bson:"registered_at,omitempty"`
	Updated_At    string             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserCode       int                `json:"user_code,omitempty" bson:"user_code,omitempty"`
	UserID         string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name           string             `json:"full_name,omitempty"  bson:"full_name,omitempty" validate:"required"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Password       string             `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	UserType       string             `json:"user_type,omitempty" bson:"user_type,omitempty" validate:"required,eq=ADMIN|eq=USER"`
	Token          string             `json:"token,omitempty" bson:"token,omitempty"`
	Refresh_Token  string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	Registered_At  string             `json:"registered_at,omitempty" bson:"registered_at,omitempty"`
	Updated_At     string             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	ExpirationTime int64              `json:"expiration_time,omitempty" bson:"expiration_time,omitempty"`
}

type LoginStruct struct {
	UserID   string `json:"user_id,omitempty" bson:"user_id,omitempty" validate:"required"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	UserType string `json:"user_type,omitempty" bson:"user_type,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password,omitempty" validate:"required"`
	NewPassword string `json:"new_password,omitempty" validate:"required"`
}
