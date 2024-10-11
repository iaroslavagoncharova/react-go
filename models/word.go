package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Word struct {
    ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    CollectionID primitive.ObjectID `json:"collectionId" bson:"collectionId"`
    Word         string             `json:"word" bson:"word" validate:"required,min=1,max=50"`
    Translation  string             `json:"translation" bson:"translation" validate:"required,min=1,max=100"`
    Difficulty   int                `json:"difficulty" bson:"difficulty" validate:"required,gte=1,lte=5"`
}

type UpdateWord struct {
    CollectionID *primitive.ObjectID `json:"collectionId" bson:"collectionId" validate:"omitempty"`
    Word         *string             `json:"word" bson:"word" validate:"omitempty,min=1,max=50"`
    Translation  *string             `json:"translation" bson:"translation" validate:"omitempty,min=1,max=100"`
    Difficulty   *int                `json:"difficulty" bson:"difficulty" validate:"omitempty,gte=1,lte=5"`
}