package confighelper

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// maybe a manifest with yaml can facility us life.
type Config struct {
	StartAt time.Time
	Assets  string
	AWS     AWSConfig
	MQTT    MQTTConfig
	OS      OSConfig
}

type AWSConfig struct {
	Region    string
	Bucket    string
	EndPoint  string
	AccessKey string
	SecretKey string
	IsAWS     bool
}

type MQTTConfig struct {
	Broker   string
	Topic    string
	ClientID string
	User     string
	Password string
}

type OSConfig struct {
	IsContainer bool
}

var config Config

func Get() Config {
	return config
}
func New() Config {
	logrus.Info("Preparing config")
	config = Config{
		Assets: os.Getenv("ASSETS_FOLDER"),
		MQTT: MQTTConfig{
			Broker:   os.Getenv("MQTT_BROKER"),
			Topic:    os.Getenv("MQTT_TOPIC"),
			ClientID: os.Getenv("MQTT_ID_CLIENT"),
			User:     os.Getenv("MQTT_USER"),
			Password: os.Getenv("MQTT_PASSWORD"),
		},
		AWS: AWSConfig{
			Region:    os.Getenv("AWS_REGION"),
			Bucket:    os.Getenv("AWS_BUCKET"),
			EndPoint:  os.Getenv("AWS_ENDPOINT"),
			AccessKey: os.Getenv("AWS_ACCESS_KEY"),
			SecretKey: os.Getenv("AWS_SECRET_KEY"),
			IsAWS:     os.Getenv("AWS_ENABLED") == "true",
		},
	}
	logrus.Info("Config is ready")
	return config
}
