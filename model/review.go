package model

import (
	"WBABEProject-04/logger"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Content    string             `json:"content" bson:"content"`
	MenuId     primitive.ObjectID `json:"menuid" bson:"menuid"`
	CustomerID primitive.ObjectID `json:"customerid,omitempty" bson:"customerid"`
	Grade      int                `json:"grade" bson:"grade"`
	IsWrite    bool               `json:"iswrite,omitempty" bson:"iswrite"`
	CreatedAt  time.Time          `json:"createdat,omitempty" bson:"createdat"`
}

func (m *Model) CreateReview(review Review) error {
	logger.Debug("review > CreateReview")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.collectionReview.InsertOne(ctx, &review); err != nil {
		log.Println("fail insert new review")
		return fmt.Errorf("fail, insert review")
	}
	return nil
}

func (m *Model) GetReviewByMenuID(menuID primitive.ObjectID) ([]Review, error) {
	logger.Debug("review > GetReviewByMenuID")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"menuid": menuID}
	var reviews []Review
	if cursor, err := m.collectionReview.Find(ctx, filter); err != nil {
		return reviews, err
	} else {
		if err = cursor.All(context.TODO(), &reviews); err != nil {
			log.Println("fail find review")
			return reviews, fmt.Errorf("fail, find review")
		}
		return reviews, nil
	}
}
