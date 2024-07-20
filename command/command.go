package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
		return playAudio(audio.Src)
	case "notify":
		var notify Notify
		err := json.Unmarshal(event.Data, &notify)
		if err != nil {
			return err
		}
		logrus.Printf("Received notification: %s\n", notify.Title)
		return displayNotification(notify.Title, notify.Message)
	case "command":
		var command Command
		err := json.Unmarshal(event.Data, &command)
		if err != nil {
			return err
		}
		logrus.Printf("Executing command: %s\n", command.Command)
		return executeCommand(command.Command)
	case "youtube":
		var youtube Youtube
		err := json.Unmarshal(event.Data, &youtube)
		if err != nil {
			return err
		}
		logrus.Printf("Playing YouTube video: %s\n", youtube.Src)
		return playYouTube(youtube)
	}
	return nil
}

func playYouTube(data Youtube) error {
	videoID := getYouTubeVideoID(data.Src)
	if videoID == "" {
		logrus.Warn("Invalid YouTube URL")
		return errors.New("invalid youtube url")
	}

	var filePath string
	if data.OnlyAudio {
		filePath = fmt.Sprintf("%s/%s.mp3", os.Getenv("ASSETS_FOLDER"), videoID)
	} else {
		filePath = fmt.Sprintf("%s/%s.mp4", os.Getenv("ASSETS_FOLDER"), videoID)
	}

	if data.Download {
		/// i downlaod only if doesnt exists here,
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			downloadYouTubeContent(data.Src, filePath, data.OnlyAudio)
		}
	}
	return playContent(filePath, data.OnlyAudio)
}

func playAudio(url string) error {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("afplay", url).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	default:
		return errors.New("Unsupported OS for playing audio")
	}
	return nil
}

func displayNotification(title string, message string) error {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("osascript", "-e", `display notification "`+message+`"`).Run()
	case "windows":
		exec.Command("powershell", "-Command", `New-BurntToastNotification -Text "`+message+`"`).Run()
	default:
		return errors.New("Unsupported OS for notifications")
	}
	return nil
}

func executeCommand(command string) error {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("sh", "-c", command).Run()
	case "windows":
		exec.Command("cmd", "/c", command).Run()
	default:
		return errors.New("Unsupported OS for command execution")
	}
	return nil
}
