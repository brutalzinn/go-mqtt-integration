package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/brutalzinn/go-mqtt-integration/command"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)

var (
	lastEvent       command.Event
	isUnsupportedOS = false
	wsClients       = make(map[*websocket.Conn]bool)
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var event *command.Event
	logrus.Info("command received")
	err := json.Unmarshal(msg.Payload(), &event)
	if err != nil {
		logrus.Error("Error unmarshaling MQTT message: %s", err)
		return
	}
	lastEvent = *event
	err = command.Run(event, msg.Payload())
	if err != nil {
		logrus.Error("Error unmarshaling MQTT message: %s", err)
	}
}

func connectMQTT() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(os.Getenv("MQTT_BROKER"))
	opts.SetClientID(os.Getenv("MQTT_ID_CLIENT"))
	opts.SetUsername(os.Getenv("MQTT_USER"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetDefaultPublishHandler(messagePubHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Fatal("Error connecting to MQTT broker: %v", token.Error())
	}
	return client
}

func main() {
	if os.Getenv("OS_TYPE") == "unsupported" {
		isUnsupportedOS = true
	}

	if _, err := os.Stat(os.Getenv("ASSETS_FOLDER")); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("ASSETS_FOLDER"), os.ModePerm)
	}
	mqttClient := connectMQTT()

	topic := "hotest"
	if token := mqttClient.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		logrus.Fatalf("Error subscribing to topic: %v", token.Error())
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"LastEvent": lastEvent,
		})
	})

	app.Get("/content-url", func(c *fiber.Ctx) error {
		files, err := filepath.Glob(filepath.Join(os.Getenv("ASSETS_FOLDER"), "*"))
		if err != nil {
			return err
		}

		urls := make([]string, len(files))
		for i, file := range files {
			urls[i] = filepath.Join(os.Getenv("ASSETS_FOLDER"), filepath.Base(file))
		}
		return c.JSON(urls)
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		wsClients[c] = true
		defer func() {
			delete(wsClients, c)
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("WebSocket error:", err)
				return
			}
		}
	}))

	///we need it?
	app.Static("/assets", "./assets")

	logrus.Fatal(app.Listen(":3000"))
}
