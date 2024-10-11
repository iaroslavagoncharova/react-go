package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Collection struct {
    ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Name        *string             `json:"name" bson:"name" validate:"required,min=3,max=50"`
    Description *string             `json:"description" bson:"description" validate:"required,min=3,max=200"`
    UserID      string             `json:"userId" bson:"userId"`
}

type UpdateCollection struct {
    Name        *string `json:"name" bson:"name" validate:"omitempty,min=3,max=50"`
    Description *string `json:"description" bson:"description" validate:"omitempty,min=3,max=200"`
}