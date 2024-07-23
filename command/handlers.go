package command

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	"github.com/brutalzinn/go-mqtt-integration/filemanager"
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
	var filename string
	if data.OnlyAudio {
		filename = video.ID + ".mp3"
	} else {
		filename = video.ID + ".mp4"
	}
	// if _, err := os.Stat(filename); os.IsNotExist(err) {
	if data.Download {
		if data.OnlyAudio {
			///if this will be played at Alexa..
			///amazon recommends to use mp3 with 48k bitrate and 24khz sample rate
			filesBytes, err := youtube.DownloadAudio(client, video, filename)
			if err != nil {
				return err
			}
			///convert to 48k bitrate and 24khz sample rate
			convertedBytes, err := utils.ConvertMP3ForAlexa(filesBytes)
			if err != nil {
				return err
			}
			config := confighelper.Get()
			fm := filemanager.New(filename)
			fm.SetBytes(convertedBytes)
			fm.SetAWS(&filemanager.AWS{
				Enabled: config.AWS.Enabled,
				Bucket:  config.AWS.Bucket,
				Region:  config.AWS.Region,
			})
			err = fm.Write()
			if err != nil {
				return err
			}
		} else {
			filesBytes, err := youtube.DownloadVideo(client, video, filename)
			if err != nil {
				return err
			}
			config := confighelper.Get()
			fm := filemanager.New(filename)
			fm.SetBytes(filesBytes)
			fm.SetAWS(&filemanager.AWS{
				Enabled: config.AWS.Enabled,
				Bucket:  config.AWS.Bucket,
				Region:  config.AWS.Region,
			})
			err = fm.Write()
			if err != nil {
				return err
			}
		}
	}
	// }
	err = playYoutubeCommand(filename, data.OnlyAudio)
	if err != nil {
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
