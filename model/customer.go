package model

type Customer struct {
	Nicname string  `json:"nicname" bson:"nicname"`
	Phone   string  `json:"phone" bson:"phone"`
	Address string  `json:"address" bson:"address"`
	Orders  []Order `json:"orders" bson:"orders"`
}
