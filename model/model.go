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
	mgUrl := "mongodb://127.0.0.1:27017"
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {

		return nil, err
	} else {
		db := r.client.Database("project")
		r.collectionCustomer = db.Collection("customer")
		r.collectionMenu = db.Collection("menu")
		r.collectionOrder = db.Collection("order")
		r.collectionReview = db.Collection("review")
	}

	return r, nil
}
