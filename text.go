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

func InitText(renderer *sdl.Renderer) {
	var err error
	if ttf.WasInit() || font != nil {
		panic("Should only be called once!")
	}
	PanicOnError(ttf.Init())
	font, err = ttf.OpenFont("./font/Play-Bold.ttf", 20)
	PanicOnError(err)
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

func WriteText(text string, color sdl.Color, x, y int32,
	centered bool, renderer *sdl.Renderer) {
	texture, src, dst := GetText(text, color, x, y, renderer)
	defer texture.Destroy()

	if centered {
		dst.X += newScreenWidth/2 - dst.W/2
		dst.Y++
	}

	renderer.SetRenderTarget(nil)
	renderer.Copy(texture, src, dst)
}
