package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	"github.com/brutalzinn/go-mqtt-integration/utils"
	"github.com/brutalzinn/go-mqtt-integration/wshelper"
	"github.com/brutalzinn/go-mqtt-integration/youtube"
	"github.com/sirupsen/logrus"
)

func handleYouTubeCommand(data Youtube) error {
	client, video, err := youtube.GetClient(data.Src)
	if err != nil {
		return err
	}
	title := utils.SanitizeFileName(video.Title)
	var fileName string
	if data.OnlyAudio {
		fileName = title + ".flac"
	} else {
		fileName = title + ".mp4"
	}
	///transfer this verification to filemanager. Filemanager will handle local files, aws files and Ftp files too. Why three providers? Because i need audio backups
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if data.Download {
			if data.OnlyAudio {
				err := youtube.DownloadAudio(client, video, fileName)
				if err != nil {
					logrus.Error("Error on downloing YouTube audio: %w %w\n", data.Src, err.Error())
					return err
				}
			} else {
				err := youtube.DownloadVideo(client, video, fileName)
				if err != nil {
					logrus.Error("Error on downloing YouTube audio: %w %w", data.Src, err.Error())
					return err
				}
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

	switch runtime.GOOS {
	case "darwin":
		exec.Command("osascript", "-e", `display notification "`+message+`"`).Run()
	case "windows":
		exec.Command("powershell", "-Command", `New-BurntToastNotification -Text "`+message+`"`).Run()
	default:
		notifyData := map[string]any{
			"title":   title,
			"message": message,
		}
		wshelper.NotifyClients("notify", notifyData)
	}
	return nil
}

func handleCommand(command Command) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command(command.OS.Mac.Name, command.OS.Mac.Args).Run()
	case "windows":
		return exec.Command(command.OS.Windows.Name, command.OS.Windows.Args).Run()
	default:
		return exec.Command(command.OS.Linux.Name, command.OS.Linux.Args).Run()
	}
}

func handleAudioCommand(audio Audio) error {
	config := confighelper.Get()
	path := filepath.Join(config.AssetsFolder, audio.Name)
	err := utils.DownloadAudio(audio.Src, audio.Name)
	if err != nil {
		return err
	}
	return playAudioCommand(path)
}

func playYoutubeCommand(filePath string, onlyAudio bool) error {
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
		notifyData := map[string]any{
			"src":       fmt.Sprintf("/%v", filePath),
			"onlyAudio": onlyAudio,
		}
		wshelper.NotifyClients("youtube", notifyData)
	}
	return nil
}

func playAudioCommand(url string) error {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("afplay", url).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	default:
		notifyData := map[string]any{
			"src": fmt.Sprintf("/%v", url),
		}
		wshelper.NotifyClients("audio", notifyData)
	}
	return nil
}
