package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
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
		// Check transliterations first
		if val, ok := transliterations[c]; ok {
			b.WriteString(val)
		} else {
			b.WriteRune(c)
		}
	}
	// NB this may be of length 0, caller must check
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
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to save audio to file: %w", err)
	}
	return nil
}
