package command

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func Run(event *Event, msg []byte) error {
	switch event.Type {
	case "audio":
		var audio Audio
		err := json.Unmarshal(event.Data, &audio)
		if err != nil {
			return err
		}
		return handleAudioCommand(audio)
	case "notify":
		var notify Notify
		err := json.Unmarshal(event.Data, &notify)
		if err != nil {
			return err
		}
		logrus.Printf("Received notification: %s\n", notify.Title)
		return handleNotificationCommand(notify.Title, notify.Message)
	case "command":
		var command Command
		err := json.Unmarshal(event.Data, &command)
		if err != nil {
			return err
		}
		logrus.Printf("Executing command: %s\n", command.Command)
		return handleCommand(command.Command)
	case "youtube":
		var youtube Youtube
		err := json.Unmarshal(event.Data, &youtube)
		if err != nil {
			return err
		}
		logrus.Printf("Playing YouTube video: %s\n", youtube.Src)
		return handleYouTubeCommand(youtube)
	}
	return nil
}
