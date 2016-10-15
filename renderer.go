package main

import "github.com/veandco/go-sdl2/sdl"

const (
	size         int32 = 16
	screenWidth  int32 = 640
	screenHeight int32 = 480
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
	sprites  []*Sprite
	texture  *sdl.Texture
	src, dst *sdl.Rect
	time     int64
}

type Sprite struct {
	texture *sdl.Texture
	src     []*sdl.Rect
	x, y    int32
	timeDiv int64
}

func (stage *Stage) Render(renderer *sdl.Renderer) {
	if !stage.tiles.renderedOnce {
		stage.tiles.Render(renderer)
	}
	stage.sprites.Render(renderer)
	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Copy(stage.tiles.texture, stage.tiles.src, stage.tiles.dst)
	renderer.Copy(stage.sprites.texture, stage.sprites.src, stage.sprites.dst)
}

func (tiles *TileStage) Render(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(tiles.texture)
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	for x := int32(0); x < tiles.w; x++ {
		tiles.tileDst.X = x * size
		for y := int32(0); y < tiles.h; y++ {
			tiles.tileDst.Y = y * size
			tiles.tileDst.W = size
			tiles.tileDst.H = size
			renderer.Copy(tiles.tileInfo.textures[tiles.tiles[x][y]],
				tiles.tileInfo.src[tiles.tiles[x][y]], tiles.dst)
		}
	}
	tiles.renderedOnce = true
}

func (sprites *SpriteStage) Render(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(sprites.texture)
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	for i := 0; i < len(sprites.sprites); i++ {
		s := sprites.sprites[i]
		sprites.dst.X = s.x * size
		sprites.dst.Y = s.y * size
		renderer.Copy(s.texture,
			s.src[sprites.time/s.timeDiv], sprites.dst)
	}
}

func Load(width, height int32, renderer *sdl.Renderer) *Stage {
	rect4x4 := sdl.Rect{0, 0, 4, 4}
	stageRect := sdl.Rect{0, 0, width * size, height * size}
	offsetFromScreenX := (screenWidth - width*size) / 2
	offsetFromScreenY := (screenHeight - height*size) / 2
	stageScreenRect := sdl.Rect{offsetFromScreenX, offsetFromScreenY,
		screenWidth * size, screenHeight * size}

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

	tileStage := TileStage{false, &tileInfo, nil,
		tileTexture, &stageRect, &stageScreenRect,
		&sdl.Rect{}, width, height}

	spriteTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(width*size), int(height*size))
	e(err)

	spriteStage := SpriteStage{nil, spriteTexture,
		&stageRect, &stageScreenRect, 0}

	return &Stage{&tileStage, &spriteStage}
}
