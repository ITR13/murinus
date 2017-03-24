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
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	MenuXOffset   int32 = 64 * 5 * sizeMult / sizeDiv
	MenuArrowSize int32 = 1280 // / 2

	NumberFieldX     int32 = 185
	NumberFieldWidth int32 = 80
)

var font *ttf.Font

type Menu struct {
	menuItems       []*MenuItem
	selectedElement int
	centered        bool
}

type MenuItem struct {
	texture     *sdl.Texture
	src, dst    *sdl.Rect
	numberField *NumberField
}

type NumberField struct {
	Title           *sdl.Texture
	Tsrc, Tdst      *sdl.Rect
	numberRect      *sdl.Rect
	Value, Min, Max int32
}

func (menu *Menu) Display(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(25, 25, 112, 255)
	renderer.Clear()

	first, last := menu.menuItems[0].dst.Y,
		menu.menuItems[len(menu.menuItems)-1].dst.Y
	last = newScreenHeight - last

	diff := (last + first) / 2
	diff -= first

	for i := range menu.menuItems {
		item := menu.menuItems[i]
		if menu.centered {
			item.dst.X = newScreenWidth/2 - item.dst.W/2
		} else {
			item.dst.X = newScreenWidth/2 - MenuXOffset
		}
		item.dst.Y += diff
		renderer.Copy(item.texture, item.src, item.dst)
	}

	sel := menu.selectedElement % len(menu.menuItems)
	x := menu.menuItems[sel].dst.X - MenuXOffset/38
	y := menu.menuItems[sel].dst.Y + menu.menuItems[sel].dst.H/2 - 2
	renderer.SetDrawColor(255, 255, 255, 255)
	for i := int32(0); i < MenuArrowSize/85; i++ {
		renderer.DrawLine(int(x-i), int(y-i), int(x-i), int(y+i))
	}

	renderer.Present()
}

func (menu *MenuItem) SetNumber(n int32, renderer *sdl.Renderer) {
	if menu.numberField != nil {
		nf := menu.numberField
		if n < nf.Min {
			n = nf.Min
		} else if n > nf.Max {
			n = nf.Max
		}
		nf.Value = n
		renderer.SetRenderTarget(menu.texture)
		defer renderer.SetRenderTarget(nil)

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		renderer.Copy(nf.Title, nf.Tsrc, nf.Tdst)
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.FillRect(nf.numberRect)

		numbers.WriteNumber(int64(n),
			nf.numberRect.X+nf.numberRect.W/2, 0, true, renderer)
	}
}

func InitText(renderer *sdl.Renderer) {
	var err error
	if ttf.WasInit() || font != nil {
		panic("Should only be called once!")
	}
	PanicOnError(ttf.Init())
	font, err = ttf.OpenFont("./font/Play-Bold.ttf", 20)
	PanicOnError(err)
}

func GetMenus(renderer *sdl.Renderer) []*Menu {
	ret := make([]*Menu, 5)

	mult := int32(1)
	if Arcade {
		mult = 2
	}

	ret[0] = &Menu{[]*MenuItem{
		GetMenuItem("1 Player", screenHeight/2-120*mult, renderer),
		GetMenuItem("2 Players", screenHeight/2-80*mult, renderer),
		GetMenuItem("Training (not implemented)", screenHeight/2-40*mult, renderer),
		GetMenuItem("High-Scores", screenHeight/2, renderer),
		GetMenuItem("Options", screenHeight/2+40*mult, renderer),
		GetMenuItem("Credits", screenHeight/2+80*mult, renderer),
		GetMenuItem("Quit", screenHeight/2+120*mult, renderer),
	}, 0, false}
	ret[1] = &Menu{[]*MenuItem{
		GetMenuItem("Beginner", screenHeight/2-80*mult, renderer),
		GetMenuItem("Intermediate", screenHeight/2-40*mult, renderer),
		GetMenuItem("Advanced", screenHeight/2, renderer),
		GetMenuItem("Beginner's Adventure", screenHeight/2+40*mult, renderer),
		GetMenuItem("Intermediate's Adventure", screenHeight/2+80*mult, renderer),
	}, 0, false}
	ret[2] = &Menu{[]*MenuItem{
		GetMenuItem("Set Name", screenHeight/2-80*mult, renderer),
		GetMenuItem("Highscores", screenHeight/2-40*mult, renderer),
		GetMenuItem("Continue", screenHeight/2, renderer),
		GetMenuItem("Restart", screenHeight/2+40*mult, renderer),
		GetMenuItem("Exit to menu", screenHeight/2+80*mult, renderer),
	}, 0, false}
	ret[3] = &Menu{[]*MenuItem{
		GetNumberMenuItem("Character (P1)", int32(options.CharacterP1), 0, 3,
			screenHeight/2-100*mult, renderer),
		GetNumberMenuItem("Character (P2)", int32(options.CharacterP2), 0, 3,
			screenHeight/2-60*mult, renderer),
		GetNumberMenuItem("EdgeSlip", int32(options.EdgeSlip), 0, 16,
			screenHeight/2-20*mult, renderer),
		GetNumberMenuItem("BetterSlip", int32(options.BetterSlip), 0, 512,
			screenHeight/2+20*mult, renderer),
		GetNumberMenuItem("Show Divert", int32(options.ShowDivert), 0, 1,
			screenHeight/2+60*mult, renderer),
		GetMenuItem("Reset", screenHeight/2+100*mult, renderer),
	}, 0, true}
	ret[4] = &Menu{[]*MenuItem{
		GetNumberMenuItem("Level", 0, 0, 34,
			screenHeight/2-40, renderer),
		GetNumberMenuItem("Difficulty", 0, 0, 2,
			screenHeight/2, renderer),
		GetNumberMenuItem("Lives", 3, 1, 16,
			screenHeight/2+40, renderer),
	}, 0, true}

	return ret
}

