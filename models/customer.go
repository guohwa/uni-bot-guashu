package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	Name      string             `bson:"name"`
	Token     string             `bson:"token"`
	Exchange  string             `bson:"exchange"`
	ApiKey    string             `bson:"apiKey"`
	ApiSecret string             `bson:"apiSecret"`
	Capital   float64            `bson:"capital"`
	Scale     float64            `bson:"scale"`
	Level1    float64            `bson:"level1"`
	Level2    float64            `bson:"level2"`
	Status    string             `bson:"status"`
}
