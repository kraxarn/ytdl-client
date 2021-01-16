package main

import "fmt"

const ApplicationName string = "ytdl-client"
const CurrentVersion string = "v1.1"

func ApplicationTitle() string {
	return fmt.Sprintf("ytdl-client %s", CurrentVersion)
}
