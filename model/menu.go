package model

import (
	"WBABEProject-04/logger"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Spiciness int

const (
	Normal = iota + 1
	Spicy
	Hot
	Hell
	Burning
)

type Menu struct {
	Name      string    `json:"name" bson:"name"`           // 메뉴 이름
	Available bool      `json:"available" bson:"available"` // 주문 가능 여부
	Quantity  int       `json:"quantity" bson:"quantity"`   // 수량
	Grade     int       `json:"grade" bson:"grade"`
	Origin    string    `json:"origin" bson:"origin"`       // 원산지
	Price     uint      `json:"price" bson:"price"`         // 가격
	Spiciness Spiciness `json:"spiciness" bson:"spiciness"` // 맵기 정도
	Favorites bool      `json:"favorites" bson:"favorites"` // 추천 여부
	Review    []Review  `json:"review,omitempty" bson:"review"`
}

func (m *Model) CreateMenu(menu Menu) error {
	logger.Debug("seller > CreateMenu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.collectionMenu.InsertOne(ctx, menu); err != nil {
		log.Println("fail insert new menu")
		return fmt.Errorf("fail, insert")
	}
	return nil
}

func (m *Model) GetOneMenu(name string) (Menu, error) {

	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	fmt.Println(name)

	var sMenu Menu

	if err := m.collectionMenu.FindOne(ctx, filter, opts...).Decode(&sMenu); err != nil {
		fmt.Println(sMenu, err)
		return sMenu, err
	} else {
		fmt.Println(sMenu, err)
		return sMenu, nil
	}

}
