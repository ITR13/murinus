package main

import (
	"fmt"
)

const (
	stageWidth  int32 = 25
	stageHeight int32 = 15
)

var tiles [][]Tile
var difficulty int

type PreStageData struct {
	stage          string
	px, py         int32
	difficultyData []*PreDifficultyData
}

type PreDifficultyData struct {
	playerSpeed int32
	snakes      []PreSnakeData
}

type PreSnakeData struct {
	x, y                 int32
	length               int
	ai                   AI
	moveTimerMax         int
	growTimerMax         int
	minLength, maxLength int
}

func GetPreStageData(stage string, px, py int32,
	snakes []PreSnakeData, speed [][]int32) *PreStageData {

	diffData := make([]*PreDifficultyData, len(speed))
	for i := 0; i < len(diffData); i++ {
		diffData[i] = &PreDifficultyData{
			0, make([]PreSnakeData, len(snakes)),
		}
		for j := 0; j < len(snakes); j++ {
			diffData[i].snakes[j] = snakes[j]
		}
	}

	for i := 0; i < len(speed); i++ {
		diffData[i].playerSpeed = speed[i][0]
		for j := 1; j < len(speed[i]); j++ {
			diffData[i].snakes[j-1].moveTimerMax = int(speed[i][j])
		}
	}

	return &PreStageData{stage, px, py, diffData}
}

