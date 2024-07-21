package youtube

import (
	"io"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"

	"github.com/sirupsen/logrus"
)

// /thanks to @kkdai github.com/kkdai/youtube to provide this awesome lib ><
func GetBestHighFormat(formats []youtube.Format) youtube.Format {
	var bestFormat youtube.Format
	for _, format := range formats {
		if format.Bitrate > bestFormat.Bitrate {
			bestFormat = format
		}
	}
	return bestFormat
}

func FilterFormats(formats youtube.FormatList, kind string) []youtube.Format {
	var filteredFormats []youtube.Format
	for _, format := range formats {
		if strings.Contains(format.MimeType, kind) {
			filteredFormats = append(filteredFormats, format)
		}
	}
	return filteredFormats
}

func DownloadStream(client *youtube.Client, video *youtube.Video, format *youtube.Format, filepath string) {
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		logrus.Error(err)
	}
	defer stream.Close()

	file, err := os.Create(filepath)
	if err != nil {
		logrus.Error(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		logrus.Error(err)
	}
}

func DownloadAudio(client *youtube.Client, video *youtube.Video, path string) {
	formats_highest_a := GetBestHighFormat(FilterFormats(video.Formats, "audio"))
	DownloadStream(client, video, &formats_highest_a, path)
}

func DownloadVideo(client *youtube.Client, video *youtube.Video, path string) {
	formats_highest_v := GetBestHighFormat(FilterFormats(video.Formats, "video"))
	DownloadStream(client, video, &formats_highest_v, path)
}

func GetClient(url string) (*youtube.Client, *youtube.Video, error) {
	client := &youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		logrus.Error("Error getting YouTube video:", err)
		return nil, nil, err
	}
	return client, video, nil
}
