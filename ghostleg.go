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
	hits := [3]int{1, 1, 1}
	for i := range flips {
		flips[i] = make([]bool, width-1)
		total := hits[0] + hits[1]*2 + hits[2]*3
		chances := [3]int{total - hits[0],
			total - hits[1]*2, total - hits[2]*3}
		left := 0
		for j := 0; j < len(flips[i]); j++ {
			if left == 0 {
				r := rand.Int() % total
				for k := range chances {
					if chances[k] > r {
						left = k
						hits[k]++
						break
					}
					r -= chances[k]
				}
			} else {
				left--
			}
			if left == 0 {
				flips[i][j] = true
				j++
			}
		}
	}

	return &GhostLeg{width, flips, rand.Int() % width}
}

func (gl *GhostLeg) Display(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.SetRenderTarget(nil)
	renderer.Clear()
	renderer.SetDrawColor(255, 255, 255, 255)

	s := gl.end

	for i := range gl.flips {
		for j := 0; j < gl.legs; j++ {
			if j == s {
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
				if j == s {
					renderer.SetDrawColor(255, 0, 0, 255)
					s++
				} else if j == s-1 {
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
	renderer.Present()
}
