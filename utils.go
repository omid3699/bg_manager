package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func call(wait bool, cmd string, args ...string) error {
	prc := exec.Command(cmd, args...)
	err := prc.Run()
	if err != nil {
		return err
	}
	if wait {
		prc.Wait()
	}
	if !prc.ProcessState.Success() {
		return prc.Err
	}
	return nil
}

func checkPath(file string) bool {
	path := os.Getenv("PATH")
	pathSplit := filepath.SplitList(path)
	for _, directory := range pathSplit {
		fullPath := filepath.Join(directory, file)
		// Does it exist?
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			mode := fileInfo.Mode()
			// Is it a regular file?
			if mode.IsRegular() {
				// Is it executable?
				if mode&0111 != 0 {
					return true
				}
			}
		}
	}
	return false
}

func newConfig() *Config {
	var backend string
	if checkPath("swaybg") {
		backend = "swaybg"
	} else if checkPath("swww") {
		backend = "swww"
	} else if checkPath("feh") {
		backend = "feh"
	} else {
		log.Fatal("installed wallpaper utility not found in path.\n please install swaybg or swww for wayland and feh for Xorg.")
	}
	return &Config{
		Backend:       backend,
		WaitDelay:     1,
		AcceptFormats: []string{".png", ".jpg", ".jpeg"},
		WallpaperDirs: []string{"/home/omid/wallpapers"},
	}
}

func (manager *Manager) listImages(path string) {
	filepath.Walk(path, func(wPath string, info os.FileInfo, err error) error {
		// Walk the given dir
		// without printing out.
		if wPath == path {
			return nil
		}
		// If given path is folder
		// stop list recursively and print as folder.
		if info.IsDir() {
			manager.listImages(wPath)
			return filepath.SkipDir
		}
		// cehck file extionsion and append to images slcide
		if wPath != path {
			// Check file extension and append to images slice if it matches accepted formats
			for _, format := range manager.config.AcceptFormats {
				if strings.HasSuffix(strings.ToLower(info.Name()), strings.ToLower(format)) {
					manager.images = append(manager.images, wPath)
					break
				}
			}
		}
		return nil
	})

}

func (manager *Manager) changeBg() error {
	if len(manager.images) == 0 {
		return errors.New("no image found")
	}

	for {
		img := manager.images[rand.Intn(len(manager.images))]
		fmt.Println("current bg: ", img)
		go func() {
			switch manager.config.Backend {
			case "feh":
				call(true, "killall", "feh")
				if err := call(false, "feh", "--bg-scale", img); err != nil {
					log.Fatal(err)
				}
			case "swaybg":
				call(true, "killall", "swaybg")
				if err := call(false, "swaybg", "-i", img, "-m", "fill"); err != nil {
					fmt.Println(err)
				}

			case "swww":
				call(false, "swww-daemon")
				if err := call(false, "swww", "img", img); err != nil {
					log.Fatal(err)
				}
			}
		}()
		time.Sleep(time.Minute * time.Duration(manager.config.WaitDelay))
	}
}
