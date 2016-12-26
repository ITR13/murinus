package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var font *ttf.Font

type Menu struct {
	menuItems       []*MenuItem
	selectedElement int
}

type MenuItem struct {
	texture  *sdl.Texture
	src, dst *sdl.Rect
}

func (menu *Menu) Display(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(25, 25, 112, 255)
	renderer.Clear()
	for i := range menu.menuItems {
		item := menu.menuItems[i]
		renderer.Copy(item.texture, item.src, item.dst)
	}

	sel := menu.selectedElement % len(menu.menuItems)
	x := menu.menuItems[sel].dst.X - (10*sizeMult)/sizeDiv
	y := menu.menuItems[sel].dst.Y + menu.menuItems[sel].dst.H/2 - 2
	renderer.SetDrawColor(255, 255, 255, 255)
	for i := int32(0); i < (15*sizeMult)/sizeDiv; i++ {
		renderer.DrawLine(int(x-i), int(y-i), int(x-i), int(y+i))
	}

	renderer.Present()
}

func GetMenus(renderer *sdl.Renderer) []*Menu {
	var err error
	if ttf.WasInit() || font != nil {
		panic("Should only be called once!")
	}
	e(ttf.Init())
	font, err = ttf.OpenFont("./font/AverageMono.ttf", 20)
	e(err)
	ret := make([]*Menu, 3)

	ret[0] = &Menu{[]*MenuItem{
		GetMenuItem("1 Player", screenWidth/2-screenWidth/8,
			screenHeight/2-80, renderer),
		GetMenuItem("2 Players", screenWidth/2-screenWidth/8,
			screenHeight/2-40, renderer),
		GetMenuItem("High-Scores", screenWidth/2-screenWidth/8,
			screenHeight/2, renderer),
		GetMenuItem("Options", screenWidth/2-screenWidth/8,
			screenHeight/2+40, renderer),
		GetMenuItem("Quit", screenWidth/2-screenWidth/8,
			screenHeight/2+80, renderer),
	}, 0}
	ret[1] = &Menu{[]*MenuItem{
		GetMenuItem("Easy", screenWidth/2-screenWidth/8,
			screenHeight/2-40, renderer),
		GetMenuItem("Medium", screenWidth/2-screenWidth/8,
			screenHeight/2, renderer),
		GetMenuItem("Hard", screenWidth/2-screenWidth/8,
			screenHeight/2+40, renderer),
	}, 1}
	ret[2] = &Menu{[]*MenuItem{
		GetMenuItem("Set Name", screenWidth/2-screenWidth/8,
			screenHeight/2-80, renderer),
		GetMenuItem("Highscores", screenWidth/2-screenWidth/8,
			screenHeight/2-40, renderer),
		GetMenuItem("Continue", screenWidth/2-screenWidth/8,
			screenHeight/2, renderer),
		GetMenuItem("Retry", screenWidth/2-screenWidth/8,
			screenHeight/2+40, renderer),
		GetMenuItem("Exit to menu", screenWidth/2-screenWidth/8,
			screenHeight/2+80, renderer),
	}, 0}

	return ret
}

func GetMenuItem(text string, x, y int32, renderer *sdl.Renderer) *MenuItem {
	texture, src, dst := GetText(text, sdl.Color{0, 190, 0, 255},
		x, y, renderer)
	return &MenuItem{texture, src, dst}
}

func GetText(text string, color sdl.Color, x, y int32,
	renderer *sdl.Renderer) (*sdl.Texture, *sdl.Rect, *sdl.Rect) {
	textSurface, err := font.RenderUTF8_Solid(text, color)
	e(err)
	defer textSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(textSurface)
	e(err)
	src := &sdl.Rect{0, 0, textSurface.W, textSurface.H}
	dst := &sdl.Rect{x, y - src.H/2, src.W, src.H}
	return texture, src, dst
}

func (menu *Menu) Run(renderer *sdl.Renderer, input *Input) int {
	prevVal := int32(0)
	step := 0
	repeat := true
	input.mono.a.down = false
	input.mono.b.down = false
	for !input.mono.a.down && !quit {
		if input.mono.b.down {
			return -1
		}
		menu.Display(renderer)
		input.Poll()
		val := input.mono.upDown.Val()
		if val != prevVal || repeat {
			repeat = false
			if prevVal != val {
				step = 0
			}
			prevVal = val
			if val > 0 {
				menu.selectedElement = (menu.selectedElement +
					1) % len(menu.menuItems)
			} else if val < 0 {
				menu.selectedElement = (menu.selectedElement +
					len(menu.menuItems) - 1) % len(menu.menuItems)
			}
		}
		step++
		if step > 20 && step%3 == 0 {
			repeat = true
		}
	}
	return menu.selectedElement
}
