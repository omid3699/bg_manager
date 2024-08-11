# a tool for working with background images in tiling window managers like hyprland, i3, awesome ...


# 1 install dependcies
install one of this supported apps

## I recomand swww for Wayland
```sh
paru -S swww 
# or 
apt-get install swww
# or
dnf install swww 
```

## swaybg is another choice for Wayland
```sh
paru -S swaybg 
# or 
apt-get install swaybg
# or
dnf install swaybg 
```

## feh is a good choice for Xorg
```sh
paru -S feh 
# or 
apt-get install feh
# or
dnf install feh 
```

# install
clone the project for github repository
```sh
git clone https://github.com/omid3699/bg_manager.git
```
build and install 
```sh
make install
```
# set config for runing at sartup

## hyprland
paste this line in config file in `~/.config/hypr/hyprland.conf`
```
exec-once = /usr/bin/bg_manager
```