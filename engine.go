package main

import "fmt"

const (
	PrecisionMax      int32 = 255 * 4
	EdgeSlipDefault   int   = 5
	BetterSlipDefault int32 = PrecisionMax * 13 / 40
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

func GetEngine(p1 *Player, p2 *Player, snakes []*Snake,
	stage *Stage) *Engine {
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
		engine.MoveSnakes()
	}
}

func (engine *Engine) MoveSnakes() {
	newPos := make([][2]int32, len(engine.snakes))
	for i := 0; i < len(engine.snakes); i++ {
		snake := engine.snakes[i]
		snake.moveTimer--
		if snake.moveTimer < 0 && !snake.shrinking {
			dir := snake.ai.Move(i, engine)
			engine.Stage.sprites.AlertSnakes(snake, snake.ai.CheckSignal())
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
			entity := engine.Stage.sprites.GetEntity(0, 0, snake.tail.spriteID)
			snake.body = append(snake.body, entity)
			entity.display = snake.tail.display
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
	var points uint64
	if engine.Stage.tiles.tiles[x][y] == Point {
		engine.Stage.tiles.tiles[x][y] = Empty
		engine.Stage.pointsLeft--
		points += uint64(10)
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p200 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points += uint64(200)
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p500 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points += uint64(500)
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p1000 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points += uint64(1000)
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == p2000 {
		engine.Stage.tiles.tiles[x][y] = Empty
		points += uint64(2000)
		modified = true
	} else if engine.Stage.tiles.tiles[x][y] == Powerup {
		engine.Stage.tiles.tiles[x][y] = Empty
		for i := 0; i < len(engine.snakes); i++ {
			points += uint64(75)
			snake := engine.snakes[i]
			if len(snake.body) > snake.minLength {
				for length := len(snake.body); length > 0; length /= 3 {
					points *= 2
				}
				if engine.snakes[i].shrinking {
					points /= 2
				}
				engine.snakes[i].shrinking = true
			}
		}
		modified = true
	}

	switch difficulty {
	case 0:
		player.score += 2 * points / 3
	case 1:
		player.score += 3 * points / 2
	case 2:
		player.score += points * 3
	case 3:
		player.score += points / 2
	case 4:
		player.score += points
	default:
		panic("Should not be reached")
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

	if controller.leftRight.Val() == 0 && controller.upDown.Val() == 0 && false {
		x, y := NewPos(e.x, e.y, e.dir)
		if !engine.LegalPos(x, y, false) {
			e.precision = 0
		}
		return
	}

	if /*e.precision > 15*PrecisionMax/16 ||*/
	controller.IsDirection(e.dir) && e.precision != 0 {
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
			//If a specific direction hasn't been found yet
			//Note:	Sides have a node either directly 1 tile in that direction
			//		of it, or if you can move sideways to reach such a tile.
			//		Distance is how many tiles you have to move sideways, with
			//		zero being the first type of situation described in the note
			dir := Direction(5)
			priority := -1
			for i := Up; i <= Left; i++ {
				if controller.IsDirection(i) {
					side := node.sides[i]
					if side != nil {
						distance := side.distance

						if distance != 0 {
							distance *= 2
							if e.dir == side.dirToPush {
								if e.precision >= options.BetterSlip {
									//If it's around a courner, apply BetterSlip
									distance--
								}
							} else {
								if e.precision > options.BetterSlip {
									//if it's around an inner courner, apply
									//BetterSlip (but since it would be - +
									// due to being on the far end of the
									// tile, we only have to add if it doesn't
									// apply)
									distance++
								}
							}
							if distance >= options.EdgeSlip {
								continue
							}
							distance /= 2
						}

						if priority == -1 ||
							(distance == 0 && priority != 0) {
							//If no node has been found yet, or this side leads
							//directly to a node, and the found node hasn't
							dir = side.dirToPush
							priority = side.distance
						} else if dir != 6 {
							//Else if there isn't a conflict so far
							if (priority == 0) == (distance == 0) {
								//If it has the same priority as the previous
								//found direction to move
								if dir != side.dirToPush {
									//If the found direction clashes with the
									//previously found direction, then mark
									//the conflict
									dir = 6
								}
							}
						}
					} else if priority == -1 {
						//Else if a valid direction hasn't been found yet
						if (e.dir+2)%4 == i {
							//If moving towards the center of the current tile
							priority = 1
							dir = i
						}
					} else if priority > 0 {
						//Else if it has the same priority as a previously found
						//direction, mark the conflict
						dir = 6
					}
				}
			}

			//If following the wall makes you move against the controller's
			//direction (courner where EdgeSlip ignores one direction)
			if controller.Dir(dir) == -1 {
				dir = 6
			}

			//Note:	Dir is 5 if no side was found,
			//		or if side.dirToPush is set to 5, which it is if the
			//		sideways travel-distance is equal in both directions
			if dir < 5 {
				//If a direction to move was found
				if e.dir == dir {
					//If it's already moving in that direction
					e.precision += step
				} else {
					//Else if it's within the specified distance of the edge
					//(Half a distance more due to being on the other side of
					// the current tile)
					e.precision -= step
					if e.precision < 0 {
						//If it crossed the border follow the wall
						e.precision = -e.precision
						e.dir = dir
					}
				}
			} else if dir == 6 || (dir == 5 && controller.Dir(e.dir) == -1) {
				//If there was a conflict, or if dir is 5 (see prev note) and
				//you are specifing a direction towards the middle of the tile
				e.precision -= step
				if e.precision < 0 {
					//If it crossed the border then stop
					e.precision = 0
				}
			}
		}
	}

	x, y := NewPos(e.x, e.y, e.dir)
	if engine.LegalPos(x, y, false) {
		if e.precision*2 > PrecisionMax {
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
