package main

// /home/thomas/Pictures/FRS/PXL_20240305_090139715.MP.jpg
// /home/thomas/Pictures/FRS/PXL_20240305_090144477.MP.jpg
// /home/thomas/Desktop
//

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/local"
	"goStreamer/modules/ui"
	"goStreamer/modules/web"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var server web.Server

func init() {
	config.Config.Init("config.json")
	os.MkdirAll(config.Config.Local.SourceFolder, os.ModePerm)
	os.MkdirAll(config.Config.Local.Targetfolder, os.ModePerm)
	os.MkdirAll(config.Config.Local.OutputFolder, os.ModePerm)

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := ui.New("GoStreamer")
	var content *fyne.Container

	if !config.Config.Local.Webcam.Enable { // We expect files then
		sourceEntry, sourceButton := ui.AddFolderSelector("Select a source folder", "Choose a folder...")
		targetEntry, targetButton := ui.AddFolderSelector("Select a target folder", "Choose a folder...")
		outputEntry, outputButton := ui.AddFolderSelector("Select output Folder", "Choose an output folder...")
		sourceEntry.Text = config.Config.Local.SourceFolder
		targetEntry.Text = config.Config.Local.Targetfolder
		outputEntry.Text = config.Config.Local.OutputFolder
		submitButton := ui.AddSubmitButton("Submit", func() {

			// Update files and config
			local.Files.Update(sourceEntry.Text, targetEntry.Text, outputEntry.Text)

			server.Connect(config.Config.Server.IP, config.Config.Server.DialPort)
			defer server.Conn.Close()
			ui.HandleUI(&server, ctx, -1)
		})

		getFileButton := ui.AddSubmitButton("Get swapped", func() {

			server.Connect(config.Config.Server.IP, config.Config.Server.DialPort)
			defer server.Conn.Close()
			getFile(ctx)
		})
		content = container.NewVBox(
			sourceEntry, sourceButton,
			targetEntry, targetButton,
			outputEntry, outputButton,
			submitButton,
			getFileButton,
		)
	} else { // We got webcam
		sourceEntry, sourceButton := ui.AddFileSelector("Select a source face", "Choose a file...")

		webcamTarget := ui.AddOutputFilename("Filename", "Enter webgam target (default is usually 0)")

		submitButton := ui.AddSubmitButton("Submit", func() {

			source, err := strconv.Atoi(webcamTarget.Text)
			if err != nil {
				log.Println("Wrong webcam type!")
			}
			local.Files.UpdateSingle(sourceEntry.Text, webcamTarget.Text)
			ui.HandleUI(&server, ctx, source)
		})

		content = container.NewVBox(
			webcamTarget,
			sourceEntry, sourceButton,
			submitButton,
		)
	}
	// Start UI
	ui.Run(content)

	// Connection to the server running face swapper

}
