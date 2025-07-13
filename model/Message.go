package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MessageSchema struct {
	Date   time.Time     `json:"date" bson:"date"`
	Body   string        `json:"body" bson:"body"`
	ID     bson.ObjectID `bson:"_id"`
	Sender bson.ObjectID `json:"userID" bson:"userID"`
	RoomID bson.ObjectID `json:"RoomID" bson:"RoomID"`
}

var Message *mongo.Collection

func (newMessage *MessageSchema) Save() error {
	newMessage.ID = bson.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := Message.InsertOne(ctx, newMessage); err != nil {
		return err
	}
	return nil
}
