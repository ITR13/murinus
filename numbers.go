package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type NumberData struct {
	texture *sdl.Texture
	src     []*sdl.Rect
	dst     *sdl.Rect
	W       int32
}

var numbers *NumberData

func InitNumbers(renderer *sdl.Renderer) {
	textSurface, err := font.RenderUTF8_Solid("0123456789-",
		sdl.Color{255, 255, 255, 255})
	e(err)
	defer textSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(textSurface)
	e(err)

	W, H := textSurface.W/11, textSurface.H
	src := make([]*sdl.Rect, 11)
	for i := int32(0); i < 11; i++ {
		src[i] = &sdl.Rect{i * W, 0, W, H}
	}
	dst := &sdl.Rect{0, 0, W, H}
	numbers = &NumberData{texture, src, dst, W}
}

//Note, needs renderer to have correct surface set before it's called
func (numberData *NumberData) WriteNumber(n int64, x, y int32, center bool,
	renderer *sdl.Renderer) {
	W := numberData.W
	digits, negative := digitsIn(n)

	if center {
		x += W * digits / 2
	} else {
		x += W * digits
	}

	for n > 0 {
		numberData.dst.X, numberData.dst.Y = x, y
		renderer.Copy(numberData.texture, numberData.src[n%10], numberData.dst)
		n, x = n/10, x-W
	}

	if negative {
		numberData.dst.X, numberData.dst.Y = x, y
		renderer.Copy(numberData.texture, numberData.src[10], numberData.dst)
	}
}

func digitsIn(n int64) (int32, bool) {
	ret := int32(0)
	negative := n < 0
	if negative {
		n = -n
		ret++
	}
	for n >= 10 {
		n /= 10
		ret++
	}
	return ret, negative
}
