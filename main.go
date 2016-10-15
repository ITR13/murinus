package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  int32 = 640
	screenHeight int32 = 480
	size         int32 = 16
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	e(err)
	fmt.Println("Init SDL")

	window, err := sdl.CreateWindow("Murinus", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight),
		sdl.WINDOW_SHOWN)
	e(err)
	defer window.Destroy()
	fmt.Println("Created window")

	renderer, err := sdl.CreateRenderer(window, -1,
		sdl.RENDERER_ACCELERATED)
	e(err)
	defer renderer.Destroy()
	renderer.Clear()
	fmt.Println("Created renderer")

	stage := LoadTextures(32, 24, renderer)
	fmt.Println("Created loaded stage-basis")
	tiles := make([][]Tile, 32)
	for x := 0; x < 32; x++ {
		tiles[x] = make([]Tile, 24)
		for y := 0; y < 24; y++ {
			tiles[x][y] = Tile((x + y) % 2)
		}
	}
	stage.tiles.tiles = tiles
	fmt.Println("Created tiles")
	for i := 0; i < 60*20; i++ {
		stage.Render(renderer)
		sdl.Delay(17)
		if i%60 == 0 {
			fmt.Printf("Rendered for %d seconds\n", i/60)
		}
	}
	fmt.Println("Exit")
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
