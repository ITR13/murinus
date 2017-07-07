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
	"fmt"
	"os"

	"encoding/gob"

	"github.com/veandco/go-sdl2/sdl"
)

type Stats struct {
	Singleplayer, Multiplayer                      [5]uint64
	Deaths, ExtraLives, Points, LevelsStarted      uint64
	TimesTrained, EdgeSlip, QuickCourner, Diagonal uint64
}

var stats Stats

func ShowStats() {
	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	white := sdl.Color{255, 255, 255, 255}

	sp, mp := stats.Singleplayer[0]+stats.Singleplayer[1]+
		stats.Singleplayer[2]+stats.Singleplayer[3]+stats.Singleplayer[4],
		stats.Multiplayer[0]+stats.Multiplayer[1]+stats.Multiplayer[2]+
			stats.Multiplayer[3]+stats.Multiplayer[4]

	WriteText("Singleplayer", white, 20, 20, false, renderer)
	WriteText(fmt.Sprint(sp), white, 0, 20, true, renderer)
	WriteText("Multiplayer", white, 20, 40, false, renderer)
	WriteText(fmt.Sprint(mp), white, 0, 40, true, renderer)
	WriteText("Death", white, 20, 60, false, renderer)
	WriteText(fmt.Sprint(stats.Deaths), white, 0, 60, true, renderer)
	WriteText("Extra Lives", white, 20, 80, false, renderer)
	WriteText(fmt.Sprint(stats.ExtraLives), white, 0, 80, true, renderer)
	WriteText("Points", white, 20, 100, false, renderer)
	WriteText(fmt.Sprint(stats.Points), white, 0, 100, true, renderer)
	WriteText("Levels Started", white, 20, 120, false, renderer)
	WriteText(fmt.Sprint(stats.LevelsStarted), white, 0, 120, true, renderer)
	WriteText("Times Trained", white, 20, 140, false, renderer)
	WriteText(fmt.Sprint(stats.TimesTrained), white, 0, 140, true, renderer)
	WriteText("EdgeSlip", white, 20, 160, false, renderer)
	WriteText(fmt.Sprint(stats.EdgeSlip), white, 0, 160, true, renderer)
	WriteText("QuickCourner", white, 20, 180, false, renderer)
	WriteText(fmt.Sprint(stats.QuickCourner), white, 0, 180, true, renderer)
	WriteText("Diagonal", white, 20, 200, false, renderer)
	WriteText(fmt.Sprint(stats.Diagonal), white, 0, 200, true, renderer)

	for i := 0; i < 5; i++ {
		WriteText("Singleplayer - "+menus[1].menuItems[i].text,
			white, 20, 240+int32(i*20), false, renderer)
		WriteText(fmt.Sprint(stats.Singleplayer[i]),
			white, 0, 240+int32(i*20), true, renderer)

		WriteText("Multiplayer  - "+menus[1].menuItems[i].text,
			white, 20, 360+int32(i*20), false, renderer)
		WriteText(fmt.Sprint(stats.Multiplayer[i]),
			white, 0, 360+int32(i*20), true, renderer)
	}

	renderer.SetDrawColor(55, 55, 95, 255)
	for i := 11; i < 460; i += 20 {
		renderer.DrawLine(0, i, int(screenWidth), i)
	}

	renderer.Present()

	input.Poll()
	input.mono.a.Clear()
	input.mono.b.Clear()
	for !(quit || input.mono.a.Down() || input.mono.b.Down()) {
		input.Poll()
	}

}

func (stats *Stats) Save(path string) {
	file, err := os.Create(path)
	if LogOnError(err) {
		return
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(stats)
}

func (stats *Stats) Load(path string) {
	file, err := os.Open(path)
	if LogOnError(err) {
		return
	}
	decoder := gob.NewDecoder(file)
	LogOnError(decoder.Decode(stats))
}
