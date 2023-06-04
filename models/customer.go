package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	Name      string             `bson:"name"`
	Token     string             `bson:"token"`
	ApiKey    string             `bson:"apiKey"`
	ApiSecret string             `bson:"apiSecret"`
	Base      float64            `bson:"base"`
	Capital   float64            `bson:"capital"`
	Ratio     float64            `bson:"ratio"`
	Level1    float64            `bson:"level1"`
	Level2    float64            `bson:"level2"`
	Mode      string             `bson:"mode"`
	Status    string             `bson:"status"`
}
