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
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Arcade bool = false

	sizeMult int32 = 1
	sizeDiv  int32 = 2

	timeExitHasToBeHeldToQuit int = 60 * 5
	timeExitHasToBeHeldToExit int = 90
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
	window      *sdl.Window
	renderer    *sdl.Renderer
	input       *Input
	menus       []*Menu
	stage       *Stage
	highscores  Highscores
	defaultName string
)

func main() {
	Init()
	defer CleanUp()

	for !quit {
		difficulty = -1
		menuChoice := -1
	menuLoop:
		for difficulty == -1 && !quit {
			menuChoice, _, _ = menus[0].Run(renderer, input)
			switch menuChoice {
			case -1:
				if !quit {
					quit = Arcade
				}
				break menuLoop
			case 0:
				fallthrough
			case 1:
				StartGameSession(menuChoice)
			case 2:
				StartTrainingSession()
			case 3:
				highscores.Display(-1, false, renderer, input)
			case 4:
				DoSettings(menus[3], renderer, input)
			case 5:
				ShowCredits()
			case 6:
				quit = true
			default:
				panic("Unknown menu option")
			}
		}
	}
	fmt.Println("Quit")
}

func Init() {
	screenWidth, screenHeight, blockSize, blockSizeBigBoard =
		screenWidthD, screenHeightD, blockSizeD, blockSizeBigBoardD
	newScreenWidth, newScreenHeight = screenWidth, screenHeight

	runtime.LockOSThread()
	err := sdl.Init(sdl.INIT_EVERYTHING)
	PanicOnError(err)
	fmt.Println("Init SDL")

	window, err = sdl.CreateWindow("Murinus", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int(screenWidth), int(screenHeight),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.RENDERER_PRESENTVSYNC)
	PanicOnError(err)
	fmt.Println("Created window")

	renderer, err = sdl.CreateRenderer(window, -1,
		sdl.RENDERER_ACCELERATED)
	PanicOnError(err)
	renderer.Clear()
	fmt.Println("Created renderer")

	InitText(renderer)
	fmt.Println("Initiated text")
	InitNumbers(renderer)
	fmt.Println("Initiated numbers")

	input = GetInput()
	fmt.Println("Got inputs")

	ReadOptions("options.xml", input)
	fmt.Println("Created options")

	stage = LoadTextures(renderer, input)
	fmt.Println("Loaded stage-basis")

	menus = GetMenus(renderer)
	fmt.Println("Created menus")

	highscores = Read("singleplayer.hs", "multiplayer.hs")

	fmt.Println("Loaded Highscores")

	defaultName = "\\\\\\\\\\"
	rand.Seed(time.Now().Unix())
}

func CleanUp() {
	highscores.Write("singleplayer.hs", "multiplayer.hs")
	if !Arcade {
		SaveOptions("options.xml", input)
	}
	numbers.Free()
	for i := 0; i < len(menus); i++ {
		menus[i].Free()
	}
	stage.Free()
	renderer.Destroy()
	window.Destroy()
}

func StartGameSession(menuChoice int) {
	difficulty, _, _ = menus[1].Run(renderer, input)
	stage.ID = -1
	for !quit && difficulty != -1 {
		levelsCleared := 0
		score := -ScoreMult(500)

		RunGame(menuChoice, &levelsCleared, &score)

		fmt.Printf("Game Over. Final score %d\n", score)
		stage.lostOnce = true
		input.exit.timeHeld = 0

		if !GameOverMenu(levelsCleared, score) {
			break
		}
	}
}

func StartTrainingSession() {
	difficulty = -1
	for play, ups, downs := menus[4].Run(renderer, input); !quit &&
		play != -1; play, _, _ = menus[4].Run(renderer, input) {
		ID := stage.levels[2][int(menus[4].NVal(0))][0]
		difficulty := int(menus[4].NVal(1))
		players := menus[4].NVal(3) - 1

		engine := stage.LoadSingleLevel(ID, difficulty,
			ups > 8, downs > 8, true, 0, players)
		for i := menus[4].NVal(2); i > 0; i-- {
			PlayStage(engine, window, renderer, i)
			if quit || input.exit.active || !lostLife {
				break
			}
			engine = stage.LoadSingleLevel(ID, difficulty,
				ups > 8, downs > 8, false, engine.Score, players)
		}
	}
}

func ShowCredits() {
	strings := []string{"Made by ITR",
		"Source available on github.com/ITR13/murinus", " ",
		"Other contributers:", "byllgrim"}

	textures := make([]*sdl.Texture, len(strings))
	src := make([]*sdl.Rect, len(strings))
	dst := make([]*sdl.Rect, len(strings))
	h := int32(0)

	for i := 0; i < len(textures); i++ {
		col := rand.Int()%26 + 1
		red := uint8((col % 3) * 255 / 2)
		col /= 3
		green := uint8((col % 3) * 255 / 2)
		col /= 3
		blue := uint8((col % 3) * 255 / 2)

		textures[i], src[i], dst[i] = GetText(strings[i],
			sdl.Color{red, green, blue, 255},
			newScreenWidth/2, newScreenHeight/2, renderer)
		defer textures[i].Destroy()
		h += dst[i].H * sizeMult / sizeDiv
		dst[i].X -= dst[i].W * sizeMult / sizeDiv
		dst[i].W *= 2 * sizeMult / sizeDiv
		dst[i].H *= 2 * sizeMult / sizeDiv
	}

	input.mono.a.down = false
	input.mono.b.down = false
	for i := 0; i < len(textures); i++ {
		dst[i].Y = newScreenHeight/2 - h
		h -= dst[i].H
	}

	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	for !input.mono.a.Down() && !input.mono.b.Down() && !quit {
		renderer.SetRenderTarget(nil)
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for i := 0; i < len(textures); i++ {
			PanicOnError(renderer.Copy(textures[i], src[i], dst[i]))
		}
		renderer.Present()
		input.Poll()
	}
}

