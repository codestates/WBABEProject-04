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

type Status int

const (
	Accepting = 1 + iota
	Cancellation
	Receipt
	Cooking
	AdditionalOrder
	InDelivery
	CompleteDelivery
)

type OrderNumber struct {
	OrderList []primitive.ObjectID `json:"ordernumber" bson:"ordernumber"`
}
type Customer struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Nicname string             `json:"nicname" bson:"nicname"`
	Phone   string             `json:"phone" bson:"phone"`
	Address string             `json:"address" bson:"address"`
	Orders  []Order            `json:"orders" bson:"orders"`
}

type Order struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	CustomerID primitive.ObjectID `bson:"customerid" json:"customerid,omitempty"`
	Menus      []Menu             `json:"menus" bson:"menus"`
	Status     Status             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}

func (m *Model) GetOrderStatusByOrderID(orderID primitive.ObjectID) (int, error) {
	logger.Debug("order > GetOrderStatusByOrderID")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": orderID}
	order := Order{}
	if err := m.collectionOrder.FindOne(ctx, filter).Decode(&order); err != nil {
		return int(order.Status), err
	} else {
		return int(order.Status), nil
	}
}

func (m *Model) GetOrdersByUserID(userID primitive.ObjectID) ([]Order, error) {
	logger.Debug("order > GetOrdersByUserID")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var orders []Order
	cursor, err := m.collectionOrder.Find(ctx, bson.M{"customeriD": userID})
	if err != nil {
		return orders, err
	} else {
		if err = cursor.All(context.TODO(), &orders); err != nil {
			return orders, err
		}
		return orders, nil
	}
}

func (m *Model) GetOrdersInfoByUserID(flag string, userID primitive.ObjectID) ([]Order, error) {
	logger.Debug("order > GetOrdersInfoByUserID")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var opts *options.FindOptions
	var filter bson.M
	if flag == "his" {
		condition := bson.M{"createdAt": -1}
		opts = options.Find().SetSort(condition)
		filter = bson.M{"customerid": userID, "status": bson.M{"$eq": 7}}
	} else {
		filter = bson.M{"customerid": userID, "status": bson.M{"$lt": 7}}
	}
	var orders []Order
	if cursor, err := m.collectionOrder.Find(ctx, filter, opts); err != nil {
		return orders, err
	} else {
		if err = cursor.All(context.TODO(), &orders); err != nil {
			return orders, err
		}
		fmt.Println(orders)
		return orders, nil
	}
}

func (m *Model) GetOrderByID(orderID primitive.ObjectID) (Order, error) {
	logger.Debug("order > GetOrderByID")
	opts := []*options.FindOneOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": orderID}
	order := Order{}
	if err := m.collectionOrder.FindOne(ctx, filter, opts...).Decode(&order); err != nil {
		return order, err
	} else {
		return order, nil
	}
}

func (m *Model) GetOrders() ([]Order, error) {
	logger.Debug("order > GetOrders")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var orders []Order
	if cursor, err := m.collectionOrder.Find(ctx, bson.M{"status": bson.M{"$eq": 3}}); err != nil {
		return orders, fmt.Errorf("fail, find order")
	} else {
		if err = cursor.All(context.TODO(), &orders); err != nil {
			return orders, fmt.Errorf("fail, find order")
		}
		return orders, nil
	}
}

func (m *Model) CreateOrder(order Order) error {
	logger.Debug("order > CreateOrder")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := m.collectionOrder.InsertOne(ctx, order); err != nil {
		log.Println("fail, insert order")
		return fmt.Errorf("fail, insert order")
	} else {
		return nil
	}
}

func (m *Model) CreateCustomer(customer Customer) error {
	logger.Debug("order > CreateCustomer")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := m.collectionCustomer.InsertOne(ctx, customer); err != nil {
		log.Println("fail, insert customer")
		return fmt.Errorf("fail, insert customer")
	}
	return nil
}

func (m *Model) UpdateOrder(order Order, orderID primitive.ObjectID) error {
	logger.Debug("order > UpdateOrder")
	filter := bson.M{"_id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"menus": order.Menus,
		},
	}
	if _, err := m.collectionOrder.UpdateOne(context.Background(), filter, update); err != nil {
		log.Println("fail, update order statuscode")
		return fmt.Errorf("fail, update order statuscode")
	} else {
		return nil
	}
}

func (m *Model) UpdateOrderStatus(order Order, statusCode int) error {
	logger.Debug("order > UpdateOrderStatus")
	filter := bson.M{"_id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"status": statusCode,
		},
	}
	if _, err := m.collectionOrder.UpdateOne(context.Background(), filter, update); err != nil {
		log.Println("fail, update order statuscode")
		return fmt.Errorf("fail, update order statuscode")
	}
	return nil
}
