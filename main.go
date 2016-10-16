package main

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  int32 = 640
	screenHeight int32 = 480
	size         int32 = 16
)

var quit bool

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
				if (x+y)%3 != 0 {
					tiles[x][y] = Point
				} else {
					tiles[x][y] = Empty
				}
			} else if (x%4 != 2 && y%12 == 0) || x%2 == 0 && y%4 != 1 {
				tiles[x][y] = Wall
			} else {
				if (x+y)%3 != 0 {
					tiles[x][y] = Point
				} else {
					tiles[x][y] = Empty
				}
			}
		}
	}
	stage.tiles.tiles = tiles
	fmt.Println("Created tiles")

	p1 := Player{stage.sprites.GetEntity(1, 1, Player1),
		0, 4, 32, 0}

	engine := GetEngine(&p1, nil, stage,
		stage.sprites.GetSnake(1, 23, 3, &SimpleAI{}, 0, 5, 10*2, 10*4, 100),
		stage.sprites.GetSnake(31, 23, 3, &SimpleAI{}, 0, 5, 10*2, 10*4, 100))

	stage.Render(renderer)
	Play(engine, window, renderer)
	fmt.Println("Exit")
}

func Play(engine *Engine, window *sdl.Window, renderer *sdl.Renderer) {
	quit = false
	for !quit {
		sdl.Delay(17)
		engine.Input.Poll()
		engine.Advance()
		window.SetTitle("Murinus (score: " + strconv.Itoa(int(engine.p1.score)) + ")")
		engine.Stage.Render(renderer)
	}
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
