package model

import "time"

type Status int

const (
	// 접수 중
	Accepting = 1 + iota
	// 접수
	Receipt
	// 접수 취소
	Cancellati
	// 조리 중
	Cooking
	// 추가 주문
	AdditionalOrder
	// 배달 중
	InDelivery
	// 배달완료
	CompleteDelivery
)

type Order struct {
	Customer   Customer  `json:"customer" bson:"customer"`
	Menus      []Menu    `json:"menus" bson:"menus"`
	Status     Status    `json:"status" bson:"status"`
	Created_at time.Time `json:"createdAt" bson:"createdAt"`
}
