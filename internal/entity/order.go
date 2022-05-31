package entity

import (
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	processing OrderStatus = "processing"
	transit                = "transit"
	delivered              = "delivered"
	cancelled              = "cancelled"
)

type Order struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id" mapstructure:"_id"`
	Status       string              `bson:"status" json:"status" mapstructure:"status" validate:"required"`
	OrderNo      string              `bson:"orderNo" json:"orderNo" mapstructure:"orderNo"`
	Items        []Product           `bson:"items" json:"items" mapstructure:"items"`
	CustomerId   primitive.ObjectID  `bson:"customerId" json:"customerId" mapstructure:"customerId"`
	DeliveryDate primitive.Timestamp `bson:"deliveryDate" json:"deliveryDate" mapstructure:"deliveryDate" validate:"required"`
	CreatedAt    primitive.Timestamp `bson:"createdAt" json:"createdAt" validate:"required" mapstructure:"createdAt"`
}

type Orders []*Order

const OrderCollectionName = "orders"

func NewOrder() *Order {
	return &Order{}
}

func NewOrders() *Orders {
	return &Orders{}
}

func (o *Order) Persist(db *db.MongoStore) (*Order, error) {
	isNewOrder := o.CreatedAt.IsZero()
	var err error

	if isNewOrder {
		err = db.Insert(OrderCollectionName, o)
	} else {
		err = db.Replace(OrderCollectionName, utils.KeyValue{
			"_id": o.ID,
		}, o)
	}

	return o, err
}

func (o *Order) Get(db *db.MongoStore, filter utils.KeyValue) (*Order, error) {
	err := db.Get(OrderCollectionName, filter, o)
	return o, err
}

func (os Orders) GetAll(db *db.MongoStore, filter utils.KeyValue) (Orders, error) {
	err := db.GetAll(OrderCollectionName, filter, os)
	return os, err
}
