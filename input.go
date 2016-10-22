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
	a         *Key
	b         *Key
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

func GetInput(oneP bool) *Input {
	var input *Input
	if Arcade {
		panic("Make this")
	} else {
		if oneP {
			ipc := Direction(6)
			allInputs := make([]*Key, ipc)
			allInputs[Up] = GetKey(sdl.K_w, sdl.K_UP)
			allInputs[Right] = GetKey(sdl.K_d, sdl.K_RIGHT)
			allInputs[Down] = GetKey(sdl.K_s, sdl.K_DOWN)
			allInputs[Left] = GetKey(sdl.K_a, sdl.K_LEFT)
			allInputs[4] = GetKey(sdl.K_SPACE, sdl.K_RETURN)
			allInputs[5] = GetKey(sdl.K_LSHIFT, sdl.K_RSHIFT)
			p1 := Controller{
				&Axis{allInputs[Up], allInputs[Down]},
				&Axis{allInputs[Left], allInputs[Right]},
				allInputs[4], allInputs[5],
			}
			input = &Input{&p1, nil, allInputs}
		} else {
			ipc := Direction(6)
			allInputs := make([]*Key, ipc*2)
			allInputs[Up] = GetKey(sdl.K_w)
			allInputs[Right] = GetKey(sdl.K_d)
			allInputs[Down] = GetKey(sdl.K_s)
			allInputs[Left] = GetKey(sdl.K_a)
			allInputs[4] = GetKey(sdl.K_SPACE)
			allInputs[5] = GetKey(sdl.K_LSHIFT)
			allInputs[ipc+Up] = GetKey(sdl.K_UP)
			allInputs[ipc+Right] = GetKey(sdl.K_RIGHT)
			allInputs[ipc+Down] = GetKey(sdl.K_LEFT)
			allInputs[ipc+Left] = GetKey(sdl.K_DOWN)
			allInputs[ipc+4] = GetKey(sdl.K_RETURN)
			allInputs[ipc+5] = GetKey(sdl.K_RSHIFT)

			p1 := Controller{
				&Axis{allInputs[Up], allInputs[Down]},
				&Axis{allInputs[Left], allInputs[Right]},
				allInputs[4], allInputs[5],
			}
			p2 := Controller{
				&Axis{allInputs[ipc+Up], allInputs[ipc+Down]},
				&Axis{allInputs[ipc+Left], allInputs[ipc+Right]},
				allInputs[ipc+4], allInputs[ipc+5],
			}
			input = &Input{&p1, &p2, allInputs}
		}
	}
	return input
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
