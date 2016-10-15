package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Tile uint8

const (
	Empty   Tile = iota
	Wall    Tile = iota
	Point   Tile = iota
	Powerup Tile = iota
	p200    Tile = iota
	p500    Tile = iota
	p1000   Tile = iota
	p2000   Tile = iota
)

type Stage struct {
	tiles   *TileStage
	sprites *SpriteStage
}

type TileStage struct {
	renderedOnce bool
	tileInfo     *TileInfo
	tiles        [][]Tile
	texture      *sdl.Texture
	src, dst     *sdl.Rect
	tileDst      *sdl.Rect
	w, h         int32
}

type TileInfo struct {
	dst      *sdl.Rect
	textures []*sdl.Texture
	src      []*sdl.Rect
}

type SpriteStage struct {
	sprites   []*Sprite
	entities  []*Entity
	texture   *sdl.Texture
	src, dst  *sdl.Rect
	spriteDst *sdl.Rect
	time      int64
}

type Entity struct {
	sprite *Sprite
	x, y   int32
}

type Sprite struct {
	texture  *sdl.Texture
	src      []*sdl.Rect
	timeDiv  int64
	priority int
}

func (stage *Stage) Render(renderer *sdl.Renderer) {
	if !stage.tiles.renderedOnce {
		stage.tiles.Render(renderer)
	}
	stage.sprites.Render(renderer)
	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.Clear()
	renderer.Copy(stage.tiles.texture, stage.tiles.src, stage.tiles.dst)
	renderer.Copy(stage.sprites.texture, stage.sprites.src, stage.sprites.dst)
	renderer.Present()
}

func (tiles *TileStage) Render(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(tiles.texture)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	if tiles.tiles != nil {
		for x := int32(0); x < tiles.w; x++ {
			tiles.tileDst.X = x * size
			for y := int32(0); y < tiles.h; y++ {
				tiles.tileDst.Y = y * size
				tiles.tileDst.W = size
				tiles.tileDst.H = size
				renderer.Copy(tiles.tileInfo.textures[tiles.tiles[x][y]],
					tiles.tileInfo.src[tiles.tiles[x][y]], tiles.tileDst)
			}
		}
		tiles.renderedOnce = true
	}
}

func (sprites *SpriteStage) Render(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(sprites.texture)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	if sprites.entities != nil {
		for priority := 0; priority <= 10; priority++ {
			for i := 0; i < len(sprites.entities); i++ {
				e := sprites.entities[i]
				s := e.sprite
				if s.priority == priority {
					sprites.spriteDst.X = e.x * size
					sprites.spriteDst.Y = e.y * size
					sprites.spriteDst.W = size
					sprites.spriteDst.H = size
					t := (sprites.time / s.timeDiv) % int64(len(s.src))
					renderer.Copy(s.texture,
						s.src[t], sprites.spriteDst)
				}
			}
		}
	}
}

func LoadTextures(width, height int32, renderer *sdl.Renderer) *Stage {
	rect4x4 := sdl.Rect{0, 0, 4, 4}
	stageRect := sdl.Rect{0, 0, width * size, height * size}
	offsetFromScreenX := (screenWidth - width*size) / 2
	offsetFromScreenY := (screenHeight - height*size) / 2
	stageScreenRect := sdl.Rect{offsetFromScreenX, offsetFromScreenY,
		width * size, height * size}

	tileTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(width*size), int(height*size))
	e(err)
	tileInfo := TileInfo{&sdl.Rect{},
		make([]*sdl.Texture, 8),
		make([]*sdl.Rect, 8)}
	for i := 0; i < 8; i++ {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
			sdl.TEXTUREACCESS_TARGET, 4, 4)
		e(err)
		tileInfo.textures[i] = texture
		tileInfo.src[i] = &rect4x4
	}
	renderer.SetRenderTarget(tileInfo.textures[0])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetRenderTarget(tileInfo.textures[1])
	renderer.SetDrawColor(25, 25, 112, 255)
	renderer.Clear()
	renderer.SetRenderTarget(tileInfo.textures[2])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(240, 230, 140, 255)
	renderer.DrawRect(&sdl.Rect{1, 1, 2, 2})

	tileStage := TileStage{false, &tileInfo, nil,
		tileTexture, &stageRect, &stageScreenRect,
		&sdl.Rect{}, width, height}

	spriteTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(width*size), int(height*size))
	e(err)
	spriteTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	spriteDatas := make([]*Sprite, 2)
	for i := 0; i < len(spriteDatas); i++ {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
			sdl.TEXTUREACCESS_TARGET, 4, 4)
		texture.SetBlendMode(sdl.BLENDMODE_BLEND)
		e(err)
		spriteDatas[i] = &Sprite{texture, []*sdl.Rect{&rect4x4}, 1, 10 - i}
	}
	renderer.SetRenderTarget(spriteDatas[0].texture)
	renderer.SetDrawColor(0, 95, 0, 255)
	renderer.Clear()
	renderer.SetRenderTarget(spriteDatas[1].texture)
	renderer.SetDrawColor(0, 127, 0, 255)
	renderer.Clear()

	spriteStage := SpriteStage{spriteDatas, nil, spriteTexture,
		&stageRect, &stageScreenRect, &sdl.Rect{}, 0}

	return &Stage{&tileStage, &spriteStage}
}
