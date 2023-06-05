package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Command struct {
	ID         primitive.ObjectID `bson:"_id"`
	CustomerID primitive.ObjectID `bson:"customerId"`
	Exchange   string             `bson:"exchange"`
	Action     string             `bson:"action"`
	Symbol     string             `bson:"symbol"`
	Side       string             `bson:"side"`
	Capital    float64            `bson:"capital"`
	Size       float64            `bson:"size"`
	Quantity   float64            `bson:"quantity"`
	Comment    string             `bson:"comment"`
	Status     string             `bson:"status"`
	Reason     string             `bson:"reason"`
	Time       time.Time          `bson:"time"`
}
