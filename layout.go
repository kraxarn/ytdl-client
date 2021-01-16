package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/cavaliercoder/grab"
	"os"
	"path"
	"strings"
	"time"
)

func createLayout(parent fyne.Window) fyne.CanvasObject {
	urlEntry := widget.NewEntry()
	urlEntry.PlaceHolder = "https://youtu.be/PayvWj2piKg"

	pathEntry := widget.NewEntry()
	if dir, err := defaultDirectory(); err == nil {
		pathEntry.Text = dir
	} else {
		fmt.Printf("failed to get default directory: %s\n", err)
	}

	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 1

	fileName := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Italic: true,
	})

	var button *widget.Button
	button = widget.NewButton("Download", func() {
		if len(urlEntry.Text) <= 0 {
			dialog.ShowError(fmt.Errorf("no url provided"), parent)
			return
		}

		_, err := os.Stat(pathEntry.Text)
		if os.IsNotExist(err) {
			dialog.ShowError(fmt.Errorf("target directory doesn't exist"), parent)
			return
		}

		var videoId string // 11

		if strings.HasPrefix(urlEntry.Text, "https://www.youtube.com/watch?v=") && len(urlEntry.Text) >= 43 {
			videoId = urlEntry.Text[32:43]
		}

		if len(videoId) <= 0 {
			dialog.ShowError(fmt.Errorf("invalid youtube url"), parent)
			return
		}

		button.SetText("Please wait...")
		button.Disable()

		go func() {
			file, err := fileInfo(videoId)
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			err = file.err()
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			urlEntry.Disable()
			pathEntry.Disable()
			button.SetText("Downloading...")
			fileName.SetText(file.fileName())

			err = download(file, path.Join(pathEntry.Text, file.fileName()), progress)
			urlEntry.Enable()
			pathEntry.Enable()
			button.Enable()
			progress.SetValue(1)

			if err != nil {
				button.SetText("Download failed")
				dialog.ShowError(err, parent)
				return
			}

			button.SetText("Download complete")
			time.AfterFunc(time.Second*5, func() {
				button.SetText("Download")
				fileName.SetText("")
				progress.SetValue(0)
			})
		}()
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("youtube-dl", fyne.TextAlignCenter, fyne.TextStyle{
			Bold: true,
		}),
		urlEntry,
		pathEntry,
		layout.NewSpacer(),
		fileName,
		progress,
		button,
	)
}

func download(fileInfo FileInfo, target string, progress *widget.ProgressBar) error {
	client := grab.NewClient()
	request, err := grab.NewRequest(target, fileInfo.url)
	if err != nil {
		return err
	}

	response := client.Do(request)
	ticker := time.NewTicker(time.Millisecond)
	run := true

	for run {
		select {
		// Progress
		case <-ticker.C:
			progress.SetValue(response.Progress())
		// Done
		case <-response.Done:
			if err := response.Err(); err != nil {
				ticker.Stop()
				return err
			}
			ticker.Stop()
			run = false
		}
	}

	return nil
}
