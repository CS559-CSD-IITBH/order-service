package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	OrderID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StoreID      primitive.ObjectID `bson:"storeID" json:"storeID"`
	UserID       uint               `bson:"userID" json:"userID"`
	Items        []OrderItem        `bson:"items" json:"items"`
	TotalAmount  float64            `bson:"totalAmount" json:"totalAmount"`
	Status       string             `bson:"status" json:"status"`
	DeliveryInfo DeliveryInfo       `bson:"deliveryInfo" json:"deliveryInfo"`
}

type OrderItem struct {
	ItemID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Price       float64            `bson:"price" json:"price"`
}

type DeliveryInfo struct {
	DeliveryAgentID string `bson:"deliveryAgentUID" json:"deliveryAgentUID"`
	CurrentLocation string `bson:"currentLocation" json:"currentLocation"`
}
