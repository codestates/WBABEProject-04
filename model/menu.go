package model

import (
	"WBABEProject-04/logger"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	MenuId    primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name" bson:"name"`           // 메뉴 이름
	Available bool               `json:"available" bson:"available"` // 주문 가능 여부
	Quantity  int                `json:"quantity" bson:"quantity"`   // 수량
	Grade     int                `json:"grade" bson:"grade"`
	Origin    string             `json:"origin" bson:"origin"`       // 원산지
	Price     uint               `json:"price" bson:"price"`         // 가격
	Spiciness Spiciness          `json:"spiciness" bson:"spiciness"` // 맵기 정도
	Favorites bool               `json:"favorites" bson:"favorites"` // 추천 여부
	Review    []Review           `json:"review,omitempty" bson:"review"`
}

func (m *Model) CreateMenu(menu Menu) error {
	logger.Debug("menu > CreateMenu")
	menu.MenuId = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.collectionMenu.InsertOne(ctx, menu); err != nil {
		log.Println("fail insert new menu")
		return fmt.Errorf("fail, insert")
	}
	return nil
}

func (m *Model) GetOneMenuWithName(menuName string) (Menu, error) {
	logger.Debug("menu > GetOneMenuWithName")
	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": menuName}

	var sMenu Menu

	if err := m.collectionMenu.FindOne(ctx, filter, opts...).Decode(&sMenu); err != nil {
		return sMenu, err
	} else {

		return sMenu, nil
	}
}

func (m *Model) GetOneMenu(pId primitive.ObjectID) (Menu, error) {
	logger.Debug("menu > GetOneMenu")
	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": pId}

	var sMenu Menu

	if err := m.collectionMenu.FindOne(ctx, filter, opts...).Decode(&sMenu); err != nil {
		return sMenu, err
	} else {

		return sMenu, nil
	}
}

func (m *Model) DeleteMenu(pId primitive.ObjectID) error {
	logger.Debug("menu > DeleteMenu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": pId}
	if res, err := m.collectionMenu.DeleteOne(ctx, filter); res.DeletedCount <= 0 {
		return fmt.Errorf("could not delete, not found menu %s", pId)
	} else if err != nil {
		return err
	}
	return nil
}

func (m *Model) UpdateMenu(menu Menu) error {
	logger.Debug("menu > UpdateMenu")
	filter := bson.M{"name": menu.Name}
	update := bson.M{
		"$set": bson.M{
			"origin":    menu.Origin,
			"quantity":  menu.Quantity,
			"price":     menu.Price,
			"spiciness": menu.Spiciness,
			"favorites": menu.Favorites,
			"available": menu.Available,
		},
	}
	if _, err := m.collectionMenu.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}

func (m *Model) UpdateOrderStatus(order Order, statusCode int) error {
	logger.Debug("menu > UpdateOrderStatus")
	filter := bson.M{"orderid": order.OrderId}
	update := bson.M{
		"$set": bson.M{
			"status": statusCode,
		},
	}
	if _, err := m.collectionMenu.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}

func (m *Model) GetMenu() ([]Menu, error) {
	logger.Debug("menu > GetMenu")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sMenu []Menu

	cursor, err := m.collectionMenu.Find(ctx, bson.M{})
	// 아래의 코드를 더 깔끔하게 바꿀 수 있을 것 같습니다.
	if err != nil {
		return sMenu, err
	} else {
		if err = cursor.All(context.TODO(), &sMenu); err != nil {
			return sMenu, err
		}
		fmt.Println(sMenu)
		return sMenu, nil
	}
}
