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

/*
import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type GhostLeg struct {
	legs  int
	flips [][]bool
	end   int
}

func GetRandomGhostLeg(width, height int) *GhostLeg {
	flips := make([][]bool, height)
	hits := make([]int, width)

	c := 2
	for i := range flips {
		flips[i] = make([]bool, width-1)
		for j := 0; j < len(flips[i]); j++ {
			top, bottom := hits[j]+hits[j+1], c
			if top < 0 {
				bottom -= top
				top = 1
			} else {
				bottom += top
			}

			if rand.Int()%bottom < top {
				c++
				hits[j], hits[j+1] = hits[j+1]-3, hits[j]-3
			} else {
				c = 2
				flips[i][j] = true
				hits[j], hits[j+1] = hits[j+1]+3, hits[j]+3
				j++
				if len(hits) < j+1 {
					hits[j]--
				}
			}
		}
	}

	return &GhostLeg{width, flips, rand.Int() % width}
}

func (gl *GhostLeg) DisplayPath(path int, renderer *sdl.Renderer) {
	if newScreenWidth != screenWidth || newScreenHeight != screenHeight {
		SetWindowSize(newScreenWidth, newScreenHeight, stage)
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.SetRenderTarget(nil)
	renderer.Clear()
	renderer.SetDrawColor(255, 255, 255, 255)

	for i := range gl.flips {
		for j := 0; j < gl.legs; j++ {
			if j == path {
				renderer.SetDrawColor(255, 0, 0, 255)
			} else {
				renderer.SetDrawColor(255, 255, 255, 255)
			}
			renderer.FillRect(&sdl.Rect{
				int32(j) * screenWidth / int32(gl.legs),
				int32(i) * screenHeight / int32(len(gl.flips)),
				4, screenHeight / int32(gl.legs*4)})
		}

		for j := range gl.flips[i] {
			if gl.flips[i][j] {
				if j == path {
					renderer.SetDrawColor(255, 0, 0, 255)
					path++
				} else if j == path-1 {
					renderer.SetDrawColor(255, 0, 0, 255)
					path--
				} else {
					renderer.SetDrawColor(255, 255, 255, 255)
				}
				renderer.FillRect(&sdl.Rect{
					int32(j) * screenWidth / int32(gl.legs),
					int32(i+1) * screenHeight / int32(len(gl.flips)),
					screenWidth / int32(gl.legs), 4})
			}
		}
	}
	for j := 0; j < gl.legs; j++ {
		renderer.FillRect(&sdl.Rect{
			int32(j) * screenWidth / int32(gl.legs),
			int32(len(gl.flips)) * screenHeight / int32(len(gl.flips)),
			4, screenHeight / int32(gl.legs)})
	}

	renderer.Present()
}

func (gl *GhostLeg) AnimatePath(path int) func(*sdl.Renderer) (int, int) {
	depth := 0
	skipped := 0
	return func(renderer *sdl.Renderer) (int, int) {
		s := path
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.SetDrawColor(255, 255, 255, 255)

		for i := range gl.flips {
			for j := 0; j < gl.legs; j++ {
				if j == s && i <= depth/2 {
					renderer.SetDrawColor(255, 0, 0, 255)
				} else {
					renderer.SetDrawColor(255, 255, 255, 255)
				}
				renderer.FillRect(&sdl.Rect{
					int32(j) * screenWidth / int32(gl.legs),
					int32(i) * screenHeight / int32(len(gl.flips)),
					4, screenHeight / int32(gl.legs)})
			}

			for j := range gl.flips[i] {
				if gl.flips[i][j] {
					if j == s && i*2 < depth {
						renderer.SetDrawColor(255, 0, 0, 255)
						s++
					} else if j == s-1 && i*2 < depth {
						renderer.SetDrawColor(255, 0, 0, 255)
						s--
					} else {
						renderer.SetDrawColor(255, 255, 255, 255)
					}
					renderer.FillRect(&sdl.Rect{
						int32(j) * screenWidth / int32(gl.legs),
						int32(i+1) * screenHeight / int32(len(gl.flips)),
						screenWidth / int32(gl.legs), 4})
				}
			}
		}

		if depth%2 != 0 {
			depth++
		} else if depth/2 < len(gl.flips) &&
			((s > 0 && gl.flips[depth/2][s-1]) ||
				(s < len(gl.flips[depth/2]) && gl.flips[depth/2][s])) {
			depth++
		} else {
			skipped++
			depth += 2
		}

		if depth < len(gl.flips)*2 {
			return -1, 1
		}
		return s, skipped
	}
}

func (gl *GhostLeg) Choose(renderer *sdl.Renderer) int {

}

func (gl *GhostLeg) Play(renderer *sdl.Renderer) {

}
*/
