package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type Menu struct {
	Name     string `json:"name" bson:"name"`         // 메뉴 이름
	Order    bool   `json:"order" bson:"order"`       // 주문 가능 여부
	Quantity int    `json:"quantity" bson:"quantity"` // 수량
	Origin   string `json:"origin" bson:"origin"`     // 원산지
	Price    int    `json:"price" bson:"price"`       // 가격
	Spicy    string `json:"spicy" bson:"spicy"`       // 맵기 정도
}

func NewModel() (*Model, error) {
	r := &Model{}
	var err error
	mgUrl := "mongodb://127.0.0.1:27017"
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("go-ready")
		r.collection = db.Collection("tPerson")
	}
	return r, nil
}

func (m *Model) GetOneMenu(flag, elem string) (Menu, error) {
	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M

	// 체크해줄 필요가 있는가? -> 고민해보자
	if flag == "name" {
		filter = bson.M{"name": elem}
		fmt.Println("플래그는 네임이 맞음")
	}
	var menus Menu
	if err := m.collection.FindOne(ctx, filter, opts...).Decode(&menus); err != nil {
		return menus, err
	} else {
		return menus, nil
	}
}

func (m *Model) CreateMenu(menus Menu) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.collection.InsertOne(ctx, menus); err != nil {
		fmt.Println("fail insert new menu")
		return fmt.Errorf("fail, insert")
	}
	return nil
}
