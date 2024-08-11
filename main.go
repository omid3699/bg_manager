package main

import (
	"log"
	"os"
	"path"
)

// Logger is a log.Logger instance
var Logger *log.Logger

func main() {
	LogFile := path.Join(os.TempDir(), "bg_manager.log")
	f, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	manger := &Manager{
		Config: LoadConfig(),
		Images: []string{},
	}

	Logger = log.New(f, "", log.LstdFlags)

	for _, dir := range manger.Config.WallpaperDirs {
		manger.listImages(dir)
	}
	Logger.Printf("%d images found from %d directories\n", len(manger.Images), len(manger.Config.WallpaperDirs))
	manger.changeBg()

}
