/*
    This file is part of Murinus.

    Murinus is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Murinus is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with Murinus.  If not, see <http://www.gnu.org/licenses/>.
*/
	
package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	sizeMult int32 = 1 //27
	sizeDiv  int32 = 2 //20

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
	newScreenWidth, newScreenHeight = screenWidth, screenHeight

	runtime.LockOSThread()
	err := sdl.Init(sdl.INIT_EVERYTHING)
	e(err)
	fmt.Println("Init SDL")

	window, err := sdl.CreateWindow("Murinus", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.RENDERER_PRESENTVSYNC)
	e(err)
	defer window.Destroy()
	fmt.Println("Created window")

	renderer, err := sdl.CreateRenderer(window, -1,
		sdl.RENDERER_ACCELERATED)
	e(err)
	defer renderer.Destroy()
	renderer.Clear()
	fmt.Println("Created renderer")

	InitText(renderer)
	fmt.Println("Initiated text")
	InitNumbers(renderer)
	fmt.Println("Initiated numbers")

	input := GetInput()
	fmt.Println("Got inputs")

	ReadOptions("options.xml", input)
	defer SaveOptions("options.xml", input)
	fmt.Println("Created options")

	menus := GetMenus(renderer)
	fmt.Println("Created menus")

	stage := LoadTextures(renderer, input)
	fmt.Println("Loaded stage-basis")

	higscores := Read("singleplayer.hs")
	defer higscores.Write("singleplayer.hs")
	fmt.Println("Loaded Highscores")

	defaultName := "\\\\\\"

	for !quit {
		difficulty = -1
		subMenu := -1
	menuLoop:
		for difficulty == -1 && !quit {
			subMenu = menus[0].Run(renderer, input)
			switch subMenu {
			case -1:
				if !quit {
					quit = Arcade
				}
				break menuLoop
			case 0:
				fallthrough
			case 1:
				difficulty = menus[1].Run(renderer, input)
			case 2:
				fmt.Println("Not made yet") //Training
			case 3:
				higscores.Display(renderer, input)
			case 4:
				DoSettings(menus[3], renderer, input)
			case 5:
				fmt.Println("Not made yet") //Credits
			case 6:
				quit = true
			default:
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
				score -= ScoreMult(500)
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
						}
						fmt.Printf("Won in a row counter: %d\n", wonInARow)
						engine = stage.Load(stage.ID+1, true,
							score+ScoreMult(500))
					}
					fmt.Printf("Lives: %d\n", lives)
					Play(engine, window, renderer, int32(lives))
					score = engine.p1.score
					if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
						fmt.Println("Game was quit with exit key")
						break
					}
					for score > extraLivesCounter && extraLivesCounter*2 > extraLivesCounter {
						extraLivesCounter *= 2
						//extraLives++
						//lives++
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
									levelsCleared, difficulty, time.Now()}
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

	numbers.Free()
	for i := 0; i < len(menus); i++ {
		menus[i].Free()
	}
	stage.Free()
	fmt.Println("Quit")
}

func Play(engine *Engine, window *sdl.Window, renderer *sdl.Renderer,
	lives int32) {
	quit = false
	lostLife = false
	engine.Stage.scores.score,
		engine.Stage.scores.lives = int32(engine.p1.score), lives
	for i := 0; i < 30 && !quit; i++ {
		engine.Stage.Render(renderer, false)
		engine.Input.Poll()
		if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
			fmt.Println("Round was quit with exit key")
			return
		}
	}

	for noKeysTouched >= 5 && !quit {
		engine.Stage.Render(renderer, false)
		engine.Input.Poll()
	}
	fmt.Println("Finished starting animation")

	engine.Stage.tiles.renderedOnce = false
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

		engine.Stage.scores.score = int32(engine.p1.score)
		engine.Stage.Render(renderer, true)
		if engine.Stage.pointsLeft <= 0 || lostLife {
			break
		}
	}
	fmt.Println("Exited play loop")

	for i := 0; i < len(engine.snakes); i++ {
		snake := engine.snakes[i]
		snake.head.display = true
		for j := 0; j < len(snake.body); j++ {
			snake.body[j].display = true
		}
		snake.tail.display = true
	}

	if lostLife {
		for i := 0; i < 90 && !quit; i++ {
			engine.Stage.tiles.renderedOnce = false
			engine.Stage.scores.lives = lives - int32(i/15%2)
			engine.p1.entity.display = (i / 15 % 2) == 0
			engine.Stage.Render(renderer, false)
			engine.Input.Poll()
			if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
				fmt.Println("Round was quit with exit key")
				return
			}
		}
	} else {
		for i := 0; i < 30 && !quit; i++ {
			engine.Stage.Render(renderer, false)
			engine.Input.Poll()
		}
	}
	fmt.Println("Finished exit animation")
}

func DoSettings(menu *Menu, renderer *sdl.Renderer, input *Input) {
	menu.Run(renderer, input)
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
