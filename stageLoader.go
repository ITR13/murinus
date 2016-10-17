package main

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
	if true {
		if loadTiles {
			tiles = ConvertStringToTiles("" +
				"#########################" +
				"#3*********************3#" +
				"#*#####*#########*#####*#" +
				"#*#*******************#*#" +
				"#*#*########*########*#*#" +
				"#*#*#***************#*#*#" +
				"#*#*#*######%######*#*#*#" +
				"#***#*%*****0*****%*#***#" +
				"#*#*#*######%######*#*#*#" +
				"#*#*#***************#*#*#" +
				"#*#*########*########*#*#" +
				"#*#*******************#*#" +
				"#*#####*#########*#####*#" +
				"#3*********************3#" +
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
			stage.sprites.GetSnake(1, stageHeight-2, 6, &SimpleAI{}, speed, 100/(speed+2), 10*4, 16),
			stage.sprites.GetSnake(stageWidth-2, 1, 6, &SimpleAI{}, speed, 100/(speed+2), 10*4, 16)}

		p1 = &Player{stage.sprites.GetEntity(stageWidth/2, stageHeight/2, Player1),
			0, 4, pspeed, score}
	} else if ID == 1 {
		if loadTiles {
			for i := int32(2); i < 6; i++ {
				tiles[i][2] = Wall
				tiles[2][i] = Wall
				tiles[i][stageHeight-1-2] = Wall
				tiles[2][stageHeight-1-i] = Wall
				tiles[stageWidth-1-i][2] = Wall
				tiles[stageWidth-1-2][i] = Wall
				tiles[stageWidth-1-i][stageHeight-1-2] = Wall
				tiles[stageWidth-1-2][stageHeight-1-i] = Wall
				tiles[i][stageHeight/2] = Wall
				tiles[stageWidth-i-1][stageHeight/2] = Wall

				tiles[i+5][2] = Wall
				tiles[stageWidth-i-6][2] = Wall
				tiles[i+5][stageHeight-1-2] = Wall
				tiles[stageWidth-i-6][stageHeight-1-2] = Wall

				tiles[stageWidth/2][i] = Wall
				tiles[stageWidth/2][stageHeight-i-1] = Wall
			}
		}

		p1 = &Player{stage.sprites.GetEntity(1, 1, Player1),
			0, 4, 32, score}
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
