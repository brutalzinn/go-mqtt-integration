package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/brutalzinn/go-mqtt-integration/command"
	"github.com/brutalzinn/go-mqtt-integration/wshelper"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)

var (
	lastEvent command.Event
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var event *command.Event
	err := json.Unmarshal(msg.Payload(), &event)
	if err != nil {
		logrus.Error("Error unmarshaling MQTT message: %s", err)
		return
	}
	logrus.Info("command received")
	lastEvent = *event
	err = command.Run(event, msg.Payload())
	if err != nil {
		logrus.Error("Error unmarshaling MQTT message: %s", err)
	}
}

func startMQTTClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(os.Getenv("MQTT_BROKER"))
	opts.SetClientID(os.Getenv("MQTT_ID_CLIENT"))
	opts.SetUsername(os.Getenv("MQTT_USER"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetDefaultPublishHandler(messagePubHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Fatal("Error connecting to MQTT broker: %s", token.Error())
	}
	if token := client.Subscribe(os.Getenv("MQTT_TOPIC"), 1, nil); token.Wait() && token.Error() != nil {
		logrus.Fatalf("Error subscribing to topic: %v", token.Error())
	}

	logrus.Info("MQTT CONNECTED")
	return client
}

func main() {
	assetsFolder := os.Getenv("ASSETS_FOLDER")
	logrus.Info("Starting GO MQTT INTEGRATION")
	if _, err := os.Stat(assetsFolder); os.IsNotExist(err) {
		os.Mkdir(assetsFolder, os.ModePerm)
		logrus.Info("Dir for assets is created at %s", assetsFolder)
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
		files, err := filepath.Glob(filepath.Join(assetsFolder, "*"))
		if err != nil {
			logrus.Error("WebSocket error:", err)
			return err
		}
		urls := make([]string, len(files))
		for i, file := range files {
			urls[i] = filepath.Join(assetsFolder, filepath.Base(file))
		}
		return c.JSON(urls)
	})

	app.Get("/ws", websocket.New(wshelper.WebSocketHandler))
	app.Static("/assets", filepath.Join(assetsFolder))

	go startMQTTClient()

	logrus.Fatal(app.Listen(":3000"))
}
