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
	sizeMult int32 = 3
	sizeDiv  int32 = 2

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

var quit, lostLife bool

var (
	window *sdl.Window
	renderer *sdl.Renderer
	input *Input
	menus []*Menu
	stage *Stage
	highscores Highscores
	defaultName string
)

func main() {
	Init()

	runtime.LockOSThread()
	err := sdl.Init(sdl.INIT_EVERYTHING)
	PanicOnError(err)
	fmt.Println("Init SDL")

	window, err = sdl.CreateWindow("Murinus", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.RENDERER_PRESENTVSYNC)
	PanicOnError(err)
	defer window.Destroy()
	fmt.Println("Created window")

	renderer, err = sdl.CreateRenderer(window, -1,
		sdl.RENDERER_ACCELERATED)
	PanicOnError(err)
	defer renderer.Destroy()
	renderer.Clear()
	fmt.Println("Created renderer")

	InitText(renderer)
	fmt.Println("Initiated text")
	InitNumbers(renderer)
	fmt.Println("Initiated numbers")

	input = GetInput()
	fmt.Println("Got inputs")

	ReadOptions("options.xml", input)
	if !Arcade {
		defer SaveOptions("options.xml", input)
	}
	fmt.Println("Created options")

	menus = GetMenus(renderer)
	fmt.Println("Created menus")

	stage = LoadTextures(renderer, input)
	fmt.Println("Loaded stage-basis")

	highscores = Read("singleplayer.hs", "multiplayer.hs")
	defer highscores.Write("singleplayer.hs", "multiplayer.hs")

	fmt.Println("Loaded Highscores")

	defaultName = "\\\\\\"

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
				highscores.Display(-1, false, renderer, input)
			case 4:
				DoSettings(menus[3], renderer, input)
			case 5:
				creds, src, dst := GetText("Made by ITR   -   Source available on "+
					"github.com/ITR13/murinus", sdl.Color{255, 255, 255, 255},
					newScreenWidth/2, newScreenHeight/2, renderer)
				input.mono.a.down = false
				input.mono.b.down = false
				renderer.SetRenderTarget(nil)
				renderer.SetDrawColor(0, 0, 0, 255)
				renderer.Clear()
				dst.X -= dst.W / 2
				PanicOnError(renderer.Copy(creds, src, dst))
				for !input.mono.a.down && !input.mono.b.down && !quit {
					input.Poll()
					renderer.Present()
				}
				creds.Destroy()
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
				score := -ScoreMult(500)
				wonInARow := -2
				extraLives := 0
				extraLivesCounter := int64(25000)
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
						engine = stage.Load(stage.ID, false, score, subMenu)
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
							score+ScoreMult(500), subMenu)
					}
					fmt.Printf("Lives: %d\n", lives)
					if engine == nil {
						fmt.Println("Engine nil, game was won")
						break
					}
					Play(engine, window, renderer, int32(lives))
					score = engine.Score
					if engine.Input.exit.timeHeld >
						timeExitHasToBeHeldBeforeCloseGame {
						fmt.Println("Game was quit with exit key")
						break
					}
					for score > extraLivesCounter &&
						extraLivesCounter*2 > extraLivesCounter {
						extraLivesCounter *= 2
						//extraLives++
						//lives++
					}
					fmt.Printf("Score: %d\n", score)
				}
				fmt.Printf("Game Over. Final score %d\n", score)
				stage.lostOnce = true
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
								highscores.Add(scoreData, subMenu != 0)
							} else {
								scoreData.Name = name
							}
						}
					} else if menuChoice == 1 {
						highscores.Display(difficulty, subMenu != 0,
							renderer, input)
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

func Init() {
	screenWidth, screenHeight, blockSize, blockSizeBigBoard =
		screenWidthD, screenHeightD, blockSizeD, blockSizeBigBoardD
	newScreenWidth, newScreenHeight = screenWidth, screenHeight
}

func Play(engine *Engine, window *sdl.Window, renderer *sdl.Renderer,
	lives int32) {
	p1C, p2C := options.CharacterP1, options.CharacterP2
	if engine.p1 == nil {
		p1C = p2C
	} else if engine.p2 == nil {
		p2C = p1C
	}

	quit = false
	lostLife = false
	score := int32(0)
	engine.Stage.scores.score, engine.Stage.scores.lives = engine.Score, lives
	for i := 0; i < 30 && !quit; i++ {
		engine.Stage.Render(p1C, p2C, renderer, false)
		engine.Input.Poll()
		if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
			fmt.Println("Round was quit with exit key")
			return
		}
	}

	for noKeysTouched >= 5 && !quit {
		engine.Stage.Render(p1C, p2C, renderer, false)
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
			strconv.Itoa(int(score)) +
			", left " + strconv.Itoa(engine.Stage.pointsLeft) + ")")

		engine.Stage.scores.score = engine.Score
		engine.Stage.Render(p1C, p2C, renderer, true)
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
			p1C, p2C := options.CharacterP1, options.CharacterP2
			if engine.p1 != nil {
				engine.p1.entity.display = (i / 15 % 2) == 0
			} else {
				p1C = p2C
			}

			if engine.p2 != nil {
				engine.p2.entity.display = (i / 15 % 2) == 0
			} else {
				p2C = p1C
			}

			engine.Stage.Render(p1C, p2C, renderer, false)
			engine.Input.Poll()
			if engine.Input.exit.timeHeld > timeExitHasToBeHeldBeforeCloseGame {
				fmt.Println("Round was quit with exit key")
				return
			}
		}
	} else {
		for i := 0; i < 30 && !quit; i++ {
			engine.Stage.Render(p1C, p2C, renderer, false)
			engine.Input.Poll()
		}
	}
	fmt.Println("Finished exit animation")
}

func DoSettings(menu *Menu, renderer *sdl.Renderer, input *Input) {
	for v := menu.Run(renderer, input); v != -1 &&
		!quit; v = menu.Run(renderer, input) {
		ReadOptions("", input)
		menu.menuItems[0].SetNumber(int32(options.CharacterP1), renderer)
		menu.menuItems[1].SetNumber(int32(options.CharacterP2), renderer)
		menu.menuItems[2].SetNumber(int32(options.EdgeSlip), renderer)
		menu.menuItems[3].SetNumber(int32(options.BetterSlip), renderer)
		menu.menuItems[4].SetNumber(int32(options.ShowDivert), renderer)
	}
	if quit {
		return
	}
	options.CharacterP1 = uint8(menu.menuItems[0].numberField.Value)
	options.CharacterP2 = uint8(menu.menuItems[1].numberField.Value)
	options.EdgeSlip = int(menu.menuItems[2].numberField.Value)
	options.BetterSlip = menu.menuItems[3].numberField.Value
	options.ShowDivert = uint8(menu.menuItems[4].numberField.Value)
	options.showDivert = options.ShowDivert != 0
	redrawTextures = true
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
