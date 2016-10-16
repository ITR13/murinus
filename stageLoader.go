package main

const (
	stageWidth  int32 = 25
	stageHeight int32 = 15
)

func (stage *Stage) Load(ID int, loadTiles bool, score uint64) *Engine {
	stage.sprites.entities = make([]*Entity, 0)
	if loadTiles {
		stage.pointsLeft = 0
		tiles := make([][]Tile, stageWidth)
		for x := int32(0); x < stageWidth; x++ {
			tiles[x] = make([]Tile, stageHeight)
			for y := int32(0); y < stageHeight; y++ {
				if x == 0 || y == 0 || x == stageWidth-1 ||
					y == stageHeight-1 || (x%2 == 0 && y%2 == 0) {
					tiles[x][y] = Wall
				} else {
					tiles[x][y] = Point
					stage.pointsLeft++
				}
			}
		}
		stage.tiles.tiles = tiles
	}

	p1 := Player{stage.sprites.GetEntity(1, 1, Player1),
		0, 4, 32, score}

	engine := GetEngine(&p1, nil, stage,
		stage.sprites.GetSnake(1, stageHeight-2, 3, &SimpleAI{}, 0, 5, 10*2, 10*4, 100),
		stage.sprites.GetSnake(stageWidth-2, stageHeight-2, 3, &SimpleAI{}, 0, 5, 10*2, 10*4, 100))
	return engine
}
