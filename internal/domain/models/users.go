package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password,omitempty" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

//signup and signin

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignUpResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"string"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

//public user(no signup)

type PublicUserResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
}

func (u *User) Public() *PublicUserResponse {
	return &PublicUserResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}
