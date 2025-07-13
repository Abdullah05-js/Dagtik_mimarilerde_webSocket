package kafka

import (
	"ChatAPP/v2/model"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var reader *kafka.Reader

type ConnectionInfo struct {
	RoomID bson.ObjectID   `json:"RoomID"`
	Conn   *websocket.Conn `json:"Conn"`
}

var Mux sync.Mutex
var Clients map[bson.ObjectID]ConnectionInfo

func InitReader(topic string, groupID string, brokers []string) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Topic:   topic,
		GroupID: groupID,
		Brokers: brokers,
	})
	Clients = make(map[bson.ObjectID]ConnectionInfo)
}

func CloseReader() {
	if reader != nil {
		if err := reader.Close(); err != nil {
			log.Println("Failed to close Kafka reader:", err)
		}
	}
}

func ListenMessages() {
	for {
		bytesMessage, err := reader.ReadMessage(context.Background()) // we used context.Background() because we don't want a cancellation happen read forever
		if err != nil {
			log.Fatal("failed to get message", err.Error())
		}
		var message model.MessageSchema
		if err := json.Unmarshal(bytesMessage.Value, &message); err != nil {
			log.Fatal("failed to pare  to json", err.Error())
		}
		var toDelete []bson.ObjectID

		//send to clients
		Mux.Lock()
		for userID, client := range Clients {

			if client.RoomID == message.RoomID {
				if err := client.Conn.WriteMessage(websocket.TextMessage, bytesMessage.Value); err != nil {
					log.Printf("Failed to write to client %s: %v\n", userID.String(), err)
					client.Conn.Close()
					toDelete = append(toDelete, userID)
					continue
				}
			}

		}
		Mux.Unlock()

		Mux.Lock()
		for _, id := range toDelete {
			delete(Clients, id)
		}
		Mux.Unlock()

		log.Printf("[%s] %s: %s\n", message.RoomID, message.Sender, message.Body)
	}
}
