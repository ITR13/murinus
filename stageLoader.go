package main

import "fmt"

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
	fmt.Printf("Loading stage %d, Tiles: %t\n", ID, loadTiles)

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
				8 * 4 * PrecisionMax / 127, score}
		} else if ID == 1 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"###******################" +
					"########**###############" +
					"#########**#0#0#0########" +
					"##########*#0#0#0########" +
					"##########***000000######" +
					"############*#0#0########" +
					"##########00***0000######" +
					"############0#*#0########" +
					"##########0000***00######" +
					"############0#0#0########" +
					"############0#0#0########" +
					"#########################")

			}
			p1 = &Player{stage.sprites.GetEntity(3, 3, Player1),
				8 * 4 * PrecisionMax / 127, score}
		} else if ID == 2 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#####**************3#####" +
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
				8 * 4 * PrecisionMax / 127, score}
			snakes = []*Snake{
				stage.sprites.GetSnake(stageWidth-6, stageHeight-6, 1,
					&SimpleAI{}, 20, 10000, 1, 1, 6)}
		}
	} else {
		ID -= 3

		STAGE := ID % 7
		ID /= 7

		speed := 11 - ID
		if speed-ID > 3 {
			speed -= ID
		} else {
			speed -= 2
		}
		if speed < 0 {
			speed = 0
		}
		pspeed := 8 * (int32(ID*3)/2 + 4) * PrecisionMax / 127
		if pspeed > PrecisionMax {
			pspeed = PrecisionMax
		}

		if STAGE == 0 {
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

			snakes = []*Snake{
				stage.sprites.GetSnake(1, stageHeight-2, 6, &SimpleAI{}, speed,
					100/(speed+2), 10*4, 2, 16),
				stage.sprites.GetSnake(stageWidth-2, 1, 6, &SimpleAI{}, speed,
					100/(speed+2), 10*4, 2, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
				pspeed, score}
		} else if STAGE == 1 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"######*****000*****######" +
					"######*###########*######" +
					"######*###########*######" +
					"######*###########*######" +
					"######*######0440#*######" +
					"######*######0##0#*######" +
					"######*************######" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################")
			}
			snakes = []*Snake{
				stage.sprites.GetSnake(stageWidth/2+1, stageHeight-5, 4,
					&ApproximatedAI{0, 3}, speed, 100/(speed+2), 10*4, 2, 4)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2,
				stageHeight/2-3, Player1), pspeed, score}
		} else if STAGE == 2 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"####*****************####" +
					"####*####*#####*####*####" +
					"####*####*#####*####*####" +
					"####*####*#####*####*####" +
					"####*****************####" +
					"####*####*#####*####*####" +
					"####*####*#####*####*####" +
					"####*####*#####*####*####" +
					"####*****************####" +
					"#########################" +
					"#########################" +
					"#########################")
			}
			snakes = []*Snake{
				stage.sprites.GetSnake(4, 4, 6, &ApproximatedAI{0, 3},
					speed, 100/(speed+2), 10*4, 2, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth-5,
				stageHeight-5, Player1), pspeed, score}
		} else if STAGE == 3 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"####*******###3000005####" +
					"####*#####*###0#0#0#0####" +
					"####*#####*###0000000####" +
					"####*#####*###0#0#0#0####" +
					"####*******%0%0000000####" +
					"####*#####*###0#0#0#0####" +
					"####*#####*###0000000####" +
					"####*#####*###0#0#0#0####" +
					"####*******###0000005####" +
					"#########################" +
					"#########################" +
					"#########################")
			}
			snakes = []*Snake{
				stage.sprites.GetSnake(stageWidth/2+2, stageHeight/2, 6,
					&ApproximatedAI{0, 3}, speed, 100/(speed+2), 10*4, 2, 16),
				stage.sprites.GetSnake(stageWidth/2-2, stageHeight/2, 6,
					&ApproximatedAI{0, 3}, speed, 100/(speed+2), 10*4, 2, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2,
				stageHeight/2, Player1), pspeed, score}
		} else if STAGE == 4 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#########################" +
					"#########################" +
					"#########################" +
					"########3*******3########" +
					"###########*#*###########" +
					"###########*0*###########" +
					"###########*#*###########" +
					"###########***###########" +
					"###########*#*###########" +
					"###########*#*###########" +
					"###########***###########" +
					"#########################" +
					"#########################" +
					"#########################")
			}
			snakes = []*Snake{
				stage.sprites.GetSnake(stageWidth/2, stageHeight/2+4, 6,
					&ApproximatedAI{0, 19},
					speed, 100/(speed+2), 10*4, 2, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2-1,
				Player1), pspeed, score}
		} else if STAGE == 5 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#0**********#***********#" +
					"#0#*#######*#*#*#4#4#4#*#" +
					"#0#*0003000***#*********#" +
					"#0#*#######*#*#########*#" +
					"#4#*********#***********#" +
					"#0#######*#*#*#*#########" +
					"#00000000*#*0*#*********#" +
					"#########*#*#*#*#######*#" +
					"#***********#***000000#*#" +
					"#*#########*#*#######0#*#" +
					"#000000000#***00000000#*#" +
					"#0#0#0#0#0#*#*#######0#*#" +
					"#0000000000*#00000000003#" +
					"#########################")
			}

			snakes = []*Snake{
				stage.sprites.GetSnake(1, stageHeight-2, 6, &ApproximatedAI{0, 3},
					speed, 100/(speed+2), 10*4, 2, 16),
				stage.sprites.GetSnake(stageWidth-2, 1, 6, &ApproximatedAI{0, 3},
					speed, 100/(speed+2), 10*4, 2, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
				pspeed, score}
		} else if STAGE == 6 {
			if loadTiles {
				tiles = ConvertStringToTiles("" +
					"#########################" +
					"#0**********#***********#" +
					"#0#*#######*#*#*#4#4#4#*#" +
					"#0#*0003000***#*********#" +
					"#0#*#######*#*#########*#" +
					"#4#*********#***********#" +
					"#0#######*#*#*#*#########" +
					"#00000000*#*0*#*********#" +
					"#########*#*#*#*#######*#" +
					"#***********#***000000#*#" +
					"#*#########*#*#######0#*#" +
					"#000000000#***00000000#*#" +
					"#0#0#0#0#0#*#*#######0#*#" +
					"#0000000000*#00000000003#" +
					"#########################")
			}

			snakes = []*Snake{
				stage.sprites.GetSnake(1, stageHeight-2, 6, &ApproximatedAI{0, 3},
					speed, 100/(speed+2), 10*4, 2, 16),
				stage.sprites.GetSnake(stageWidth-2, 1, 6, &ApproximatedAI{0, 3},
					speed, 100/(speed+2), 10*4, 2, 16),
				stage.sprites.GetSnake(stageWidth-2, stageHeight-2, 6, &SimpleAI{},
					speed, 100/(speed+2), 10*4, 2, 16),
				stage.sprites.GetSnake(1, 1, 6, &SimpleAI{},
					speed, 100/(speed+2), 10*4, 2, 16)}

			p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
				pspeed, score}
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
	fmt.Println("Finished loading stage")
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
