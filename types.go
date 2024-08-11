package main

// Config struct of application
type Config struct {
	Backend       string   `json:"backend"`
	WaitDelay     int      `json:"wait_delay"`
	AcceptFormats []string `json:"accpet_formats"`
	WallpaperDirs []string `json:"wallpaper_dirs"`
}

// Manager have one config and one slice of background images
type Manager struct {
	Config *Config
	Images []string
}
