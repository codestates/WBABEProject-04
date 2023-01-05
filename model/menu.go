package model

import (
	"WBABEProject-04/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryData struct {
	Condition string `form:"con" binding:"required"`
	OrderBy   int    `form:"ord" binding:"required"`
}

type Spiciness int

const (
	Normal = iota + 1
	Spicy
	Hot
	Hell
	Burning
)

type Menu struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name" bson:"name"`                     // 메뉴 이름
	SoldOut   bool               `json:"soldout,omitempty" bson:"soldout"`     // 주문 가능 여부
	Quantity  int                `json:"quantity,omitempty" bson:"quantity"`   // 수량
	Grade     int                `json:"grade" bson:"grade"`                   // 평점
	Origin    string             `json:"origin" bson:"origin"`                 // 원산지
	Price     uint               `json:"price" bson:"price"`                   // 가격
	Spiciness Spiciness          `json:"spiciness,omitempty" bson:"spiciness"` // 맵기 정도
	Favorites bool               `json:"favorites" bson:"favorites"`           // 추천 여부
	Review    []Review           `json:"review,omitempty" bson:"review"`       // 리뷰
	Count     int                `json:"count" bson:"count"`                   // 재주문수
	CreatedAt time.Time          `json:"createdat" bson:"createdat"`
}

func (m *Model) CreateMenu(menu Menu) error {
	logger.Debug("menu > CreateMenu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	menu.ID = primitive.NewObjectID()
	menu.CreatedAt = time.Now()

	if _, err := m.collectionMenu.InsertOne(ctx, &menu); err != nil {
		log.Println("fail insert new menu")
		return fmt.Errorf("fail, insert")
	}
	return nil
}

func (m *Model) GetOneMenu(flag, elem string) (Menu, error) {
	logger.Debug("menu > GetOneMenu")

	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M
	if flag == "name" {
		filter = bson.M{"name": elem}
	}

	var sMenu Menu

	if err := m.collectionMenu.FindOne(ctx, filter, opts...).Decode(&sMenu); err != nil {
		return sMenu, err
	} else {
		return sMenu, nil
	}
}

func (m *Model) DeleteMenu(name string) error {
	logger.Debug("menu > DeleteMenu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	if res, err := m.collectionMenu.DeleteOne(ctx, filter); res.DeletedCount <= 0 {
		return fmt.Errorf("could not delete, not found menu %s", name)
	} else if err != nil {
		return err
	}
	return nil
}
func (m *Model) UpdateMenuGrade(grade int, menuName string) error {
	// logger.Debug("menu > UpdateMenuGrade")
	// filter := bson.M{"name": menuName}
	return nil
}
func (m *Model) UpdateMenu(menu Menu, menuName string) error {

	logger.Debug("menu > UpdateMenu")
	filter := bson.M{"name": menuName}
	update := bson.M{
		"$set": bson.M{
			"origin":    menu.Origin,
			"price":     menu.Price,
			"favorites": menu.Favorites,
			"soldout":   menu.SoldOut,
		},
	}

	if _, err := m.collectionMenu.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	} else {
		return nil
	}
}

func (m *Model) GetMenu(flag string) ([]Menu, error) {
	logger.Debug("menu > GetMenu")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}

	var sMenu []Menu

	if cursor, err := m.collectionMenu.Find(ctx, filter); err != nil {
		return sMenu, err
	} else {
		if err = cursor.All(context.TODO(), &sMenu); err != nil {
			return sMenu, err
		}
		return sMenu, nil
	}
}

func (m *Model) GetSortedMenu(query QueryData) ([]Menu, error) {
	filter := bson.D{}
	condition := bson.D{{query.Condition, query.OrderBy}}
	opts := options.Find().SetSort(condition)
	if cursor, err := m.collectionMenu.Find(context.TODO(), filter, opts); err != nil {
		log.Printf("fail, find %v", query)
		return nil, fmt.Errorf("find condition : %s", query.Condition)
	} else {
		var result []Menu
		if err = cursor.All(context.TODO(), &result); err != nil {
			log.Printf("fail, find %v", query)
			return nil, fmt.Errorf("search criteria not found.. %s", query.Condition)
		}
		for _, result := range result {
			res, _ := json.Marshal(result)
			fmt.Println(string(res))
		}
		return result, nil
	}
}
func (m *Model) IncreaseMenuVolume(menu Menu) error {
	opts := []*options.FindOneOptions{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": menu.Name}
	var sMenu Menu

	if err := m.collectionMenu.FindOne(ctx, filter, opts...).Decode(&sMenu); err != nil {
		return err
	} else {
		update := bson.M{
			"$set": bson.M{
				"count": sMenu.Count + 1,
			},
		}
		if _, err := m.collectionMenu.UpdateOne(context.Background(), filter, update); err != nil {
			return err
		} else {
			return nil
		}
	}
}
