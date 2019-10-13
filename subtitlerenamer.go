//go:generate goversioninfo -icon=icon.ico -64=true
package main

import (
    "fmt"
    "os"
    "log"
    "path/filepath"
    "io"
    "github.com/conradfr/dialog"
    "github.com/jinzhu/configor"
    "errors"
)

type target struct {
    folder string
    fileName string
    fileExt string
    baseName string
}

const configFileName = "config.yml"

var Config = struct {
    DistFolder string `yaml:"DistFolder"`
}{}

// Not sure of the best or transparent way to handle differences in current path when executed as executable or go run
func getConfigFilePath() (string, error) {
    if _, err := os.Stat(configFileName); os.IsNotExist(err) {
        ex, _ := os.Executable()
        execFolder := filepath.Dir(ex)

        if _, err := os.Stat(execFolder + "/" + configFileName); err == nil {
            return execFolder + "/" + configFileName, nil
        } else {
            return "", errors.New("Config file not found.")
        }
    } else {
        return configFileName, nil;
    }
}

// Manage config file & get destination folder
func getDestinationFolder() string {
    var distFolder string;
    configFilePath, err := getConfigFilePath()
    if err == nil {
        configor.Load(&Config, configFilePath)
        distFolder = Config.DistFolder;
    } else {
        // Use subtitle's folder as default
        distFolder = filepath.Dir(os.Args[1])
    }

    // Ensure distFolder has a trailing backslash
    if (distFolder[len(distFolder) - 1:len(distFolder)] != "\\" ) {
        distFolder = distFolder + "\\";
    }

    return distFolder
}

func getFinalSrtPath(targetVideo string) string {
    targetFolder := filepath.Dir(targetVideo);
    targetFileName := filepath.Base(targetVideo)
    targetFileExt := filepath.Ext(targetVideo)
    targetBaseName := targetFileName[:len(targetFileName)-len(targetFileExt)]

    return targetFolder + "/" + targetBaseName + ".srt"
}

func main() {
    // Ensure a subtitle file has been given
    // TODO: If not, open a file dialog for selecting one
    if len(os.Args) < 2 {
        log.Fatal("No subtitle file provided.")
    }

    // Ensure file exists
    if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
        log.Fatal("Subtitle file not found")
    }

    destinationFolder := getDestinationFolder()
    
    srtFilePath := os.Args[1]
    fileExt := filepath.Ext(srtFilePath)

    if fileExt == ".srt" {
        // Open original subtitle
        in, err := os.Open(srtFilePath)
        if err != nil { log.Fatal(err) }
        defer in.Close()

        // Select file dialog
        videoDialog := dialog.File()
        videoDialog.StartDir = destinationFolder
        //videoDialog.ValidateNames(false)
        targetVideo, err := videoDialog.Title("Select target video").Filter("Video file", "*.mkv;*.avi;*.mp4").Load()

        if err != nil {
            log.Fatal(err)
        }

        // Create destination subtitle
        finalStrPath := getFinalSrtPath(targetVideo)
        out, err := os.Create(finalStrPath)
        if err != nil { log.Fatal(err) }
        defer out.Close()

        // Copy
        _, err = io.Copy(out, in)
        if err != nil { log.Fatal(err) }
    } else {
        log.Fatal("Unsupported filetype." + "\n")
    }

    fmt.Println("Done.")
}
