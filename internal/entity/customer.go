package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Lat  int32 `bson:"lat" json:"lat" validate:"required"`
	Long int32 `bson:"long" json:"long" validate:"required"`
}

type Customer struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" validate:"required"`
	Position Location           `bson:"position" json:"position" validate:"required"`
}

func NewCustomer() *Customer {
	return &Customer{}
}

const CustomerCollectionName = "customers"
