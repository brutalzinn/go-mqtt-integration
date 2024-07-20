package command

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"runtime"
)

func playContent(filePath string, onlyAudio bool) error {
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
		log.Println("Unsupported OS for playing content")
		return errors.New("unsupported OS for playing content")
	}
	return nil
}

func getYouTubeVideoID(url string) string {
	re := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func downloadYouTubeContent(url, outputPath string, onlyAudio bool) error {
	var cmd *exec.Cmd
	if onlyAudio {
		cmd = exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "--output", outputPath, url)
	} else {
		cmd = exec.Command("youtube-dl", "--format", "mp4", "--output", outputPath, url)
	}
	err := cmd.Run()
	return err
}
