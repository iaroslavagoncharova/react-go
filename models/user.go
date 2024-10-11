package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username" validate:"required,min=3,max=50"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password" bson:"password" validate:"required,min=6"`
	Role     string             `json:"role" bson:"role"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserWithoutPassword struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Role     string             `json:"role" bson:"role"`
}

type UpdateUser struct {
	Username *string `json:"username" bson:"username" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email" bson:"email" validate:"omitempty,email"`
	Password *string `json:"password" bson:"password" validate:"omitempty,min=6"`
}