func RunGame(menuChoice int, levelsCleared *int, score *int64) {
	lostLife = false
	lives := 3
	wonInARow := -2
	extraLives := 0
	extraLivesCounter := int64(25000)

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
			engine = stage.Load(stage.ID, false, *score, menuChoice)
			window.SetTitle("Score: " + strconv.Itoa(int(*score)) +
				" Lives: " + strconv.Itoa(lives))
		} else {
			*levelsCleared++
			wonInARow++
			if wonInARow == 3 {
				if lives-extraLives < 4 {
					wonInARow = 0
					lives++
				}
			}
			fmt.Printf("Won in a row counter: %d\n", wonInARow)
			engine = stage.Load(stage.ID+1, true,
				*score+ScoreMult(500), menuChoice)
		}
		fmt.Printf("Lives: %d\n", lives)
		if engine == nil {
			fmt.Println("Engine nil, game was won")
			break
		}
		PlayStage(engine, window, renderer, int32(lives))
		*score = engine.Score
		if engine.Input.exit.active {
			fmt.Println("Game was quit with exit key")
			break
		}
		for *score > extraLivesCounter &&
			extraLivesCounter*2 > extraLivesCounter {
			extraLivesCounter *= 2
			//extraLives++
			//lives++
		}
		fmt.Printf("Score: %d\n", *score)
	}
}

func PlayStage(engine *Engine, window *sdl.Window, renderer *sdl.Renderer,
	lives int32) {
	p1C, p2C := engine.GetPlayerSpriteID()

	quit = false
	lostLife = false
	score := int32(0)
	engine.Stage.scores.score, engine.Stage.scores.lives = engine.Score, lives
	for i := 0; i < 30 && !quit; i++ {
		engine.Stage.Render(p1C, p2C, renderer, false)
		engine.Input.Poll()
		if engine.Input.exit.active {
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
		if engine.Input.exit.active {
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
			p1C, p2C := engine.GetPlayerSpriteID()
			if engine.p1 != nil {
				engine.p1.entity.display = (i / 15 % 2) == 0
			}
			if engine.p2 != nil {
				engine.p2.entity.display = (i / 15 % 2) == 0
			}
			engine.Stage.Render(p1C, p2C, renderer, false)
			engine.Input.Poll()
			if engine.Input.exit.active {
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
	for v, _, _ := menu.Run(renderer, input); v != -1 &&
		!quit; v, _, _ = menu.Run(renderer, input) {
		ReadOptions("", input)
		menu.menuItems[0].SetNumber(int32(options.CharacterP1), renderer)
		menu.menuItems[1].SetNumber(int32(options.CharacterP2), renderer)
		menu.menuItems[2].SetNumber(int32(options.UseTap), renderer)
		menu.menuItems[3].SetNumber(int32(options.EdgeSlip), renderer)
		menu.menuItems[4].SetNumber(int32(options.BetterSlip), renderer)
		menu.menuItems[5].SetNumber(int32(options.ShowDivert), renderer)
	}
	if quit {
		return
	}
	options.CharacterP1 = uint8(menu.NVal(0))
	options.CharacterP2 = uint8(menu.NVal(1))
	options.UseTap = uint8(menu.NVal(2))
	options.useTap = options.UseTap != 0
	options.EdgeSlip = int(menu.NVal(3))
	options.BetterSlip = menu.NVal(4)
	options.ShowDivert = uint8(menu.NVal(5))
	options.showDivert = options.ShowDivert != 0
	redrawTextures = true
}

func GameOverMenu(levelsCleared int, score int64) (resume bool) {
	menuChoice := -1
	var scoreData *ScoreData
	menus[2].selectedElement = 0
	for !quit && menuChoice < 2 {
		menuChoice, _, _ = menus[2].Run(renderer, input)
		if menuChoice == 0 { // Set name
			name := GetName(defaultName, renderer, input)
			if name != "" {
				defaultName = name
				if scoreData == nil {
					scoreData = &ScoreData{score, name,
						levelsCleared, difficulty,
						time.Now()}
					highscores.Add(scoreData,
						menuChoice != 0, true)
				} else {
					scoreData.Name = name
				}
			}
		} else if menuChoice == 1 { // Highscores
			highscores.Display(difficulty, menuChoice != 0,
				renderer, input)
		} else if menuChoice == -1 {
			menuChoice = 4
		}
	}

	if quit {
		return false
	}

	resume = true
	if menuChoice == 2 { // Continue
		stage.ID--
	} else if menuChoice == 3 { // Restart
		stage.ID = -1
	} else if menuChoice == 4 { // Exit to menu
		resume = false
	} else {
		panic("Unknown menu option")
	}

	return
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogOnError(err error) bool {
	if err != nil {
		fmt.Println(err)
	}
	return err != nil
}
