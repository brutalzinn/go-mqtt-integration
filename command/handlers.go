package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/brutalzinn/go-mqtt-integration/utils"
	"github.com/brutalzinn/go-mqtt-integration/wshelper"
	"github.com/brutalzinn/go-mqtt-integration/youtube"
	"github.com/sirupsen/logrus"
)

var isCompatiblePlayback = os.Getenv("OS_TYPE") == "unsupported"

func handleYouTubeCommand(data Youtube) error {
	client, video, err := youtube.GetClient(data.Src)
	if err != nil {
		return err
	}
	title := utils.SanitizeFileName(video.Title)
	var path string
	if data.OnlyAudio {
		path = os.Getenv("ASSETS_FOLDER") + "/" + title + ".flc"
	} else {
		path = os.Getenv("ASSETS_FOLDER") + "/" + title + ".mp4"
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if data.Download {
			if data.OnlyAudio {
				youtube.DownloadAudio(client, video, path)
			} else {
				youtube.DownloadVideo(client, video, path)
			}
		}
	}

	err = playYoutubeCommand(path, data.OnlyAudio)
	if err != nil {
		logrus.Printf("Playing YouTube video: %s\n", data.Src)
		return err
	}
	return nil
}

func handleNotificationCommand(title string, message string) error {
	if isCompatiblePlayback {
		notifyData := map[string]any{
			"title":   title,
			"message": message,
		}
		///incorrect way. But now i can identify that i am using this project as hyperfocus too. why my mind do this?!!
		wshelper.NotifyClients("notify", notifyData)
		return errors.New("is this app running at container? can only reproduce on web")
	}
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

func handleCommand(command string) error {
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

func handleAudioCommand(audio Audio) error {
	path := filepath.Join(os.Getenv("ASSETS_FOLDER"), audio.Name)
	err := utils.DownloadAudio(audio.Src, audio.Name)
	if err != nil {
		return err
	}
	return playAudioCommand(path)
}

func playYoutubeCommand(filePath string, onlyAudio bool) error {
	if isCompatiblePlayback {
		notifyData := map[string]any{
			"src":       fmt.Sprintf("/%v", filePath),
			"onlyAudio": onlyAudio,
		}
		///incorrect way. But now i can identify that i am using this project as hyperfocus too. why my mind do this?!!
		wshelper.NotifyClients("youtube", notifyData)
		return errors.New("is this app running at container? can only reproduce on web")
	}
	switch runtime.GOOS {
	case "darwin":
		if onlyAudio {
			exec.Command("afplay", filePath).Start()
		} else {
			exec.Command("open", filePath).Start()
		}
	case "windows":
		exec.Command("cmd", "/c", "start", filePath).Start()
	default:
		logrus.Println("Unsupported OS for playing content")
		return errors.New("unsupported OS for playing content")
	}
	return nil
}

func playAudioCommand(url string) error {
	if isCompatiblePlayback {
		notifyData := map[string]any{
			"src": fmt.Sprintf("/%v", url),
		}
		///incorrect way. But now i can identify that i am using this project as hyperfocus too. why my mind do this?!!
		wshelper.NotifyClients("audio", notifyData)
		return errors.New("is this app running at container? can only reproduce on web")
	}
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
