package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FileInfo struct {
	title     string
	extension string
	url       string
}

func (i *FileInfo) fileName() string {
	return fmt.Sprintf("%s.%s", i.title, i.extension)
}

func (i *FileInfo) err() error {
	if len(i.title) <= 0 || len(i.extension) <= 0 {
		return fmt.Errorf("no information found")
	}
	if len(i.url) <= 0 {
		return fmt.Errorf("no audio found to download")
	}

	return nil
}

const BaseUrl string = "https://kraxarn.com/yt/info/%s"

func fileInfo(videoId string) (FileInfo, error) {
	response, err := http.Get(fmt.Sprintf(BaseUrl, videoId))
	if err != nil {
		return FileInfo{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return FileInfo{}, err
	}

	err = response.Body.Close()
	if err != nil {
		return FileInfo{}, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return FileInfo{}, err
	}

	if errStr, ok := data["error"].(string); ok {
		return FileInfo{}, fmt.Errorf(errStr)
	}

	var file FileInfo
	if title, ok := data["title"].(string); ok {
		file.title = title
	}
	if audio, ok := data["audio"].(map[string]interface{}); ok {
		if extension, ok := audio["codec"].(string); ok {
			file.extension = extension
		}
		if url, ok := audio["url"].(string); ok {
			file.url = url
		}
	}

	return file, nil
}
