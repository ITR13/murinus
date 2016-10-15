package main

type Engine struct {
	players []*Player
	snakes  []*Snake
	Stage   *Stage
}

type Player struct {
	entity       *Entity
	timeout, max int
}

type Snake struct {
	head                    *Entity
	body                    []*Entity
	tail                    *Entity
	ai                      AI
	moveTimer, moveTimerMax int
	growTimer, growTimerMax int
	maxLength               int
}

func (engine *Engine) Advance() {
	if engine.players != nil {
		for i := 0; i < len(engine.players); i++ {
			panic("Do this")
		}
	}
	if engine.snakes != nil {
		for i := 0; i < len(engine.snakes); i++ {
			snake := engine.snakes[i]
			snake.moveTimer--
			if snake.moveTimer < 0 {
				snake.moveTimer = snake.moveTimerMax
				dir := snake.ai.Move(i, engine)
				x, y := NewPos(snake.head.x, snake.head.y, dir)
				snake.Move(x, y, engine)
			}
		}
	}
}

func (engine *Engine) LegalPos(x, y int32) bool {
	if engine.Stage.tiles.tiles[x][y] == Wall {
		return false
	}
	for i := 0; i < len(engine.snakes); i++ {
		snake := engine.snakes[i]
		if snake.head.Is(x, y) || snake.tail.Is(x, y) {
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
		if snake.growTimer <= 0 && len(snake.body) < snake.maxLength {
			snake.growTimer = snake.growTimerMax
			entity := Entity{snake.body[0].sprite, 0, 0}
			engine.Stage.sprites.entities = append(
				engine.Stage.sprites.entities, &entity)
			snake.body = append(snake.body, &entity)
			last++
		} else {
			snake.tail.x, snake.tail.y = snake.body[last].x, snake.body[last].y
		}
		for i := last; i > 0; i-- {
			snake.body[i].x, snake.body[i].y = snake.body[i-1].x, snake.body[i-1].y
		}
		snake.body[0].x, snake.body[0].y = snake.head.x, snake.head.y
	}

	if engine.LegalPos(x, y) {
		snake.head.x, snake.head.y = x, y
		snake.growTimer--
	}
}
