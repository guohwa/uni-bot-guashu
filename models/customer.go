package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	Name      string             `bson:"name"`
	Capital   float64            `bson:"capital"`
	ApiKey    string             `bson:"apiKey"`
	ApiSecret string             `bson:"apiSecret"`
	Status    string             `bson:"status"`
}
