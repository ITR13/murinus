package main

import "math/rand"

const (
	stageWidth  int32 = 25
	stageHeight int32 = 15
)

func (stage *Stage) Load(ID int, loadTiles bool, score uint64) *Engine {
	var p1, p2 *Player
	var tiles [][]Tile
	var snakes []*Snake

	stage.sprites.entities = make([]*Entity, 0)
	stage.ID = ID

	if ID < 3 {
		if ID == 0 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########****############" +
					"#########*##*############" +
					"#########*##*############" +
					"#########***0***#########" +
					"############*##*#########" +
					"############*##*#########" +
					"############****#########" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################")

			}
			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
				0, 4, 8 * 4 * PrecisionMax / 127, score}
		} else if ID == 1 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"###******################" +
					"########**###############" +
					"#########**##############" +
					"##########*##############" +
					"##########***0000########" +
					"############*#0#0########" +
					"############***00########" +
					"############0#*#0########" +
					"############00***########" +
					"#########################" +
					"#########################" +
					"#########################")

			}
			p1 = &Player{stage.sprites.GetEntity(3, 3, Player1),
				0, 4, 8 * 4 * PrecisionMax / 127, score}
		} else if ID == 2 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#####**************4#####" +
					"#####*#############*#####" +
					"#####*#############*#####" +
					"#####*#############*#####" +
					"#####4*************3#####" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################")

			}
			p1 = &Player{stage.sprites.GetEntity(5, 5, Player1),
				0, 4, 8 * 4 * PrecisionMax / 127, score}
			snakes = []*Snake{
				stage.sprites.GetSnake(stageWidth-6, stageHeight-6, 6,
					&SimpleAI{}, 20, 10000, 10000, 6)}
		}
	} else {
		ID -= 3
		if ID%2 == 0 {
			ID /= 2
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#*******0000400000000003#" +
					"#*#####*#########0#####0#" +
					"#*#0000***************#0#" +
					"#*#0########*########*#0#" +
					"#*#0#***************#*#0#" +
					"#*#0#*######%######*#*#0#" +
					"#***#*%*****0*****%*#***#" +
					"#0#*#*######%######*#0#*#" +
					"#0#*#***************#0#*#" +
					"#0#*########*########0#*#" +
					"#0#***************0000#*#" +
					"#0#####0#########*#####*#" +
					"#3000000000040000*******#" +
					"#########################")

			}

			speed := 18 - ID*5
			if speed < 0 {
				speed = 0
			}

			pspeed := 8 * int32(ID+4) * PrecisionMax / 127
			if pspeed > PrecisionMax/2 {
				pspeed = PrecisionMax / 2
			}

			snakes = []*Snake{
				stage.sprites.GetSnake(1, stageHeight-2, 6, &SimpleAI{}, speed,
					100/(speed+2), 10*4, 16),
				stage.sprites.GetSnake(stageWidth-2, 1, 6, &SimpleAI{}, speed,
					100/(speed+2), 10*4, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
				0, 4, pspeed, score}
		} else {
			ID /= 2
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#***********#***********#" +
					"#*#*#######*#*#*#*#*#*#*#" +
					"#*#***********#*********#" +
					"#*#*#######*#*#########*#" +
					"#*#*********#***********#" +
					"#*#######*#*#*#*#########" +
					"#*********#*0*#*********#" +
					"#########*#*#*#*#######*#" +
					"#***********#*********#*#" +
					"#*#########*#*#######*#*#" +
					"#*********#***********#*#" +
					"#*#*#*#*#*#*#*#######*#*#" +
					"#***********#***********#" +
					"#########################")
			}

			speed := 18 - ID*5
			if speed < 0 {
				speed = 0
			}

			pspeed := 8 * int32(ID+4) * PrecisionMax / 127
			if pspeed > PrecisionMax/2 {
				pspeed = PrecisionMax / 2
			}

			snakes = []*Snake{
				stage.sprites.GetSnake(1, stageHeight-2, 6, &RandomAI{
					rand.New(rand.NewSource(0))}, speed, 100/(speed+2), 10*4, 16),
				stage.sprites.GetSnake(stageWidth-2, 1, 6, &RandomAI{
					rand.New(rand.NewSource(0))}, speed, 100/(speed+2), 10*4, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
				0, 4, pspeed, score}
		}
	}

	if loadTiles {
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
	}
	engine := GetEngine(p1, p2, snakes, stage)

	return engine
}

func ConvertStringToTiles(s string) [][]Tile {
	tiles := make([][]Tile, stageWidth)
	for x := int32(0); x < stageWidth; x++ {
		tiles[x] = make([]Tile, stageHeight)
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
	return tiles
}
