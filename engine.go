package main

import "fmt"

const (
	PrecisionMax int32 = 255 * 4
	EdgeSlip     int32 = 5
)

type Engine struct {
	p1, p2 *Player
	snakes []*Snake
	Stage  *Stage
	Input  *Input
	Graph  *Graph
}

type Player struct {
	entity *Entity
	//timeout, max int
	step  int32
	score uint64
}

type Snake struct {
	head                               *Entity
	body                               []*Entity
	tail                               *Entity
	ai                                 AI
	shrinking                          bool
	moveTimer, moveTimerMax            int
	speedUpTimer, speedUpTimerMax      int
	growTimer, growTimerMax            int
	normalLengthGrowTimer              int
	minLength, normalLength, maxLength int
}

func GetEngine(p1 *Player, p2 *Player, snakes []*Snake, stage *Stage) *Engine {
	fmt.Println("Getting input")
	input := GetInput(true)
	fmt.Println("Making graph")
	graph := stage.tiles.MakeGraph(false)
	fmt.Println("Returning engine")
	return &Engine{p1, p2, snakes, stage, input, graph}
}

func (engine *Engine) Advance() {
	if engine.p1 != nil {
		player := engine.p1
		player.Control(engine.Input.p1, engine)
		engine.CheckCollisions(player)
	}
	if engine.p2 != nil {
		player := engine.p2
		player.Control(engine.Input.p2, engine)
		engine.CheckCollisions(player)
	}

	if engine.snakes != nil {
		newPos := make([][2]int32, len(engine.snakes))
		for i := 0; i < len(engine.snakes); i++ {
			snake := engine.snakes[i]
			snake.moveTimer--
			if snake.moveTimer < 0 && !snake.shrinking {
				dir := snake.ai.Move(i, engine)
				x, y := NewPos(snake.head.x, snake.head.y, dir)
				newPos[i] = [2]int32{x, y}
			}
		}
		for i := 0; i < len(engine.snakes); i++ {
			snake := engine.snakes[i]
			if snake.moveTimer < 0 {
				snake.moveTimer = snake.moveTimerMax
				snake.Move(newPos[i][0], newPos[i][1], engine)
			}
		}
	}
}

func (engine *Engine) LegalPos(x, y int32, isSnake bool) bool {
	if x < 0 || y < 0 || x >= engine.Stage.tiles.w || y >= engine.Stage.tiles.h {
		return false
	}
	if engine.Stage.tiles.tiles[x][y] == Wall ||
		(isSnake && engine.Stage.tiles.tiles[x][y] == SnakeWall) {
		return false
	}
	for i := 0; i < len(engine.snakes); i++ {
		snake := engine.snakes[i]
		if snake.head.Is(x, y) || (!isSnake && snake.tail.Is(x, y)) {
			return false
		}
		for k := 0; k < len(snake.body); k++ {
			if snake.body[k].Is(x, y) {
				return false
			}
		}
	}
	return true
}

func (entity *Entity) Is(x, y int32) bool {
	return x == entity.x && y == entity.y
}

func (snake *Snake) Move(x, y int32, engine *Engine) {
	if snake.body == nil || len(snake.body) == 0 {
		snake.tail.x, snake.tail.y = snake.head.x, snake.head.y
	} else {
		last := len(snake.body) - 1
		grow := false
		if !snake.shrinking {
			if len(snake.body) < snake.normalLength {
				if snake.normalLengthGrowTimer <= 0 {
					snake.normalLengthGrowTimer = 4
					grow = true
				}
			} else if snake.growTimer <= 0 && len(snake.body) < snake.maxLength {
				snake.growTimer = snake.growTimerMax
				snake.normalLength++
				grow = true
			}
		}
		if grow {
			entity := engine.Stage.sprites.GetEntity(0, 0, SnakeBody)
			snake.body = append(snake.body, entity)
			last++
		} else {
			snake.tail.x, snake.tail.y = snake.body[last].x, snake.body[last].y
		}

		for i := last; i > 0; i-- {
			snake.body[i].x, snake.body[i].y = snake.body[i-1].x, snake.body[i-1].y
		}
		if snake.shrinking {
			snake.body[0].display = false
			snake.body = snake.body[1:]
			if len(snake.body) <= snake.minLength {
				snake.shrinking = false
				snake.normalLengthGrowTimer = 1
			}
		} else {
			snake.body[0].x, snake.body[0].y = snake.head.x, snake.head.y
		}
	}

	if !snake.shrinking && engine.LegalPos(x, y, true) {
		snake.head.x, snake.head.y = x, y
		snake.growTimer--
		snake.normalLengthGrowTimer--
		snake.speedUpTimer--
		if snake.speedUpTimer <= 0 && snake.moveTimerMax > 0 {
			if snake.moveTimerMax > 1 {
				snake.speedUpTimerMax = (snake.moveTimerMax*300)/
					(snake.moveTimerMax-1) + snake.speedUpTimerMax/2
			}
			snake.speedUpTimer = snake.speedUpTimerMax
			snake.moveTimerMax--
			fmt.Printf("Sped up to %d!\tNext speed up: %d\n", snake.moveTimerMax, snake.speedUpTimer)
		}
	}
}

