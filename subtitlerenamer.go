package main

import (
    "fmt"
    "os"
    "log"
    "path/filepath"
    "io"
    "github.com/sqweek/dialog"
    "github.com/jinzhu/configor"
    "github.com/kardianos/osext"
)

const configFileName = "config.yml"

var Config = struct {
    DistFolder string `yaml:"DistFolder"`
}{}

func main() {
    // Ensure a subtitle file has been given
    // TODO: If not, open a file dialog for selecting one
    if len(os.Args) < 2 {
        log.Fatal("No file provided.")
    }


    // Find the config file
    // Not sure of the best or transparent way to handle differences in current path when executed as executable or go run
    var configFilePath string;
    if _, err := os.Stat(configFileName); os.IsNotExist(err) {
        var execFolder, _ = osext.ExecutableFolder()
        if _, err := os.Stat(execFolder + "/" + configFileName); err == nil {
            configFilePath = execFolder + "/" + configFileName;
        } else {
            log.Fatal("Config file not found.");
        }

    } else {
        configFilePath = configFileName;
    }

    // Load conf & ensured distFolder has a trailing backslash
    configor.Load(&Config, configFilePath)

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

        targetFolder := filepath.Dir(targetVideo);
        targetFileName := filepath.Base(targetVideo)
        targetFileExt := filepath.Ext(targetVideo)
        targetBaseName := targetFileName[:len(targetFileName)-len(targetFileExt)]

        // Open original subtitle
        in, err := os.Open(srtFilePath)
        if err != nil { log.Fatal(err) }
        defer in.Close()

        // Create dest subtitle
        out, err := os.Create(targetFolder + "/" + targetBaseName + ".srt")
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
