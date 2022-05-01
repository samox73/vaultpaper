package main

type backend interface {
	setWallpaper(name string) error
}
