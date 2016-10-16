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

	stage := LoadTextures(33, 25, renderer)
	fmt.Println("Created loaded stage-basis")
	tiles := make([][]Tile, 33)
	for x := 0; x < 33; x++ {
		tiles[x] = make([]Tile, 25)
		for y := 0; y < 25; y++ {
			if x == 0 || y == 0 || x == 32 || y == 24 {
				tiles[x][y] = Wall
			} else if x == 1 || y == 1 || y == 11 || x == 31 || y == 23 {
				tiles[x][y] = Point
			} else if (x%4 != 2 && y%12 == 0) || x%2 == 0 && y%4 != 1 {
				tiles[x][y] = Wall
			} else {
				tiles[x][y] = Point
			}
		}
	}
	stage.tiles.tiles = tiles
	fmt.Println("Created tiles")

	engine := Engine{nil, []*Snake{
		stage.sprites.GetSnake(1, 23, 3, &SimpleAI{}, 0, 5, 10*2, 10*4, 100),
		stage.sprites.GetSnake(31, 23, 3, &SimpleAI{}, 0, 5, 10*2, 10*4, 100)},
		stage}

	stage.Render(renderer)
	for i := uint64(0); true; i++ {
		sdl.Delay(17)
		if i >= 30 {
			engine.Advance()
			stage.Render(renderer)
			if i%60 == 0 {
				fmt.Printf("Rendered for %d seconds\n", i/60)
			}
			for sdl.PollEvent() != nil {
			}
		}
	}
	fmt.Println("Exit")
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
