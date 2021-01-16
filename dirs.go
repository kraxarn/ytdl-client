package main

import (
	"os"
	"path"
)

func defaultDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	music := path.Join(home, "Music")
	_, err = os.Stat(music)
	if err != nil {
		return home, nil
	}

	return music, nil
}
