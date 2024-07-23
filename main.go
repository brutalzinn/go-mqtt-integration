package main

import (
	"os"
	"path/filepath"

	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	"github.com/brutalzinn/go-mqtt-integration/mqtthelper"
	"github.com/brutalzinn/go-mqtt-integration/wshelper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting GO MQTT INTEGRATION PROJECT")
	config := confighelper.New()
	if _, err := os.Stat(config.AssetsFolder); os.IsNotExist(err) {
		os.Mkdir(config.AssetsFolder, os.ModePerm)
		logrus.Info("Dir for assets is created at %s", config.AssetsFolder)
	}
	//// Check if AWS is enabled
	//// yeep. I know. I should have used a flag
	if config.AWS.Enabled == false {
		logrus.Warn("AWS is not enabled. Files is stored locally and cannot be accessed from the HO and Alexa devices.")
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"Title":     "GO MQTT INTEGRATION",
			"Body":      "Welcome to the GO MQTT INTEGRATION",
			"Assets":    "/assets",
			"WS":        "/ws",
			"Content":   "/content-url",
			"AWS":       config.AWS.Enabled,
			"StartTime": config.StartAt.Format("01/02/2006 15:04:05"),
		})
	})
	app.Get("/content-url", func(c *fiber.Ctx) error {
		files, err := filepath.Glob(filepath.Join(config.AssetsFolder, "*"))
		if err != nil {
			logrus.Error("WebSocket error:", err)
			return err
		}
		urls := make([]string, len(files))
		for i, file := range files {
			urls[i] = filepath.Join(config.AssetsFolder, filepath.Base(file))
		}
		return c.JSON(urls)
	})
	app.Get("/ws", websocket.New(wshelper.WebSocketHandler))
	app.Static("/assets", filepath.Join(config.AssetsFolder))
	go mqtthelper.StartMQTTClient(config)
	logrus.Fatal(app.Listen(":3000"))
}
