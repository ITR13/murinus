package main

import (
	"math/rand"
)

type Direction uint8

const (
	Up    Direction = iota
	Right Direction = iota
	Down  Direction = iota
	Left  Direction = iota
)

func NewPos(x, y int32, d Direction) (int32, int32) {
	d = d % 4
	if d == Up {
		return x, y - 1
	} else if d == Down {
		return x, y + 1
	} else if d == Left {
		return x - 1, y
	} else if d == Right {
		return x + 1, y
	}
	return x, y
}

//Only use with snakes
func (engine *Engine) LegalDir(x, y int32, d Direction) int {
	for i := 0; true; i++ {
		x, y = NewPos(x, y, d)
		if i > 3 || !engine.LegalPos(x, y, true) {
			return i
		}
	}
	return 0
}

type AI interface {
	Move(snakeID int, engine *Engine) Direction
	Reset()
}

type SimpleAI struct {
	lastDirection Direction
	turnedRight   bool
	ignore        int
}

func (simpleAI *SimpleAI) Move(snakeID int, engine *Engine) Direction {
	snake := engine.snakes[snakeID]
	X, Y := snake.head.x, snake.head.y
	options := make([]int, 4)
	legalOptions := 0
	for i := Up; i <= Left; i++ {
		options[i] = engine.LegalDir(X, Y, i)
		if options[i] > 0 {
			legalOptions++
		}
	}
	if legalOptions == 0 {
		return Up
	} else if legalOptions == 1 {
		for i := Up; i <= Left; i++ {
			if options[i] > 0 {
				simpleAI.lastDirection = i
				return i
			}
		}
	} else if legalOptions == 2 {
		fx, fy := NewPos(X, Y, simpleAI.lastDirection)
		if engine.LegalPos(fx, fy, true) {
			for i := Up; i <= Left; i++ {
				if options[i] > 0 && i != simpleAI.lastDirection {
					simpleAI.lastDirection = i
					return i
				}
			}
		}
	}

	if simpleAI.turnedRight {
		dir := (simpleAI.lastDirection + Left) % 4
		for i := 2; i > 0; i-- {
			if options[dir] >= i {
				simpleAI.turnedRight = false
				simpleAI.lastDirection = dir
				return dir
			}
			dir = (dir + 2) % 4
			if options[dir] >= i {
				simpleAI.turnedRight = true
				simpleAI.lastDirection = dir
				return dir
			}
			dir = (dir + 3) % 4
			if options[dir] >= i+1 {
				simpleAI.turnedRight = true
				simpleAI.lastDirection = dir
				return dir
			}
			dir = (dir + 3) % 4
		}
	} else {
		dir := (simpleAI.lastDirection + Right) % 4
		for i := 2; i > 0; i-- {
			if options[dir] >= i {
				simpleAI.turnedRight = true
				simpleAI.lastDirection = dir
				return dir
			}
			dir = (dir + 2) % 4
			if options[dir] >= i {
				simpleAI.turnedRight = false
				simpleAI.lastDirection = dir
				return dir
			}
			dir = (dir + 1) % 4
			if options[dir] >= i+1 {
				simpleAI.turnedRight = false
				simpleAI.lastDirection = dir
				return dir
			}
			dir = (dir + 1) % 4
		}
	}
	return simpleAI.lastDirection
}

func (simpleAI *SimpleAI) Reset() {
	simpleAI.lastDirection = Up
	simpleAI.turnedRight = false
	simpleAI.ignore = 0
}

type ApproximatedAI struct {
	divertTimer, divertTimerMax int
}

func (approx *ApproximatedAI) Move(snakeID int, engine *Engine) Direction {
	snake := engine.snakes[snakeID]
	X, Y := snake.head.x, snake.head.y
	options := make([]int, 4)
	legalOptions := 0
	for i := Up; i <= Left; i++ {
		options[i] = engine.LegalDir(X, Y, i)
		if options[i] > 0 {
			legalOptions++
		}
	}
	if legalOptions == 0 {
		return Up
	} else if legalOptions == 1 {
		for i := Up; i <= Left; i++ {
			if options[i] > 0 {
				return i
			}
		}
	}

	UpDownDir := Up
	LeftRightDir := Left
	dx := X - engine.p1.entity.x
	dy := Y - engine.p1.entity.y
	if dx < 0 {
		dx = -dx
		LeftRightDir = Right
	}
	if dy < 0 {
		dy = -dy
		UpDownDir = Down
	}
	approx.divertTimer--
	if approx.divertTimer < 0 {
		UpDownDir = (UpDownDir + 2) % 2
		LeftRightDir = (LeftRightDir + 2) % 2
		dx, dy = dy, dx
		approx.divertTimer = approx.divertTimerMax
	}

	if dx > dy {
		if options[LeftRightDir] > 0 {
			return LeftRightDir
		}
		dir := UpDownDir - LeftRightDir
		for i := UpDownDir; i != LeftRightDir; i = (i + dir) % 4 {
			if options[i] > 0 {
				return i
			}
		}
	} else {
		if options[UpDownDir] > 0 {
			return UpDownDir
		}
		dir := LeftRightDir - UpDownDir
		for i := LeftRightDir; i != UpDownDir; i = (i + dir) % 4 {
			if options[i] > 0 {
				return i
			}
		}
	}
	return Up
}

func (approx *ApproximatedAI) Reset() {
	approx.divertTimer = approx.divertTimerMax
}

type RandomAI struct {
	seed int64
	r    *rand.Rand
}

func (randAI *RandomAI) Move(snakeID int, engine *Engine) Direction {
	snake := engine.snakes[snakeID]
	X, Y := snake.head.x, snake.head.y
	options := make([]int, 4)
	legalOptions := 0
	for i := Up; i <= Left; i++ {
		options[i] = engine.LegalDir(X, Y, i)
		if options[i] > 0 {
			legalOptions++
		}
	}
	if legalOptions == 0 {
		return Up
	} else if legalOptions == 1 {
		for i := Up; i <= Left; i++ {
			if options[i] > 0 {
				return i
			}
		}
	}
	dir := randAI.r.Int() % legalOptions
	for i := Up; i <= Left; i++ {
		if options[i] > 0 {
			if dir == 0 {
				return i
			} else {
				dir--
			}
		}
	}
	return Up
}

func (randAI *RandomAI) Reset() {
	randAI.r = rand.New(rand.NewSource(randAI.seed))
}
