package entity

import (
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Lat  float32 `bson:"lat" json:"lat" validate:"required"`
	Long float32 `bson:"long" json:"long" validate:"required"`
}

type Customer struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name      string              `bson:"name" json:"name" validate:"required"`
	Position  Location            `bson:"position" json:"position" validate:"required"`
	CreatedAt primitive.Timestamp `bson:"createdAt" json:"createdAt"`
}

func NewCustomer() *Customer {
	return &Customer{}
}

func (c *Customer) GetCustomer(db *db.MongoStore, filter utils.KeyValue) (*Customer, error) {
	err := db.Get(CustomerCollectionName, filter, c)
	return c, err
}

func (c *Customer) Persist(db *db.MongoStore) (*Customer, error) {
	isNewCustomer := c.CreatedAt.IsZero()
	var err error
	if isNewCustomer {
		err = db.Insert(CustomerCollectionName, c)
	} else {
		err = db.Replace(CustomerCollectionName, utils.KeyValue{"_id": c.ID}, c)
	}

	return c, err
}

const CustomerCollectionName = "customers"
