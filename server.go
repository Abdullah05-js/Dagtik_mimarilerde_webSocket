package main

import (
	"ChatAPP/v2/database"
	"ChatAPP/v2/kafka"
	"ChatAPP/v2/model"
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("ERROR FROM server.go/15 %v:", err)
	}

	if err := database.InitMongo(); err != nil {
		log.Fatal(err.Error())
	}

	Broker := os.Getenv("BROKER")
	Topic := os.Getenv("TOPIC")
	Group := os.Getenv("GROUP_ID")

	kafka.InitWriter(Topic, []string{Broker})
	kafka.InitReader(Topic, Group, []string{Broker}) // burda group önemli diğer kopya backende aynı groupu verirsek mesajları paylaşırlar

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, //max 10mb
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Frontend
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			response, isOK := err.(*fiber.Error)

			if !isOK {
				return c.Status(500).JSON(fiber.Map{
					"Message": "Internal Server Error",
				})
			}

			return c.Status(response.Code).JSON(fiber.Map{
				"Message": response.Message,
			})
		}
		return nil
	})

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		var userID bson.ObjectID // dışarı tanımla
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			var message model.MessageSchema
			if err := json.Unmarshal(msg, &message); err != nil {
				c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid message format"}`))
				continue
			}
			if userID.IsZero() {
				userID = message.Sender
				kafka.Mux.Lock()
				kafka.Clients[userID] = kafka.ConnectionInfo{
					RoomID: message.RoomID,
					Conn:   c,
				}
				kafka.Mux.Unlock()
			}

			go message.Save()
			_ = kafka.SendMessage(message.RoomID, msg)
		}

		kafka.Mux.Lock()
		delete(kafka.Clients, userID)
		kafka.Mux.Unlock()
	}))

	go kafka.ListenMessages()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	app.Listen(":" + port)
}
