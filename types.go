package main

import "log"

type Config struct {
	Backend       string   `json:"backend"`
	WaitDelay     int      `json:"wait_delay"`
	AcceptFormats []string `json:"accpet_formats"`
	WallpaperDirs []string `json:"wallpaper_dirs"`
}

type Manager struct {
	Config *Config
	Images []string
	Logger *log.Logger
}