func GetPreStageDatas() []*PreStageData {
	pSpeed := PrecisionMax / 4
	return []*PreStageData{
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########****############"+
			"#########*##*############"+
			"#########*##*############"+
			"#########***0***#########"+
			"############*##*#########"+
			"############*##*#########"+
			"############****#########"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################",
			stageWidth/2, stageHeight/2,
			nil, [][]int32{{pSpeed}, {pSpeed * 2}}),
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"###0*****################"+
			"########**###############"+
			"#########**#0#0#0########"+
			"##########*#0#0#0########"+
			"##########***000000######"+
			"############*#0#0########"+
			"##########00***0000######"+
			"############0#*#0########"+
			"##########0000***00######"+
			"############0#0#0########"+
			"############0#0#0########"+
			"#########################",
			3, 3,
			nil, [][]int32{{pSpeed}, {pSpeed * 2}}),
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#####0*************3#####"+
			"#####*#############*#####"+
			"#####*#############*#####"+
			"#####*#############*#####"+
			"#####4*************3#####"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################",
			5, 5,
			[]PreSnakeData{{stageWidth - 6, stageHeight - 6,
				1, &SimpleAI{}, 20, 1, 1, 6}},
			[][]int32{{pSpeed}, {pSpeed * 2}}),
		GetPreStageData(""+
			"#########################"+
			"#*******0000400000000003#"+
			"#*#####*#########0#####0#"+
			"#*#0000***************#0#"+
			"#*#0########*########*#0#"+
			"#*#0#***************#*#0#"+
			"#*#0#*######%######*#*#0#"+
			"#***#*%*****0*****%*#***#"+
			"#0#*#*######%######*#0#*#"+
			"#0#*#***************#0#*#"+
			"#0#*########*########0#*#"+
			"#0#***************0000#*#"+
			"#0#####0#########*#####*#"+
			"#3000000000040000*******#"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{1, stageHeight - 2, 6, &SimpleAI{},
				0, 10 * 4, 2, 16},
				{stageWidth - 2, 1, 6, &SimpleAI{},
					0, 10 * 4, 2, 16}},
			[][]int32{{pSpeed, 9, 9},
				{pSpeed * 3 / 2, 4, 4},
				{pSpeed * 2, 0, 0}}),
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################"+
			"######*****000*****######"+
			"######*###########*######"+
			"######*###########*######"+
			"######*###########*######"+
			"######*######0440#*######"+
			"######*######0##0#*######"+
			"######*************######"+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################",
			stageWidth/2, stageHeight/2-3,
			[]PreSnakeData{{stageWidth/2 + 1, stageHeight - 5, 4,
				&ApproximatedAI{0, 3}, 1, 10 * 4, 2, 4}},
			[][]int32{{pSpeed, 9},
				{pSpeed * 3 / 2, 4},
				{pSpeed * 3 / 2, 2}}),
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"####*****************####"+
			"####*####*#####*####*####"+
			"####*####*#####*####*####"+
			"####*####*#####*####*####"+
			"####*****************####"+
			"####*####*#####*####*####"+
			"####*####*#####*####*####"+
			"####*####*#####*####*####"+
			"####*****************####"+
			"#########################"+
			"#########################"+
			"#########################",
			stageWidth-5, stageHeight-5,
			[]PreSnakeData{{4, 4, 6, &ApproximatedAI{0, 3},
				1, 10 * 4, 2, 12}},
			[][]int32{{pSpeed, 9},
				{pSpeed * 3 / 2, 4},
				{pSpeed * 2, 1}}),
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"####*******###3000005####"+
			"####*#####*###0#0#0#0####"+
			"####*#####*###0000000####"+
			"####*#####*###0#0#0#0####"+
			"####*******%0%0000000####"+
			"####*#####*###0#0#0#0####"+
			"####*#####*###0000000####"+
			"####*#####*###0#0#0#0####"+
			"####*******###0000005####"+
			"#########################"+
			"#########################"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{stageWidth/2 + 2, stageHeight / 2, 6,
				&ApproximatedAI{0, 3}, 1, 10 * 4, 2, 10},
				{stageWidth/2 - 2, stageHeight / 2, 6,
					&ApproximatedAI{0, 3}, 1, 10 * 4, 2, 14}},
			[][]int32{{pSpeed, 9, 7},
				{pSpeed * 3 / 2, 4, 3},
				{pSpeed * 2, 1, 1}}),
		GetPreStageData(""+
			"#########################"+
			"#########################"+
			"#########################"+
			"#########################"+
			"########3*******3########"+
			"###########*#*###########"+
			"###########*0*###########"+
			"###########*#*###########"+
			"###########***###########"+
			"###########*#*###########"+
			"###########*#*###########"+
			"###########***###########"+
			"#########################"+
			"#########################"+
			"#########################",
			stageWidth/2, stageHeight/2-1,
			[]PreSnakeData{{stageWidth / 2, stageHeight/2 + 4, 6,
				&ApproximatedAI{0, 19}, 1, 10 * 4, 2, 6}},
			[][]int32{{pSpeed, 9},
				{pSpeed * 3 / 2, 6},
				{pSpeed * 2, 3}}),
		GetPreStageData(""+
			"#########################"+
			"#0**********#***********#"+
			"#0#*#######*#*#*#4#4#4#*#"+
			"#0#*0003000***#*********#"+
			"#0#*#######*#*#########*#"+
			"#4#*********#***********#"+
			"#0#######*#*#*#*#########"+
			"#00000000*#*0*#*********#"+
			"#########*#*#*#*#######*#"+
			"#***********#***000000#*#"+
			"#0#########*#*#######0#*#"+
			"#000000000#***00000000#*#"+
			"#0#0#0#0#0#0#0#######0#*#"+
			"#00000000000#00000000003#"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{1, stageHeight - 2, 6, &ApproximatedAI{0, 4},
				1, 10 * 4, 2, 16},
				{stageWidth - 2, 1, 6, &ApproximatedAI{0, 4},
					1, 10 * 4, 2, 16}},
			[][]int32{{pSpeed, 9, 9},
				{pSpeed * 3 / 2, 5, 5},
				{pSpeed * 2, 3, 3}}),
		GetPreStageData(""+
			"#########################"+
			"#0**********#00000000000#"+
			"#0#*#######*#0#0#6#6#6#0#"+
			"#0#*0003000*00#000000000#"+
			"#0#*#######*#0#########0#"+
			"#4#*********#00000000000#"+
			"#0#######*#*#0#0#########"+
			"#00000000*#***#000000000#"+
			"#########*#*#*#0#######0#"+
			"#***********#*00000000#0#"+
			"#0#########*#*#######0#0#"+
			"#000000000#***00000000#0#"+
			"#0#0#0#0#0#0#0#######0#0#"+
			"#00000000000#00000000003#"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{1, stageHeight - 2, 3, &ApproximatedAI{0, 3},
				1, 10 * 4, 1, 3},
				{stageWidth - 2, 1, 3, &ApproximatedAI{0, 3},
					1, 10 * 4, 1, 3},
				{stageWidth - 2, stageHeight - 2, 6, &SimpleAI{},
					1, 10 * 4, 2, 16},
				{1, 1, 6, &SimpleAI{},
					1, 10 * 4, 2, 16}},
			[][]int32{{pSpeed, 11, 11, 9, 9},
				{pSpeed * 3 / 2, 7, 7, 5, 5},
				{pSpeed * 2, 5, 5, 3, 3}}),
		GetPreStageData(""+
			"#########################"+
			"#***********************#"+
			"#*#0#0#0#0#0#0#0#0#0#0#0#"+
			"#***********************#"+
			"#0#0#0#0#0#0#0#0#0#0#0#*#"+
			"#***********************#"+
			"#*#0#0#0#0#0#0#0#0#0#0#0#"+
			"#***********000000000000#"+
			"#0#0#0#0#0#0#0#0#0#0#0#0#"+
			"#00000000000000000000000#"+
			"#0#0#5#0#0#0#0#0#0#5#0#0#"+
			"#00000000040004000000000#"+
			"#0#0#0#0#4#0#0#4#0#0#0#0#"+
			"#40000000000000000000004#"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{1, 1, 3, &ApproximatedAI{0, 3},
				1, 10 * 4, 1, 3},
				{stageWidth - 2, 1, 3, &ApproximatedAI{0, 3},
					1, 10 * 4, 1, 3},
				{stageWidth - 2, stageHeight - 2, 6, &SimpleAI{},
					1, 10 * 4, 2, 16},
				{1, stageHeight - 2, 6, &SimpleAI{},
					1, 10 * 4, 2, 16}},
			[][]int32{{pSpeed, 11, 11, 9, 9},
				{pSpeed * 3 / 2, 7, 7, 5, 5},
				{pSpeed * 2, 5, 5, 3, 3}}),
		GetPreStageData(""+
			"#########################"+
			"#50000000000%00000000005#"+
			"#0#########0#0#########0#"+
			"#0#*******************#0#"+
			"#0#*###*###*#*###*###*#0#"+
			"#0#*#***###*#*###***#*#0#"+
			"#0#*#*#*###*#*###*#*#*#0#"+
			"#0#***#***********#***#0#"+
			"#0#*###*#########*###*#0#"+
			"#0#*******************#0#"+
			"#0#*#####*#####*#####*#0#"+
			"#0#*******************#0#"+
			"#0#####################0#"+
			"#50000000000000000000005#"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{3, stageHeight - 4, 3, &ApproximatedAI{0, 3},
				1, 10 * 4, 1, 3},
				{stageWidth - 4, stageHeight - 4, 3, &ApproximatedAI{0, 3},
					1, 10 * 4, 1, 3},
				{stageWidth - 2, 1, 6, &ApproximatedAI{0, 5},
					1, 10 * 4, 2, 6},
				{1, 1, 6, &ApproximatedAI{0, 5},
					1, 10 * 4, 2, 10}},
			[][]int32{{pSpeed, 11, 11, 9, 9},
				{pSpeed * 3 / 2, 7, 7, 5, 5},
				{pSpeed * 2, 5, 5, 3, 3}}),
		GetPreStageData(""+
			"#########################"+
			"#****####0004000####****#"+
			"#*##**###0#####0###**##*#"+
			"#*###**##0#####0##**###*#"+
			"#*####**#0004000#**####*#"+
			"#*#####*#0#####0#*#####*#"+
			"#***###***********###***#"+
			"###***######*######***###"+
			"###*#***************#*###"+
			"###*########*########*###"+
			"###*##****##*##****##*###"+
			"###*##*##*******##*##*###"+
			"###*##*#####*#####*##*###"+
			"#00*******************00#"+
			"#########################",
			stageWidth/2, stageHeight/2,
			[]PreSnakeData{{stageWidth - 2, stageHeight - 2, 3, &ApproximatedAI{0, 3},
				1, 10 * 4, 1, 3},
				{stageWidth - 2, 1, 3, &ApproximatedAI{0, 3},
					1, 10 * 4, 1, 3},
				{1, stageHeight - 2, 6, &SimpleAI{},
					1, 10 * 4, 2, 16},
				{1, 1, 6, &SimpleAI{},
					1, 10 * 4, 2, 16}},
			[][]int32{{pSpeed, 11, 11, 9, 9},
				{pSpeed * 3 / 2, 7, 7, 5, 5},
				{pSpeed * 2, 5, 5, 3, 3}}),
	}
}

