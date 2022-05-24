package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderStatus string

const (
	processing OrderStatus = "processing"
	transit                = "transit"
	delivered              = "delivered"
	cancelled              = "cancelled"
)

type Order struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Status       OrderStatus         `bson:"status" json:"status" validate:"required"`
	DeliveryDate primitive.Timestamp `bson:"deliveryDate" json:"deliveryDate" validate:"required"`
	CreatedAt    primitive.Timestamp `bson:"createdAt" json:"createdAt" validate:"required"`
}

const OrderCollectionName = "orders"

func NewOrder() *Order {
	return &Order{}
}
