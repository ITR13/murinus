/*
   This file is part of Murinus.

   Murinus is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Murinus is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Murinus.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Input struct {
	mono   *Controller
	p1, p2 *Controller
	exit   *DurationKey

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

type DurationKey struct {
	key      *Key
	active   bool
	timeHeld int
}

type Key struct {
	KeyCode []sdl.Keycode
	tapped  bool
	down    bool
}

func (axis *Axis) Val() int32 {
	val := int32(0)
	if axis.up.Down() {
		val -= 1
	}
	if axis.down.Down() {
		val += 1
	}
	return val
}

func (key *Key) Down() bool {
	return key.down || key.tapped
}

func (key *Key) Clear() {
	key.down = false
	key.tapped = false
}

func GetKey(keyCode ...sdl.Keycode) *Key {
	return &Key{keyCode, false, false}
}

func GetInput() *Input {
	ipc := Direction(6)
	allInputs := make([]*Key, ipc*3+1)
	if Arcade {
		allInputs[Up] = GetKey(sdl.K_w, sdl.K_UP)
		allInputs[Right] = GetKey(sdl.K_d, sdl.K_RIGHT)
		allInputs[Down] = GetKey(sdl.K_s, sdl.K_DOWN)
		allInputs[Left] = GetKey(sdl.K_a, sdl.K_LEFT)
		allInputs[4] = GetKey(sdl.K_SPACE, sdl.K_RETURN, sdl.K_g, sdl.K_l)
		allInputs[5] = GetKey(sdl.K_f, sdl.K_r, sdl.K_k, sdl.K_o)

		allInputs[ipc+Up] = GetKey(sdl.K_w)
		allInputs[ipc+Right] = GetKey(sdl.K_d)
		allInputs[ipc+Down] = GetKey(sdl.K_s)
		allInputs[ipc+Left] = GetKey(sdl.K_a)
		allInputs[ipc+4] = GetKey(sdl.K_g, sdl.K_SPACE)
		allInputs[ipc+5] = GetKey(sdl.K_f, sdl.K_r)

		allInputs[ipc*2+Up] = GetKey(sdl.K_UP)
		allInputs[ipc*2+Right] = GetKey(sdl.K_RIGHT)
		allInputs[ipc*2+Down] = GetKey(sdl.K_DOWN)
		allInputs[ipc*2+Left] = GetKey(sdl.K_LEFT)
		allInputs[ipc*2+4] = GetKey(sdl.K_l, sdl.K_RETURN)
		allInputs[ipc*2+5] = GetKey(sdl.K_k, sdl.K_o)
		allInputs[ipc*3] = GetKey(sdl.K_t, sdl.K_p)
	} else {
		allInputs[Up] = GetKey(sdl.K_w, sdl.K_UP)
		allInputs[Right] = GetKey(sdl.K_d, sdl.K_RIGHT)
		allInputs[Down] = GetKey(sdl.K_s, sdl.K_DOWN)
		allInputs[Left] = GetKey(sdl.K_a, sdl.K_LEFT)
		allInputs[4] = GetKey(sdl.K_SPACE, sdl.K_RETURN)
		allInputs[5] = GetKey(sdl.K_LSHIFT, sdl.K_RSHIFT, sdl.K_BACKSPACE)

		allInputs[ipc+Up] = GetKey(sdl.K_w)
		allInputs[ipc+Right] = GetKey(sdl.K_d)
		allInputs[ipc+Down] = GetKey(sdl.K_s)
		allInputs[ipc+Left] = GetKey(sdl.K_a)
		allInputs[ipc+4] = GetKey(sdl.K_SPACE)
		allInputs[ipc+5] = GetKey(sdl.K_LSHIFT)

		allInputs[ipc*2+Up] = GetKey(sdl.K_UP)
		allInputs[ipc*2+Right] = GetKey(sdl.K_RIGHT)
		allInputs[ipc*2+Down] = GetKey(sdl.K_DOWN)
		allInputs[ipc*2+Left] = GetKey(sdl.K_LEFT)
		allInputs[ipc*2+4] = GetKey(sdl.K_RETURN)
		allInputs[ipc*2+5] = GetKey(sdl.K_RSHIFT, sdl.K_BACKSPACE)
		allInputs[ipc*3] = GetKey(sdl.K_ESCAPE)
	}
	mono := Controller{
		&Axis{allInputs[Up], allInputs[Down]},
		&Axis{allInputs[Left], allInputs[Right]},
		allInputs[4], allInputs[5],
	}
	p1 := Controller{
		&Axis{allInputs[ipc+Up], allInputs[ipc+Down]},
		&Axis{allInputs[ipc+Left], allInputs[ipc+Right]},
		allInputs[ipc+4], allInputs[ipc+5],
	}
	p2 := Controller{
		&Axis{allInputs[ipc*2+Up], allInputs[ipc*2+Down]},
		&Axis{allInputs[ipc*2+Left], allInputs[ipc*2+Right]},
		allInputs[ipc*2+4], allInputs[ipc*2+5],
	}
	escape := DurationKey{allInputs[ipc*3], false, 9}
	return &Input{&mono, &p1, &p2, &escape, allInputs}
}

var noKeysTouched int
var newScreenWidth, newScreenHeight int32
var redrawTextures bool

func (input *Input) Poll() {
	sdl.Delay(17)

	noKeysTouched++
	if Arcade && noKeysTouched > 60*60*15 {
		quit = true
	}

	for i := 0; i < len(input.allInputs); i++ {
		input.allInputs[i].tapped = false
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			quit = true
		case *sdl.KeyDownEvent:
			if t.Repeat == 0 {
				for i := 0; i < len(input.allInputs); i++ {
					key := input.allInputs[i]
					for k := 0; k < len(key.KeyCode); k++ {
						if key.KeyCode[k] == t.Keysym.Sym {
							key.down = true
							key.tapped = true
							noKeysTouched = 0
							break
						}
					}
				}
			}
		case *sdl.KeyUpEvent:
			if t.Repeat == 0 {
				for i := 0; i < len(input.allInputs); i++ {
					key := input.allInputs[i]
					for k := 0; k < len(key.KeyCode); k++ {
						if key.KeyCode[k] == t.Keysym.Sym {
							key.down = false
							break
						}
					}
				}
			}
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
				newScreenWidth, newScreenHeight = t.Data1, t.Data2
			}

		case *sdl.RenderEvent:
			redrawTextures = true
		}
	}

	if input.exit.key.Down() {
		input.exit.timeHeld++
		if input.exit.timeHeld > timeExitHasToBeHeldToExit {
			input.exit.active = true
			noKeysTouched = 0
			if input.exit.timeHeld > timeExitHasToBeHeldToQuit {
				quit = true
			}
		}
	} else {
		input.exit.timeHeld = 0
		input.exit.active = false
	}
}

func (input *Input) Clear() {
	for i := 0; i < len(input.allInputs); i++ {
		input.allInputs[i].Clear()
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

func (key *Key) Stepper(start, mod int) func() bool {
	c := 0
	return func() bool {
		if key.Down() {
			c++
			return c == 1 || (c > start && c%mod == 0)
		} else {
			c = 0
		}
		return false
	}
}

func (axis *Axis) Stepper(start, mod int) func() int32 {
	c := 0
	return func() int32 {
		val := axis.Val()
		if val == 0 {
			c = 0
		} else {
			c++
			if c == 1 || (c > start && c%mod == 0) {
				return val
			}
		}
		return 0
	}
}
