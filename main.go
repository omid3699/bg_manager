package main

func main() {
	manger := &Manager{
		config: newConfig(),
		images: []string{},
	}

	for _, dir := range manger.config.WallpaperDirs {
		manger.listImages(dir)
	}
	manger.changeBg()

}
