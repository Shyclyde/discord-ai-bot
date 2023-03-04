package openai

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadOpenAIImage(URL string) (string, error) {
	filename := "temp_image.png"
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("can't download image, status: %d", response.StatusCode)
	}

	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}
