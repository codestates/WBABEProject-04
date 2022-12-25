package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	CustomerId primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Nicname    string             `json:"nicname" bson:"nicname"`
	Phone      string             `json:"phone" bson:"phone"`
	Address    string             `json:"address" bson:"address"`
	Orders     []Order            `json:"orders" bson:"orders"`
}
