package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	sizeMult int32 = 1 //27
	sizeDiv  int32 = 1 //20

	timeExitHasToBeHeldBeforeGameEnd   int = 60 * 5
	timeExitHasToBeHeldBeforeCloseGame int = 90
)

const (
	screenWidthD       int32 = (1280 * sizeMult) / sizeDiv
	screenHeightD      int32 = (800 * sizeMult) / sizeDiv
	blockSizeD         int32 = (48 * sizeMult) / sizeDiv
	blockSizeBigBoardD int32 = (24 * sizeMult) / sizeDiv
	gSize              int32 = 12
)

var screenWidth, screenHeight, blockSize, blockSizeBigBoard int32

var quit bool
var lostLife bool

func main() {
	screenWidth, screenHeight, blockSize, blockSizeBigBoard =
		screenWidthD, screenHeightD, blockSizeD, blockSizeBigBoardD

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

	input := GetInput()
	fmt.Println("Got inputs")

	ReadOptions("options.xml", input)
	defer SaveOptions("options.xml", input)
	fmt.Println("Created options")

	stage := LoadTextures(renderer, input)
	fmt.Println("Loaded stage-basis")

	menus := GetMenus(renderer)
	higscores := Read("singleplayer.hs")
	defer higscores.Write("singleplayer.hs")
	fmt.Println("Loaded Highscores")

	defaultName := "\\\\\\"

	for !quit {
		difficulty = -1
		subMenu := -1
		for difficulty == -1 && !quit {
			subMenu = menus[0].Run(renderer, input)
			if subMenu == -1 && !quit {
				quit = Arcade
				break
			} else if subMenu == 0 || subMenu == 1 {
				difficulty = menus[1].Run(renderer, input)
			} else if subMenu == 2 {
				higscores.Display(renderer, input)
			} else if subMenu == 3 {
				fmt.Println("Not made yet")
			} else if subMenu == 4 {
				quit = true
			} else {
				panic("Unknown menu option")
			}
		}
		if quit {
			break
		}

		stage.ID = -1
		if subMenu == 0 || subMenu == 1 {
			for !quit {
				lostLife = false
				lives := 3
				score := uint64(0)
				score -= 500 * (uint64(difficulty*difficulty) + 1)
				wonInARow := -2
				extraLives := 0
				extraLivesCounter := uint64(25000)
				levelsCleared := 0
				for !quit && (lives != 1 || !lostLife) {
					var engine *Engine
					if lostLife {
						wonInARow = -1
						lostLife = false
						lives--
						if lives < extraLives {
							extraLives = lives
						}
						if lives == 0 {
							panic("Should not reach this statement")
						}
						engine = stage.Load(stage.ID, false, score)
						window.SetTitle("Score: " + strconv.Itoa(int(score)) +
							" Lives: " + strconv.Itoa(lives))
					} else {
						levelsCleared++
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
						engine = stage.Load(stage.ID+1, true, score+500*
							(uint64(difficulty*difficulty)+1))
					}
					fmt.Printf("Lives: %d\n", lives)
					Play(engine, window, renderer, int32(lives))
					score = engine.p1.score
					if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
						fmt.Println("Game was quit with exit key")
						break
					}
					if score > extraLivesCounter {
						extraLivesCounter *= 2
						extraLives++
						lives++
					}
					fmt.Printf("Score: %d\n", score)
				}
				fmt.Printf("Game Over. Final score %d\n", score)
				input.exit.timeHeld = 0

				menuChoice := -1
				var scoreData *ScoreData
				menus[2].selectedElement = 0
				for !quit && menuChoice < 2 {
					menuChoice = menus[2].Run(renderer, input)
					if menuChoice == 0 {
						name := GetName(defaultName, renderer, input)
						if name != "" {
							defaultName = name
							if scoreData == nil {
								scoreData = &ScoreData{score, name,
									levelsCleared, difficulty}
								higscores.Add(scoreData)
								higscores.Sort()
							} else {
								scoreData.Name = name
							}
						}
					} else if menuChoice == 1 {
						higscores.Display(renderer, input)
					} else if menuChoice == -1 {
						menuChoice = 4
					}
				}
				if quit {
					break
				} else if menuChoice == 2 {
					stage.ID--
				} else if menuChoice == 3 {
					stage.ID = -1
				} else if menuChoice == 4 {
					break
				} else {
					panic("Unknown menu option")
				}
			}
		}
	}
	fmt.Println("Quit")
}

func Play(engine *Engine, window *sdl.Window, renderer *sdl.Renderer,
	lives int32) {
	quit = false
	lostLife = false
	for i := 0; i < 30 && !quit; i++ {
		engine.Stage.Render(renderer, lives, int32(engine.p1.score))
		engine.Input.Poll()
		if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
			fmt.Println("Round was quit with exit key")
			return
		}
	}

	for noKeysTouched >= 5 && !quit {
		engine.Stage.Render(renderer, lives, int32(engine.p1.score))
		engine.Input.Poll()
	}

	fmt.Println("Finished starting animation")
	for !quit {
		engine.Input.Poll()
		if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
			fmt.Println("Round was quit with exit key")
			return
		}

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
			engine.Input.Poll()
			if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
				fmt.Println("Round was quit with exit key")
				return
			}
		}
	} else {
		for i := 0; i < 30 && !quit; i++ {
			engine.Stage.Render(renderer, lives, int32(engine.p1.score))
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
