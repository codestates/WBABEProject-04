package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client             *mongo.Client
	collectionMenu     *mongo.Collection
	collectionOrder    *mongo.Collection
	collectionCustomer *mongo.Collection
	collectionReview   *mongo.Collection
}

func NewModel() (*Model, error) {
	r := &Model{}
	var err error
	//URL과 같은 Acronym은 대문자를 사용합니다.
	mgUrl := "mongodb://127.0.0.1:27017"
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("project")
		r.collectionCustomer = db.Collection("custoemr")
		r.collectionMenu = db.Collection("menu")
		r.collectionOrder = db.Collection("order")
		r.collectionReview = db.Collection("review")
	}

	return r, nil
}

// func (m *Model) GetOneStore(name string) (Store, error) {

// 	logger.Debug("model > GetOneStore")
// 	opts := []*options.FindOneOptions{}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	filter := bson.M{"name": "과일가게"}
// 	fmt.Println(name)

// 	var store Store

// 	if err := m.collectionStore.FindOne(ctx, filter, opts...).Decode(&store); err != nil {
// 		fmt.Println(store, err)
// 		return store, errㅏ
// 	} else {
// 		fmt.Println(store, err)
// 		return store, nil
// 	}

// }
