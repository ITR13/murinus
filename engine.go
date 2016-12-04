package main

import "fmt"

const (
	PrecisionMax int32 = 255 * 4
	EdgeSlip     int   = 7
	BetterSlip   int32 = PrecisionMax * 13 / 40
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
	step   int32
	score  uint64
}

type Snake struct {
	head                               *Entity
	body                               []*Entity
	tail                               *Entity
	ai                                 AI
	shrinking                          bool
	moveTimer, moveTimerMax            int
	growTimer, growTimerMax            int
	normalLengthGrowTimer              int
	minLength, normalLength, maxLength int
}

func GetEngine(p1 *Player, p2 *Player, snakes []*Snake, stage *Stage) *Engine {
	fmt.Println("Making graph")
	graph := stage.tiles.MakeGraph(false)
	fmt.Println("Returning engine")
	return &Engine{p1, p2, snakes, stage, stage.input, graph}
}

func (engine *Engine) Advance() {
	if engine.p1 != nil {
		player := engine.p1
		if engine.p2 != nil {
			player.Control(engine.Input.p1, engine)
		} else {
			player.Control(engine.Input.mono, engine)
		}
		engine.CheckCollisions(player)
	}
	if engine.p2 != nil {
		player := engine.p2
		if engine.p1 != nil {
			player.Control(engine.Input.p2, engine)
		} else {
			player.Control(engine.Input.mono, engine)
		}
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
	}
}

func (engine *Engine) CheckCollisions(player *Player) {
	x, y := player.entity.x, player.entity.y
	modified := false
	if engine.Stage.tiles.tiles[x][y] == Point {
		engine.Stage.tiles.tiles[x][y] = Empty
		engine.Stage.pointsLeft--
		points := uint64(10)
		if difficulty == 0 {
			player.score += points / 2
		} else if difficulty == 1 {
			player.score += points
		} else if difficulty == 2 {
			player.score += points * 5
		}
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p200 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points := uint64(200)
		if difficulty == 0 {
			player.score += points / 2
		} else if difficulty == 1 {
			player.score += points
		} else if difficulty == 2 {
			player.score += points * 5
		}
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p500 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points := uint64(500)
		if difficulty == 0 {
			player.score += points / 2
		} else if difficulty == 1 {
			player.score += points
		} else if difficulty == 2 {
			player.score += points * 5
		}
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p1000 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points := uint64(1000)
		if difficulty == 0 {
			player.score += points / 2
		} else if difficulty == 1 {
			player.score += points
		} else if difficulty == 2 {
			player.score += points * 5
		}
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p2000 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points := uint64(2000)
		if difficulty == 0 {
			player.score += points / 2
		} else if difficulty == 1 {
			player.score += points
		} else if difficulty == 2 {
			player.score += points * 5
		}
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == Powerup {
		engine.Stage.tiles.tiles[x][y] = Empty
		for i := 0; i < len(engine.snakes); i++ {
			points := uint64(75)
			if len(engine.snakes[i].body) > engine.snakes[i].minLength {
				for length := len(engine.snakes[i].body); length > 0; length /= 3 {
					points *= 2
				}
				if engine.snakes[i].shrinking {
					points /= 2
				}
				engine.snakes[i].shrinking = true
			}
			if difficulty == 0 {
				player.score += points / 2
			} else if difficulty == 1 {
				player.score += points
			} else if difficulty == 2 {
				player.score += points * 5
			}
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
		step = (step * 3) / 5
	}

	if controller.leftRight.Val() == 0 && controller.upDown.Val() == 0 {
		return
	}

	if e.precision > 15*PrecisionMax/16 ||
		(controller.IsDirection(e.dir) && e.precision != 0) {
		e.precision += step
	} else {
		node := engine.Graph.nodes[e.x][e.y]
		if node == nil {
			panic("You are inside a wall")
		}
		useGraph := true
		if e.dir == Up || e.dir == Down {
			val := controller.leftRight.Val()
			if val != 0 {
				if engine.Graph.nodes[e.x+val][e.y] != nil {
					useGraph = false
					e.precision -= step
					if e.precision < 0 {
						e.precision = -e.precision
						e.dir = Direction(2 - val)
					}
				}
			}
		} else {
			val := controller.upDown.Val()
			if val != 0 {
				if engine.Graph.nodes[e.x][e.y+val] != nil {
					useGraph = false
					e.precision -= step
					if e.precision < 0 {
						e.precision = -e.precision
						e.dir = Direction((3 - val) % 4)
					}
				}
			}
		}

		if useGraph {
			dir := Direction(5)
			priority := -1
			for i := Up; i <= Left; i++ {
				if controller.IsDirection(i) {
					side := node.sides[i]
					if side != nil {
						if priority == -1 || (side.distance == 0 && priority != 0) {
							dir = side.dirToPush
							priority = side.distance
						} else if dir != 6 {
							if (priority == 0) == (side.distance == 0) {
								if dir != side.dirToPush {
									dir = 6
								}
							}
						}
					} else if priority == -1 {
						if (e.dir+2)%4 == i {
							priority = 1
							dir = i
						}
					} else if priority > 0 {
						dir = 6
					}
				}
			}
			if dir < 5 {
				if e.dir == dir {
					if priority*2 < EdgeSlip {
						e.precision += step
					}
				} else {
					e.precision -= step
					if e.precision < 0 {
						if priority*2+1 < EdgeSlip {
							e.precision = -e.precision
							e.dir = dir
						} else {
							e.precision = 0
						}
					}
				}
			} else if dir == 6 || (dir == 5 && controller.Dir(e.dir) == -1) {
				e.precision -= step
				if e.precision < 0 {
					e.precision = 0
				}
			} else if dir == 5 {
				if priority*2 < EdgeSlip {
					if e.precision >= BetterSlip {
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