func (engine *Engine) CheckCollisions(player *Player) {
	x, y := player.entity.x, player.entity.y
	modified := false
	if engine.Stage.tiles.tiles[x][y] == Point {
		engine.Stage.tiles.tiles[x][y] = Empty
		engine.Stage.pointsLeft--
		player.score += 10
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p200 {
		engine.Stage.tiles.tiles[x][y] = Empty
		player.score += 200
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p500 {
		engine.Stage.tiles.tiles[x][y] = Empty
		player.score += 500
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p1000 {
		engine.Stage.tiles.tiles[x][y] = Empty
		player.score += 1000
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p2000 {
		engine.Stage.tiles.tiles[x][y] = Empty
		player.score += 2000
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == Powerup {
		engine.Stage.tiles.tiles[x][y] = Empty
		for i := 0; i < len(engine.snakes); i++ {
			points := uint64(75)
			for length := len(engine.snakes[i].body); length > 0; length /= 3 {
				points *= 2
			}
			if engine.snakes[i].shrinking {
				points /= 2
			}
			player.score += points
			engine.snakes[i].shrinking = true
		}
		modified = true
	}

	if modified {
		engine.Stage.tiles.renderedOnce = false
	}

	for i := 0; i < len(engine.snakes); i++ {
		if engine.snakes[i].head.Is(x, y) {
			lostLife = true
			break
		}
	}
}

func (player *Player) Control(controller *Controller, engine *Engine) {
	e := player.entity
	step := player.step
	if controller.b.down {
		step = (step * 2) / 3
	}

	if e.precision > 15*PrecisionMax/16 || controller.IsDirection(e.dir) {
		e.precision += controller.Dir(e.dir) * step
	} else {
		checkDir := false
		perpMove := false
		if e.dir == Up || e.dir == Down {
			val := controller.leftRight.Val()
			if val != 0 {
				perpMove = true
				if !engine.LegalPos(e.x+val, e.y, false) {
					val = 0
				}
			}
			if val != 0 {
				e.precision -= step
				if e.precision < 0 {
					e.precision = -e.precision
					if val < 0 {
						e.dir = Left
					} else {
						e.dir = Right
					}
				}
			} else {
				checkDir = true
			}
		} else if e.dir == Right || e.dir == Left {
			val := controller.upDown.Val()
			if val != 0 {
				perpMove = true
				if !engine.LegalPos(e.x, e.y+val, false) {
					val = 0
				}
			}
			if val != 0 {
				e.precision -= step
				if e.precision < 0 {
					e.precision = -e.precision
					if val < 0 {
						e.dir = Up
					} else {
						e.dir = Down
					}
				}
			} else {
				checkDir = true
			}
		}
		if checkDir {
			val := controller.Dir(e.dir)
			if val != 0 {
				e.precision += val * step
				if e.precision < 0 {
					dir := (e.dir + 2) % 4
					x, y := NewPos(e.x, e.y, dir)
					if engine.LegalPos(x, y, false) {
						e.precision = -e.precision
						e.dir = dir
					} else {
						e.precision = 0
					}
				}
			} else if perpMove {
				edge := engine.Graph.edge[e.x][e.y]
				if edge != nil && edge.distance > 0 &&
					edge.distance < EdgeSlip {
					if e.dir != edge.dir && edge.me != nil {
						if edge.distance < (EdgeSlip-1)/2 {
							e.precision -= step
							if e.precision < 0 {
								e.precision = -e.precision
								e.dir = edge.dir
							}
						}
					} else {
						e.precision += step
					}
				}
			}
		}
	}
	x, y := NewPos(e.x, e.y, e.dir)
	if engine.LegalPos(x, y, false) {
		if e.precision > PrecisionMax/2 {
			e.precision = PrecisionMax - e.precision
			e.x, e.y = x, y
			e.dir = (e.dir + 2) % 4
		}
	} else {
		e.precision = 0
	}
	if e.precision < 0 {
		panic("Should not reach this point")
	}
}
