package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io"
	"github.com/sqweek/dialog"
	"github.com/jinzhu/configor"
)

var Config = struct {
	DistFolder string `yaml:"DistFolder"`
}{}

func main() {
	// Ensure a subtitle file has been given
	// TODO: If not, open a file dialog for selecting one
	if len(os.Args) < 2 {
		log.Fatal("No file provided.")
	}

	// Load conf & ensured distFolder has a trailing backslash
	configor.Load(&Config, "config.yml")
	var distFolder = Config.DistFolder;
	if (distFolder[len(distFolder) - 1:len(distFolder)] != "\\" ) {
		distFolder = distFolder + "\\";
	}

	var srtFilePath string = os.Args[1]
	fileExt := filepath.Ext(srtFilePath)

	if fileExt == ".srt" {
		// Select file dialog
		var videoDialog = dialog.File()
		videoDialog.StartDir = distFolder
		targetVideo, err := videoDialog.Title("Select target video").Filter("Video file", "*.mkv;*.avi;*.mp4").Load()

		if err != nil {
			log.Fatal(err)
		}

		targetFileName := filepath.Base(targetVideo)
		targetFileExt := filepath.Ext(targetVideo)
		targetBaseName := targetFileName[:len(targetFileName)-len(targetFileExt)]

		// Open original subtitle
		in, err := os.Open(srtFilePath)
		if err != nil { log.Fatal(err) }
		defer in.Close()

		// Create dest subtitle
		out, err := os.Create(distFolder + targetBaseName + ".srt")
		if err != nil { log.Fatal(err) }
		defer out.Close()

		// Copy
		_, err = io.Copy(out, in)
		out.Close()
		if err != nil { log.Fatal(err) }
	} else {
		log.Fatal("Unsupported filetype." + "\n")
	}

	fmt.Println("Done.")
}