func GetMenuItem(text string, y int32, renderer *sdl.Renderer) *MenuItem {
	texture, src, dst := GetText(text, sdl.Color{0, 190, 0, 255},
		-1, y, renderer)
	if Arcade {
		dst.W, dst.H = dst.W*2, dst.H*2
	}
	return &MenuItem{texture, src, dst, nil}
}

func GetNumberMenuItem(text string, value, min, max int32,
	y int32, renderer *sdl.Renderer) *MenuItem {

	title, tsrc, tdst := GetText(text, sdl.Color{0, 190, 0, 255},
		0, 0, renderer)

	numberRect := &sdl.Rect{NumberFieldX, 0, NumberFieldWidth, tdst.H}
	numberField := &NumberField{title, tsrc, tdst, numberRect,
		value, min, max}

	src := &sdl.Rect{0, 0, numberRect.X + numberRect.W, numberRect.H}
	dst := &sdl.Rect{-1, y + tdst.Y, src.W, src.H}
	if Arcade {
		dst.W, dst.H = dst.W*2, dst.H*2
	}
	tdst.Y = 0

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(src.W), int(src.H))
	PanicOnError(err)
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	menuItem := &MenuItem{texture, src, dst, numberField}
	menuItem.SetNumber(value, renderer)

	return menuItem
}

func GetText(text string, color sdl.Color, x, y int32,
	renderer *sdl.Renderer) (*sdl.Texture, *sdl.Rect, *sdl.Rect) {
	textSurface, err := font.RenderUTF8_Solid(text, color)
	PanicOnError(err)
	defer textSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(textSurface)
	PanicOnError(err)
	src := &sdl.Rect{0, 0, textSurface.W, textSurface.H}
	dst := &sdl.Rect{x, y - src.H/2, src.W, src.H}
	return texture, src, dst
}

func (menu *Menu) Run(renderer *sdl.Renderer, input *Input) int {
	vStepper, mStepper := input.mono.upDown.Stepper(20, 5),
		input.mono.leftRight.Stepper(20, 4)
	input.mono.a.down = false
	input.mono.b.down = false

	for !quit {
		selected := menu.menuItems[menu.selectedElement]
		menu.Display(renderer)
		input.Poll()

		if input.mono.b.down {
			return -1
		}
		if input.mono.a.down && selected.numberField == nil {
			break
		}

		mod := mStepper()
		if mod != 0 {
			if selected.numberField != nil {
				if input.mono.a.down {
					mod *= 10
				}
				selected.SetNumber(selected.numberField.Value+mod,
					renderer)
			}
		}

		val := vStepper()
		if val != 0 {
			if val > 0 {
				menu.selectedElement = (menu.selectedElement +
					1) % len(menu.menuItems)
			} else if val < 0 {
				menu.selectedElement = (menu.selectedElement +
					len(menu.menuItems) - 1) % len(menu.menuItems)
			}
		}

	}
	return menu.selectedElement
}

func (menu *Menu) Free() {
	for i := 0; i < len(menu.menuItems); i++ {
		if menu.menuItems[i].texture != nil {
			menu.menuItems[i].texture.Destroy()
			if menu.menuItems[i].numberField != nil {
				if menu.menuItems[i].numberField.Title != nil {
					menu.menuItems[i].numberField.Title.Destroy()
				}
				menu.menuItems[i].numberField = nil
			}
		}
	}
	menu.menuItems = nil
}
