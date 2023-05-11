package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id, omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Username string             `json:"username,omitempty" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
	Phone    string             `json:"phone,omitempty" validate:"required"`
	Cnpj     string             `json:"cnpj,omitempty" validate:"required"`
	Firmware string             `json:"firmware,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
