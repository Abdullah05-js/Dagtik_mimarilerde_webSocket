package kafka

import (
	"context"
	"fmt"
	"time"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var writer *kafka.Writer

func InitWriter(Topic string, Brokers []string) {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  Brokers,
		Topic:    Topic,
		Balancer: &kafka.Hash{},
	})
}

func CloseWriter() {
	if writer != nil {
		if err := writer.Close(); err != nil {
			fmt.Println("Failed to close Kafka writer:", err)
		}
	}
}

func SendMessage(roomID bson.ObjectID, body []byte) error {
	message := kafka.Message{
		Key:   []byte(roomID.String()),
		Value: body,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // in 5 seconds if the message not sended wireMessages will return err
	defer cancel()
	if err := writer.WriteMessages(ctx, message); err != nil {
		return err
	}

	return nil
}
