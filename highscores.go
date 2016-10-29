package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type HighscoreList struct {
	scores       []*ScoreData
	uniqueScores []*ScoreData
}

type ScoreData struct {
	score         uint64
	name          string
	levelsCleared int
}

func GetHighscoreList() HighscoreList {
	return HighscoreList{
		make([]*ScoreData, 0),
		make([]*ScoreData, 0),
	}
}

func GetName(characters int, renderer *sdl.Renderer, input *Input) string {
	currentCharacter := 0
	//charsToChooseFrom := []string{" ", ",", ".", "-"}

	for !quit && currentCharacter != characters {

	}
	return ""
}

func (list *HighscoreList) Get() *ScoreData {
	return &ScoreData{0, "", 0}
}

func (list *HighscoreList) Add(score *ScoreData) {
	list.scores = append(list.scores, score)
	for i := range list.uniqueScores {
		if list.uniqueScores[i].name == score.name {
			if list.uniqueScores[i].score < score.score {
				list.uniqueScores[i] = score
			} else if list.uniqueScores[i].score == score.score {
				if list.uniqueScores[i].levelsCleared < score.levelsCleared {
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
	textureHeight := screenHeight / 7
	names := make([]*sdl.Texture, 9)
	for i := 0; i < len(names); i++ {
		names[i] = list.RenderScore(i-1, false, renderer)
		defer names[i].Destroy()
		renderer.SetRenderTarget(names[i])
		renderer.SetDrawColor(uint8(i*255/9), 255, 255, 255)
		renderer.Clear()
	}

	src := &sdl.Rect{0, 0, screenWidth, textureHeight}
	dst := &sdl.Rect{0, 0, screenWidth, textureHeight}
	renderer.SetRenderTarget(nil)
	for !input.mono.a.down && !input.mono.b.down && !quit {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for i := 0; i < len(names); i++ {
			if names[i] != nil {
				y := textureHeight*int32(i-1) + subPixel
				dst.Y = y
				renderer.Copy(names[i], src, dst)
			}
		}
		renderer.Present()
		sdl.Delay(4)
		input.Poll()
		dir := input.mono.upDown.Val()
		if dir != 0 {
			subPixel += dir * textureHeight / 16
			for subPixel < 0 {
				subPixel += textureHeight
				currentIndex--
			}
			for subPixel >= textureHeight {
				subPixel -= textureHeight
				currentIndex++
			}
		}

	}
}

func (list *HighscoreList) RenderScore(index int, unique bool,
	renderer *sdl.Renderer) *sdl.Texture {
	if index < 0 {
		return nil
	}
	if unique {
		if index >= len(list.uniqueScores) {
			return nil
		}
		return list.uniqueScores[index].Render(renderer)
	} else {
		if index >= len(list.scores) {
			return nil
		}
		return list.scores[index].Render(renderer)
	}
}

func (score *ScoreData) Render(renderer *sdl.Renderer) *sdl.Texture {
	surface, err := font.RenderUTF8_Solid("Sample Text",
		sdl.Color{255, 255, 255, 255})
	e(err)
	defer surface.Free()
	texture, err := renderer.CreateTextureFromSurface(surface)
	e(err)
	return texture
}
