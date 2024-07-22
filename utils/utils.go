package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
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
	fm.SetAWS(filemanager.AWS{
		Region: config.AWS.Region,
		Bucket: config.AWS.Bucket,
	})
	fm.Write()
	return nil
}
