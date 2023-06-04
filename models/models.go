package models

import (
	"context"
	"time"

	"bot/config"
	"bot/log"
	"bot/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var (
	UserCollection     *mongo.Collection
	CommandCollection  *mongo.Collection
	CustomerCollection *mongo.Collection
	SessionCollection  *mongo.Collection
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Config.Database.Host))
	if err != nil {
		log.Fatal(err)
	}

	UserCollection = Client.Database(config.Config.Database.Name).Collection("users")
	CommandCollection = Client.Database(config.Config.Database.Name).Collection("commands")
	CustomerCollection = Client.Database(config.Config.Database.Name).Collection("customers")
	SessionCollection = Client.Database(config.Config.Database.Name).Collection("sessions")

	res := UserCollection.FindOne(context.Background(), bson.M{"username": "admin"})
	if err = res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			user := User{
				ID:       primitive.NewObjectID(),
				Username: "admin",
				Password: utils.Encrypt("admin"),
				Role:     "Admin",
				Status:   "Enable",
			}
			result, err := UserCollection.InsertOne(context.Background(), &user)
			if err != nil {
				log.Fatal(err)
			}
			if id, ok := result.InsertedID.(primitive.ObjectID); ok {
				log.Infof("system init, admin id: %s", id.Hex())
			} else {
				log.Fatal("system init failed")
			}
		} else {
			log.Fatal(err)
		}
	}
}
