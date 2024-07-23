package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	"github.com/brutalzinn/go-mqtt-integration/filemanager"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// /thanks to kennygrant to provide this lib for sanitize https://github.com/kennygrant/sanitize
func SanitizeFileName(filename string) string {
	filePath := strings.ToLower(filename)
	filePath = strings.Replace(filePath, "..", "_", -1)
	filePath = strings.Replace(filePath, " ", "", -1)
	filePath = strings.Replace(filePath, "-", "", -1)
	filePath = path.Clean(filePath)
	// Remove illegal characters for paths, flattening accents
	// and replacing some common separators with -
	b := bytes.NewBufferString("")
	for _, c := range filePath {
		if val, ok := transliterations[c]; ok {
			b.WriteString(val)
		} else {
			b.WriteRune(c)
		}
	}
	return filePath
}

// /amazon recommends to use mp3 with 48k bitrate and 24khz sample rate
func ConvertMP3ForAlexa(input []byte) ([]byte, error) {
	inputFile, err := os.CreateTemp("", "input-*.mp3")
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()
	defer os.Remove(inputFile.Name())
	_, err = inputFile.Write(input)
	if err != nil {
		return nil, err
	}
	outputFile, err := os.CreateTemp("", "output-*.mp3")
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()
	defer os.Remove(outputFile.Name())
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile.Name(), "-ac", "2", "-codec:a", "libmp3lame", "-b:a", "48k", "-ar", "16000", outputFile.Name())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %v, output: %s", err, out.String())
	}
	convertedBytes, err := os.ReadFile(outputFile.Name())
	if err != nil {
		return nil, err
	}
	return convertedBytes, nil
}

func DownloadAudio(url string, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download audio: received non-200 status code %d", response.StatusCode)
	}
	config := confighelper.Get()
	fm := filemanager.New(filename)
	fm.SetReader(response.Body)
	fm.SetAWS(&filemanager.AWS{
		Enabled: config.AWS.Enabled,
		Bucket:  config.AWS.Bucket,
		Region:  config.AWS.Region,
	})
	err = fm.Write()
	if err != nil {
		return err
	}
	return nil
}
