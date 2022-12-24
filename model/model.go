package model

import (
	"WBABEProject-04/logger"
	"context"
	"fmt"
	"log"
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
	Price    uint   `json:"price" bson:"price"`       // 가격
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

// 메뉴이름을 받아 메뉴를 가져온다.
func (m *Model) GetOneMenu(flag, elem string) (Menu, error) {
	logger.Debug("Model > GetOneMenu")
	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M

	if flag == "name" {
		filter = bson.M{"name": elem}
	}
	var menus Menu
	if err := m.collection.FindOne(ctx, filter, opts...).Decode(&menus); err != nil {
		return menus, err
	} else {
		return menus, nil
	}
}

func (m *Model) GetMenuList() []Menu {
	logger.Debug("Model > GetMenuList")
	fmt.Println("GetMenuList")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	var menus []Menu
	if err = cursor.All(ctx, &menus); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	return menus
}

// 메뉴를 생성한다.
func (m *Model) CreateMenu(menus Menu) error {
	logger.Debug("Model > CreateMenu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.collection.InsertOne(ctx, menus); err != nil {
		log.Println("fail insert new menu")
		return fmt.Errorf("fail, insert")
	}
	return nil
}

func (m *Model) DeleteMenu(smenu string) error {
	logger.Debug("Model > DeleteMenu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": smenu}
	if res, err := m.collection.DeleteOne(ctx, filter); res.DeletedCount <= 0 {
		return fmt.Errorf("could not delete, not found menu %s", smenu)
	} else if err != nil {
		return err
	}
	return nil
}

func (m *Model) UpdateMenu(menu Menu) error {
	fmt.Println("UpdateMenu : ", menu)
	filter := bson.M{"name": menu.Name}
	update := bson.M{
		"$set": bson.M{
			"order":    menu.Order,
			"quantity": menu.Quantity,
			"price":    menu.Price,
			"spicy":    menu.Spicy,
			"origin":   menu.Origin,
		},
	}
	if _, err := m.collection.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}
	fmt.Println(">>??")
	return nil

}
