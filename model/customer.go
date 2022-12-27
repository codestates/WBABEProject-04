package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	// Acronym은 Uppercase를 사용합니다
	CustomerId primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	// nick name field에 오탈자가 있습니다
	Nicname    string             `json:"nicname" bson:"nicname"`
	Phone      string             `json:"phone" bson:"phone"`
	Address    string             `json:"address" bson:"address"`
	Orders     []Order            `json:"orders" bson:"orders"`
}
