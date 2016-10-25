package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	sizeMult int32 = 1 //27
	sizeDiv  int32 = 2 //20
)

const (
	screenWidth  int32 = (1280 * sizeMult) / sizeDiv
	screenHeight int32 = (800 * sizeMult) / sizeDiv
	blockSize    int32 = (48 * sizeMult) / sizeDiv
)

var quit bool
var lostLife bool

func main() {
	runtime.LockOSThread()
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

	stage := LoadTextures(stageWidth, stageHeight, renderer)
	fmt.Println("Loaded stage-basis")

	lostLife = false
	lives := 3
	score := uint64(0)
	score -= 1000
	wonInARow := -2
	extraLives := 0
	extraLivesCounter := uint64(25000)
	for !quit {
		var engine *Engine
		if lostLife {
			wonInARow = -1
			lostLife = false
			lives--
			if lives < extraLives {
				extraLives = lives
			}
			if lives == 0 {
				fmt.Printf("Game Over. Final score %d\n", score)
				lives = 3
				score = 0
				score -= 1000
				wonInARow = -2
				extraLives = 0
				extraLivesCounter = 50000
				engine = stage.Load(0, true, 0)
			} else {
				engine = stage.Load(stage.ID, false, score)
			}
			window.SetTitle("Score: " + strconv.Itoa(int(score)) +
				" Lives: " + strconv.Itoa(lives))
		} else {
			wonInARow++
			if wonInARow == 3 {
				if lives-extraLives < 4 {
					wonInARow = 0
					lives++
				}
			} else if wonInARow == 10 {
				if lives-extraLives < 5 {
					wonInARow = 0
					lives++
				}
			}
			fmt.Printf("Won in a row counter: %d\n", wonInARow)
			engine = stage.Load(stage.ID+1, true, score+1000)
		}
		fmt.Printf("Lives: %d\n", lives)
		Play(engine, window, renderer, int32(lives))
		score = engine.p1.score
		if score > extraLivesCounter {
			extraLivesCounter *= 2
			extraLives++
			lives++
		}
		fmt.Printf("Score: %d\n", score)
	}
	fmt.Println("Exit")
}

func Play(engine *Engine, window *sdl.Window, renderer *sdl.Renderer,
	lives int32) {
	quit = false
	lostLife = false
	for i := 0; i < 90 && !quit; i++ {
		engine.Stage.Render(renderer, lives, int32(engine.p1.score))
		sdl.Delay(17)
		engine.Input.Poll()
	}
	fmt.Println("Finished starting animation")
	for !quit {
		sdl.Delay(17)
		engine.Input.Poll()
		engine.Advance()
		window.SetTitle("Murinus (score: " +
			strconv.Itoa(int(engine.p1.score)) +
			", left " + strconv.Itoa(engine.Stage.pointsLeft) + ")")
		engine.Stage.Render(renderer, lives, int32(engine.p1.score))
		if engine.Stage.pointsLeft <= 0 || lostLife {
			break
		}
	}
	fmt.Println("Exited play loop")
	if lostLife {
		for i := 0; i < 90 && !quit; i++ {
			engine.Stage.Render(renderer, lives-int32(i/15%2),
				int32(engine.p1.score))
			sdl.Delay(17)
			engine.Input.Poll()
		}
	} else {
		for i := 0; i < 30 && !quit; i++ {
			engine.Stage.Render(renderer, lives, int32(engine.p1.score))
			sdl.Delay(17)
			engine.Input.Poll()
		}
	}
	fmt.Println("Finished exit animation")
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
