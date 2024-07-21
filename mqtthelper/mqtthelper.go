package mqtthelper

import (
	"encoding/json"
	"os"

	"github.com/brutalzinn/go-mqtt-integration/command"
	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var event *command.Event
	err := json.Unmarshal(msg.Payload(), &event)
	if err != nil {
		logrus.Errorln("Error unmarshaling MQTT message: %v", err)
		return
	}
	logrus.Info("command received")
	err = command.Run(event, msg.Payload())
	if err != nil {
		logrus.Errorln("Error unmarshaling MQTT message: %v", err)
	}
}

func StartMQTTClient(config confighelper.Config) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MQTT.Broker)
	opts.SetClientID(config.MQTT.ClientID)
	opts.SetUsername(config.MQTT.User)
	opts.SetPassword(config.MQTT.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Fatalf("Error connecting to MQTT broker: %s", token.Error())
	}
	if token := client.Subscribe(os.Getenv("MQTT_TOPIC"), 1, nil); token.Wait() && token.Error() != nil {
		logrus.Fatalf("Error subscribing to topic: %v", token.Error())
	}

	logrus.Info("MQTT client started")
	return client
}
