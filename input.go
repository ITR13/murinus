package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	Arcade bool = false
)

type Input struct {
	p1, p2 *Controller
	//mute *[]Key

	allInputs []*Key
}

type Controller struct {
	upDown    *Axis
	leftRight *Axis
	//	a         *Key
	//	b         *Key
}

type Axis struct {
	up, down *Key
}

type Key struct {
	keyCode []sdl.Keycode
	down    bool
}

func (axis *Axis) Val() int32 {
	val := int32(0)
	if axis.up.down {
		val -= 1
	}
	if axis.down.down {
		val += 1
	}
	return val
}

func GetKey(keyCode ...sdl.Keycode) *Key {
	return &Key{keyCode, false}
}

func GetInput() *Input {
	if Arcade {
		panic("Make this")
	} else {
		ipc := Direction(4)
		allInputs := make([]*Key, ipc*2)
		allInputs[Up] = GetKey(sdl.K_w)
		allInputs[Right] = GetKey(sdl.K_d)
		allInputs[Down] = GetKey(sdl.K_s)
		allInputs[Left] = GetKey(sdl.K_a)
		allInputs[ipc+Up] = GetKey(sdl.K_UP)
		allInputs[ipc+Right] = GetKey(sdl.K_RIGHT)
		allInputs[ipc+Down] = GetKey(sdl.K_LEFT)
		allInputs[ipc+Left] = GetKey(sdl.K_DOWN)

		p1 := Controller{
			&Axis{allInputs[Up], allInputs[Down]},
			&Axis{allInputs[Left], allInputs[Right]},
		}
		p2 := Controller{
			&Axis{allInputs[ipc+Up], allInputs[ipc+Down]},
			&Axis{allInputs[ipc+Left], allInputs[ipc+Right]},
		}

		return &Input{&p1, &p2, allInputs}
	}
}

func (input *Input) Poll() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) { //Add window resizing
		case *sdl.QuitEvent:
			quit = true
		case *sdl.KeyDownEvent:
			for i := 0; i < len(input.allInputs); i++ {
				key := input.allInputs[i]
				for k := 0; k < len(key.keyCode); k++ {
					if key.keyCode[k] == t.Keysym.Sym {
						key.down = true
						break
					}
				}
			}
		case *sdl.KeyUpEvent:
			for i := 0; i < len(input.allInputs); i++ {
				key := input.allInputs[i]
				for k := 0; k < len(key.keyCode); k++ {
					if key.keyCode[k] == t.Keysym.Sym {
						key.down = false
						break
					}
				}
			}
		}
	}
}

func (controller *Controller) IsDirection(dir Direction) bool {
	if dir == Up {
		return controller.upDown.Val() < 0
	} else if dir == Down {
		return controller.upDown.Val() > 0
	} else if dir == Left {
		return controller.leftRight.Val() < 0
	} else if dir == Right {
		return controller.leftRight.Val() > 0
	}
	return false
}

func (controller *Controller) Dir(dir Direction) int32 {
	if dir == Up {
		return -controller.upDown.Val()
	} else if dir == Down {
		return controller.upDown.Val()
	} else if dir == Left {
		return -controller.leftRight.Val()
	} else if dir == Right {
		return controller.leftRight.Val()
	}
	return 0
}
