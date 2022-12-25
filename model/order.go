package model

import (
	"WBABEProject-04/logger"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int

const (
	// 접수 중
	Accepting = 1 + iota
	// 접수 취소
	Cancellation
	// 접수
	Receipt

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
	OrderId    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Menus      []Menu             `json:"menus" bson:"menus"`
	Status     Status             `json:"status" bson:"status"`
	Created_at time.Time          `json:"createdAt" bson:"createdAt"`
}

func (m *Model) GetOrders() ([]Order, error) {
	logger.Debug("order > GetOrders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orders []Order

	cursor, err := m.collectionMenu.Find(ctx, bson.M{"status": bson.M{"$gte": 3}})
	if err != nil {
		return orders, err
	} else {
		if err = cursor.All(context.TODO(), &orders); err != nil {
			return orders, err
		}
		fmt.Println(orders)
		return orders, nil
	}
}

func (m *Model) UpdateStatus() ([]Order, error) {
	logger.Debug("order > GetStatus")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orders []Order

	cursor, err := m.collectionMenu.Find(ctx, bson.M{"status": bson.M{"$gte": 3}})
	if err != nil {
		return orders, err
	} else {
		if err = cursor.All(context.TODO(), &orders); err != nil {
			return orders, err
		}
		fmt.Println(orders)
		return orders, nil
	}
}
