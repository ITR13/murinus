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
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type HighscoreList struct {
	scores       []*ScoreData
	uniqueScores []*ScoreData
}

type ScoreData struct {
	Score         int64
	Name          string
	LevelsCleared int
	Difficulty    int
	Date          time.Time
}

func GetName(defaultName string, renderer *sdl.Renderer, input *Input) string {
	characters := int32(len(defaultName))
	input.mono.a.down = false
	input.mono.b.down = false
	currentCharacter := int32(0)
	charList := make([][13]byte, characters)
	for i := 0; i < len(charList); i++ {
		for j := 0; j < 13; j++ {
			charList[i][j] = defaultName[i] + byte(j-6)
		}
	}
	curCharUnderline := &sdl.Rect{0, screenHeight/2 + 12, 12, 4}

	draw := func() {
		renderer.SetRenderTarget(nil)
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for y := int32(0); y < 13; y++ {
			for x := int32(0); x < characters; x++ {
				Y := (y - 6)
				c := y - 6
				if Y < 0 {
					Y--
					c = -c
				} else if Y > 0 {
					Y++
				}
				c = 255 - c*16

				texture, src, dst := GetText(string(charList[x][y]),
					sdl.Color{uint8(c), uint8(c), uint8(c), 255},
					x*40-40*characters/2+screenWidth/2,
					Y*24+screenHeight/2, renderer)
				renderer.Copy(texture, src, dst)
				texture.Destroy()
			}
		}
		renderer.SetDrawColor(215, 10, 10, 255)
		curCharUnderline.X = currentCharacter*40 - 40*characters/2 +
			screenWidth/2
		renderer.FillRect(curCharUnderline)
	}

	prevLR := int32(0)
	prevUD := int32(0)
	for !quit && currentCharacter != characters {
		draw()
		renderer.Present()
		input.Poll()
		ud := input.mono.upDown.Val()
		if ud != prevUD {
			prevUD = ud
			if ud < 0 {
				for i := 12; i > 0; i-- {
					charList[currentCharacter][i] =
						charList[currentCharacter][i-1]
				}
				charList[currentCharacter][0]--
				if charList[currentCharacter][0] < 32 {
					charList[currentCharacter][0] = 126
				}
			} else if ud > 0 {
				for i := 0; i < 12; i++ {
					charList[currentCharacter][i] =
						charList[currentCharacter][i+1]
				}
				charList[currentCharacter][12]++
				if charList[currentCharacter][12] > 126 {
					charList[currentCharacter][12] = 32
				}
			}
		}

		lr := input.mono.leftRight.Val()
		if lr != prevLR {
			prevLR = lr
			if lr > 0 {
				if currentCharacter < characters-1 {
					currentCharacter++
				}
			} else if lr < 0 {
				if currentCharacter > 0 {
					currentCharacter--
				}
			}
		}
		if input.mono.a.down {
			input.mono.a.down = false
			currentCharacter++
		}
		if input.mono.b.down {
			input.mono.b.down = false
			currentCharacter--
			if currentCharacter < 0 {
				return ""
			}
		}
	}
	name := ""
	for i := 0; i < len(charList); i++ {
		name += string(charList[i][6])
	}
	return name
}

func (list *HighscoreList) Add(score *ScoreData) {
	list.scores = append(list.scores, score)
	for i := range list.uniqueScores {
		if list.uniqueScores[i].Name == score.Name {
			if list.uniqueScores[i].Score < score.Score {
				list.uniqueScores[i] = score
			} else if list.uniqueScores[i].Score == score.Score {
				if list.uniqueScores[i].LevelsCleared < score.LevelsCleared {
					list.uniqueScores[i] = score
				}
			}
			return
		}
	}
	list.uniqueScores = append(list.uniqueScores, score)
}

