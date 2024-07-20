package command

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func serveContent(src string) {
	contentURL := "/content/" + filepath.Base(src)
	// Write the content URL to a file to be served by the web server
	err := os.WriteFile("/app/content_url.txt", []byte(contentURL), 0644)
	if err != nil {
		log.Printf("Error writing content URL: %v\n", err)
		return
	}
}

func getContentURL() string {
	data, err := ioutil.ReadFile("/app/content_url.txt")
	if err != nil {
		log.Printf("Error reading content URL: %v\n", err)
		return ""
	}
	return string(data)
}
