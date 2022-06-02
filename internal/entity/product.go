package entity

import (
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"_id" mapstructure:"_id"`
	Name        string              `bson:"name" json:"name" mapstructure:"name"`
	Description string              `bson:"description" json:"description" mapstructure:"description"`
	Price       float64             `bson:"price" json:"price" mapstructure:"price"`
	SKU         string              `bson:"sku" json:"sku" mapstructure:"sku"`
	CreatedAt   primitive.Timestamp `bson:"createdAt" json:"createdAt" mapstructure:"createdAt"`
}

const ProductCollectionName = "products"

func NewProduct() Product {
	return Product{}
}

func (p *Product) GetProduct(db *db.MongoStore, filter utils.KeyValue) (*Product, error) {
	err := db.Get(ProductCollectionName, filter, p)
	return p, err
}

func (p *Product) Persist(db *db.MongoStore) (*Product, error) {
	isNew := p.CreatedAt.IsZero()
	var err error
	if isNew {
		err = db.Insert(ProductCollectionName, p)
	} else {
		err = db.Replace(ProductCollectionName, utils.KeyValue{
			"_id": p.ID,
		}, p)
	}

	return p, err
}