func (list *HighscoreList) Display(renderer *sdl.Renderer, input *Input) {
	input.mono.a.down = false
	input.mono.b.down = false
	subPixel := int32(0)
	currentIndex := -1
	storedIndex := -1
	unique := false
	l := 18
	textureHeight := screenHeight / int32(l-2)
	names := make([]*sdl.Texture, l)
	for i := 0; i < len(names); i++ {
		names[i] = list.RenderScore(i-1, false, true, renderer)
	}
	src := &sdl.Rect{0, 0, screenWidth, textureHeight}
	dst := &sdl.Rect{0, 0, screenWidth, textureHeight}
	renderer.SetRenderTarget(nil)
	scrollMult := int32(210)
	update := false
	for !input.mono.a.down && !input.mono.b.down && !quit {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for i := 0; i < len(names); i++ {
			if names[i] != nil {
				y := textureHeight*int32(i-1) + subPixel
				_, _, w, h, err := names[i].Query()
				e(err)
				dst.Y = y
				dst.W, dst.H = w, h
				src.W, src.H = w, h
				renderer.Copy(names[i], src, dst)
			}
		}
		renderer.Present()
		input.Poll()
		dir := -input.mono.upDown.Val()
		if dir != 0 {
			subPixel += scrollMult * dir * textureHeight / (5 * 210)
			scrollMult++
			for subPixel < 0 {
				subPixel += textureHeight
				currentIndex++
				update = true
			}
			for subPixel >= textureHeight {
				subPixel -= textureHeight
				currentIndex--
				update = true
			}
			if unique {
				if currentIndex < -l-2 {
					currentIndex = len(list.uniqueScores) + l + 2
				} else if currentIndex > len(list.uniqueScores)+l+2 {
					currentIndex = -l - 2
				}
			} else {
				if currentIndex < -16 {
					currentIndex = len(list.scores) + 16
				} else if currentIndex > len(list.scores)+16 {
					currentIndex = -16
				}
			}
		} else {
			scrollMult = 210
		}
		val := input.mono.leftRight.Val()
		if (val > 0 && !unique) || (val < 0 && unique) {
			unique = !unique
			if unique {
				storedIndex = currentIndex
			} else {
				currentIndex = storedIndex
			}
			update = true
		}
		if update {
			update = false
			for i := 0; i < len(names); i++ {
				names[i].Destroy()
				names[i] = list.RenderScore(i+currentIndex,
					unique, true, renderer)
			}
		}
	}
	for i := 0; i < len(names); i++ {
		if names[i] != nil {
			names[i].Destroy()
		}
	}
}

func (list *HighscoreList) RenderScore(index int, unique, multi bool,
	renderer *sdl.Renderer) *sdl.Texture {
	if index < 0 {
		return nil
	}
	if unique {
		if index >= len(list.uniqueScores) {
			return nil
		}
		return list.uniqueScores[index].Render(index, multi, renderer)
	} else {
		if index >= len(list.scores) {
			return nil
		}
		return list.scores[index].Render(index, multi, renderer)
	}
}

func (score *ScoreData) Render(i int, multi bool,
	renderer *sdl.Renderer) *sdl.Texture {
	text := "[" + strconv.Itoa(i+1) + "]"
	for len(text) < 5 {
		text += " "
	}
	text += score.Name + " | "

	for v := int64(1000000000); v > 0; v /= 1000 {
		if v > score.Score {
			text += "000"
		} else {
			val := int((score.Score / v) % 1000)
			if val < 10 {
				text += "00"
			} else if val < 100 {
				text += "0"
			}
			text += strconv.Itoa(val)
		}
		if v > 999 {
			text += "."
		}
	}
	text += fmt.Sprintf(" | %02d | %04d/%02d/%02d - %02d:%02d:%02d",
		score.LevelsCleared,
		score.Date.Year(), score.Date.Month(), score.Date.Day(),
		score.Date.Hour(), score.Date.Minute(), score.Date.Second())

	r, g, b := 255, 255, 255
	if i == 0 {
		r, g, b = 255, 255, 51
	} else if i == 1 {
		r, g, b = 255, 0, 0
	} else if i == 2 {
		r, g, b = 0, 190, 0
	}

	if multi {
		diff := score.Difficulty*2 + 2
		if diff > 6 {
			diff -= 5
		}
		r = r * (diff) / 8
		g = g * (diff) / 8
		if score.Difficulty == 0 {
			b = b * (diff + 1) / 10
		} else {
			b = b * (diff) / 8
		}
	}

	surface, err := font.RenderUTF8_Solid(text,
		sdl.Color{uint8(r), uint8(g), uint8(b), 255})
	e(err)
	defer surface.Free()
	texture, err := renderer.CreateTextureFromSurface(surface)
	e(err)
	return texture
}

func Read(path string) *HighscoreList {
	list := HighscoreList{
		make([]*ScoreData, 0),
		make([]*ScoreData, 0),
	}
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		e(err)
		defer file.Close()
		decoder := gob.NewDecoder(file)
		datas := make([]*ScoreData, 0)
		e(decoder.Decode(&datas))
		for i := 0; i < len(datas); i++ {
			list.Add(datas[i])
		}
		sort.Sort(SortByScore(list.scores))
		sort.Sort(SortByScore(list.uniqueScores))
	}
	return &list
}

func (list *HighscoreList) Write(path string) {
	file, err := os.Create(path)
	e(err)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	e(encoder.Encode(list.scores))
}

func (list *HighscoreList) Sort() {
	sort.Sort(SortByScore(list.scores))
	sort.Sort(SortByScore(list.uniqueScores))
}

type SortByScore []*ScoreData

func (s SortByScore) Len() int {
	return len(s)
}

func (s SortByScore) Less(i, j int) bool {
	if s[i].Score == s[j].Score {
		return s[i].LevelsCleared > s[j].LevelsCleared
	}
	return s[i].Score > s[j].Score
}

func (s SortByScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
