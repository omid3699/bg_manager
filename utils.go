package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"strings"
	"time"
)

func call(wait bool, cmd string, args ...string) error {
	prc := exec.Command(cmd, args...)
	if wait {
		if err := prc.Run(); err != nil {
			return fmt.Errorf("command failed: %w", err)
		}
	} else {
		if err := prc.Start(); err != nil {
			return fmt.Errorf("command failed: %w", err)
		}
	}
	return nil
}

func checkPath(file string) bool {
	path := os.Getenv("PATH")
	pathSplit := filepath.SplitList(path)
	for _, directory := range pathSplit {
		fullPath := filepath.Join(directory, file)
		// Does it exist and is it executable?
		if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.Mode().IsRegular() && fileInfo.Mode()&0111 != 0 {
			return true
		}
	}
	return false
}

// LoadConfig loads configuration from JSON file
func LoadConfig() *Config {
	var config *Config

	configPath := fmt.Sprintf("%s/.config/bg_manager.json", os.Getenv("HOME"))
	f, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("error in reading config", err)
		config = DefaultConfig()
	}
	err = json.Unmarshal(f, &config)
	if err != nil {
		fmt.Println("error failed to unmarshal config", err)
		config = DefaultConfig()
	}
	SaveConfig(configPath, config)
	fmt.Println("you can edit config in: ", configPath)
	return config

}

// SaveConfig saves the configuration to a JSON file
func SaveConfig(filename string, config *Config) {
	bytes, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		Logger.Printf("could not marshal config: %v\n", err)
		return
	}

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		Logger.Printf("could not write config to file: %v\n", err)
	}
}

// DefaultConfig is our default configuration file
func DefaultConfig() *Config {
	var backend string
	if checkPath("swaybg") {
		backend = "swaybg"
	} else if checkPath("swww") {
		backend = "swww"
	} else if checkPath("feh") {
		backend = "feh"
	} else {
		Logger.Panic("installed wallpaper utility not found in path.\n please install swaybg or swww for wayland and feh for Xorg.")
	}

	return &Config{
		Backend:       backend,
		WaitDelay:     1,
		AcceptFormats: []string{".png", ".jpg", ".jpeg"},
		WallpaperDirs: []string{fmt.Sprintf("%s/wallpapers/", os.Getenv("HOME")), "/usr/share/backgrounds/", "/usr/share/wallpapers/"},
	}
}

func (manager *Manager) listImages(path string) {
	Logger.Printf("scanning %s", path)
	filepath.Walk(path, func(wPath string, info os.FileInfo, err error) error {
		if err != nil {
			Logger.Println(err)
			return nil
		}
		if info.IsDir() && wPath != path {
			manager.listImages(wPath)
			return filepath.SkipDir
		}
		for _, format := range manager.Config.AcceptFormats {
			if strings.HasSuffix(strings.ToLower(info.Name()), strings.ToLower(format)) {
				manager.Images = append(manager.Images, wPath)
				Logger.Println("found:", wPath)
				break
			}
		}
		return nil
	})
}

func (manager *Manager) changeBg() error {
	if len(manager.Images) == 0 {
		return errors.New("no image found")
	}

	for {
		img := manager.Images[rand.Intn(len(manager.Images))]
		fmt.Println("current bg:", img)
		Logger.Println("current bg:", img)

		switch manager.Config.Backend {
		case "feh":
			_ = call(true, "killall", "feh")
			if err := call(false, "feh", "--bg-scale", img); err != nil {
				Logger.Println(err)
			}
		case "swaybg":
			_ = call(true, "killall", "swaybg")
			if err := call(false, "swaybg", "-i", img, "-m", "fill"); err != nil {
				Logger.Println(err)
			}
		case "swww":
			_ = call(false, "swww-daemon")
			if err := call(false, "swww", "img", img); err != nil {
				Logger.Println(err)
			}
		}

		time.Sleep(time.Minute * time.Duration(manager.Config.WaitDelay))
	}
}