func (stage *Stage) Load(ID int, loadTiles bool, score uint64) *Engine {
	var p1, p2 *Player
	var snakes []*Snake
	stage.sprites.entities = make([]*Entity, 0)

	stage.ID = ID
	fmt.Printf("Loading stage %d, Tiles: %t\n", ID, loadTiles)

	levelIndex := ID
	diffIndex := 0
	if ID < 3 {
		diffIndex = difficulty / 2
	} else {
		if difficulty == 2 {
			if ID >= len(stage.levels) {
				stage.ID = 0
				ID = 0
				levelIndex = 0
			}
			diffIndex = len(stage.levels[levelIndex].difficultyData) - 1
		} else {
			if ID < 9 {

			} else if ID < 15 {
				diffIndex = 1
				levelIndex = ID - 6
			} else if ID < 16 {
				levelIndex = ID - 6
			} else if ID < 22 {
				levelIndex = ID - 13
				diffIndex = 2
			} else if ID < 23 {
				levelIndex = ID - 13
				diffIndex = 1
			} else if ID < 26 {
				levelIndex = ID - 13
			} else if ID < 27 {
				levelIndex = ID - 17
				diffIndex = 2
			} else if ID < 30 {
				levelIndex = ID - 17
				diffIndex = 1
			} else if ID < 33 {
				levelIndex = ID - 20
				diffIndex = 2
			} else {
				levelIndex = 0
				diffIndex = 1
				stage.ID = 0
				ID = 0
			}
		}
	}

	level := stage.levels[levelIndex]
	diffData := level.difficultyData[diffIndex]

	if loadTiles {
		ConvertStringToTiles(level.stage)
	}
	p1 = &Player{stage.sprites.GetEntity(level.px, level.py, Player1),
		diffData.playerSpeed, score}
	if diffData.snakes != nil {
		snakes = make([]*Snake, len(diffData.snakes))
		for i := 0; i < len(snakes); i++ {
			snake := diffData.snakes[i]
			snakes[i] = stage.sprites.GetSnake(snake.x, snake.y,
				snake.length, snake.ai, snake.moveTimerMax,
				snake.growTimerMax, snake.minLength, snake.maxLength)
			snakes[i].ai.Reset()
		}
	}
	fmt.Println("Exited set-up of stage")

	if loadTiles {
		fmt.Println("Calculating points left")
		stage.tiles.renderedOnce = false
		stage.pointsLeft = 0
		for x := int32(0); x < stageWidth; x++ {
			for y := int32(0); y < stageHeight; y++ {
				if tiles[x][y] == Point {
					stage.pointsLeft++
				}
			}
		}

		stage.tiles.tiles = tiles
		fmt.Printf("Replacing tiles\tPoints: %d\n", stage.pointsLeft)
	}
	fmt.Println("Getting engine")

	engine := GetEngine(p1, p2, snakes, stage)
	fmt.Println("Finished loading stage ", stage.ID)
	fmt.Printf("Level: %d\tDifficulty: %d\n", levelIndex, diffIndex)
	return engine
}

func ConvertStringToTiles(s string) {
	if tiles == nil {
		tiles = make([][]Tile, stageWidth)
		for x := int32(0); x < stageWidth; x++ {
			tiles[x] = make([]Tile, stageHeight)
		}
	}
	for x := int32(0); x < stageWidth; x++ {
		for y := int32(0); y < stageHeight; y++ {
			if x == 0 || y == 0 || x == stageWidth-1 ||
				y == stageHeight-1 {
				tiles[x][y] = Wall
			} else {
				tiles[x][y] = Point
			}
		}
	}

	for y := int32(0); y < stageHeight; y++ {
		for x := int32(0); x < stageWidth; x++ {
			c := s[y*stageWidth+x]
			if c == '#' {
				tiles[x][y] = Wall
			} else if c == '*' {
				tiles[x][y] = Point
			} else if c == '%' {
				tiles[x][y] = SnakeWall
			} else if c >= '0' && c <= '9' {
				tiles[x][y] = Tile(c - '0')
			}
		}
	}
}
