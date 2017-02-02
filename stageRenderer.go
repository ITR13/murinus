package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Tile uint8

const (
	Empty     Tile = iota
	Wall      Tile = iota
	Point     Tile = iota
	Powerup   Tile = iota
	p200      Tile = iota
	p500      Tile = iota
	p1000     Tile = iota
	p2000     Tile = iota
	SnakeWall Tile = iota

	NumTiles int = iota
)

type SpriteID uint8

const (
	Player1        SpriteID = iota
	Player2        SpriteID = iota
	SnakeHead      SpriteID = iota
	SnakeBody      SpriteID = iota
	SnakeHeadSig   SpriteID = iota
	SnakeBodySig   SpriteID = iota
	SnakeHeadGhost SpriteID = iota
	SnakeBodyGhost SpriteID = iota

	NumSprites int = iota
)

type Stage struct {
	input          *Input
	tiles          *TileStage
	sprites        *SpriteStage
	scoreField     *ScoreField
	stages         []*PreStageData
	levels         [3][][2]int
	pointsLeft, ID int
	hideWalls      bool
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
	sprites   SpriteInfo
	entities  [][]*Entity
	texture   *sdl.Texture
	src, dst  *sdl.Rect
	spriteDst *sdl.Rect
	time      int64
}

type Entity struct {
	sprite    *Sprite
	x, y      int32
	precision int32
	display   bool
	dir       Direction
}

type Sprite struct {
	texture  *sdl.Texture
	src      []*sdl.Rect
	timeDiv  int64
	priority int
}

type SpriteInfo []*Sprite

type ScoreField struct {
	rect           *sdl.Rect
	lives          *sdl.Rect
	xOffset, xMult int32
}

func (spriteStage *SpriteStage) GetEntity(x, y int32, id SpriteID) *Entity {
	entity := &Entity{
		spriteStage.sprites[id],
		x, y, 0, true, Right}
	l := spriteStage.sprites[id].priority
	spriteStage.entities[l] = append(spriteStage.entities[l], entity)
	return entity
}

func (spriteStage *SpriteStage) GetSnake(x, y int32, length int, ai AI,
	moveTimerMax, growTimerMax,
	minLength, maxLength int) *Snake {
	length += 2
	entities := make([]*Entity, length)
	entities[0] = spriteStage.GetEntity(x, y, SnakeHead)
	for i := 1; i < length; i++ {
		entities[i] = spriteStage.GetEntity(x, y, SnakeBody)
	}
	return &Snake{
		entities[0], entities[1 : length-1], entities[length-1], ai, false,
		moveTimerMax / 2, moveTimerMax,
		growTimerMax / 3, growTimerMax, 0, minLength, length - 2, maxLength}
}

func (stage *Stage) Render(renderer *sdl.Renderer,
	lives, score int32, hideWalls bool) {
	defer renderer.Present()
	renderedOnce := stage.tiles.renderedOnce

	if redrawTextures {
		stage.tiles.tileInfo.Draw(renderer)
		stage.sprites.sprites.Draw(renderer)
		renderedOnce = false
	}

	if newScreenWidth != screenWidth || newScreenHeight != screenHeight {
		SetWindowSize(newScreenWidth, newScreenHeight, stage)
		renderedOnce = false
	}

	if !renderedOnce {
		stage.tiles.Render(renderer, hideWalls && stage.hideWalls)
	}
	stage.sprites.Render(renderer)
	renderer.SetRenderTarget(nil)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Copy(stage.tiles.texture, stage.tiles.src, stage.tiles.dst)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.Copy(stage.sprites.texture, stage.sprites.src, stage.sprites.dst)
	for i := int32(0); i < lives; i++ {
		renderer.SetDrawColor(255, 182, 193, 255)
		stage.scoreField.lives.X = stage.scoreField.xOffset +
			stage.scoreField.xMult*i
		renderer.FillRect(stage.scoreField.lives)
	}

}

func (tiles *TileStage) Render(renderer *sdl.Renderer, hideWalls bool) {
	e(renderer.SetRenderTarget(tiles.texture))
	e(renderer.SetDrawBlendMode(sdl.BLENDMODE_NONE))
	e(renderer.SetDrawColor(0, 0, 0, 255))
	e(renderer.Clear())
	if tiles.tiles != nil {
		for x := int32(0); x < tiles.w; x++ {
			tiles.tileDst.X = x * gSize
			for y := int32(0); y < tiles.h; y++ {
				tiles.tileDst.Y = y * gSize
				tiles.tileDst.W = gSize
				tiles.tileDst.H = gSize
				tile := tiles.tiles[x][y]
				if hideWalls && tile == Wall {
					tile = Empty
				}

				e(renderer.Copy(tiles.tileInfo.textures[tile],
					tiles.tileInfo.src[tile], tiles.tileDst))
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
		for priority := 0; priority < len(sprites.entities); priority++ {
			for i := 0; i < len(sprites.entities[priority]); i++ {
				e := sprites.entities[priority][i]
				s := e.sprite
				if s.priority == priority && e.display {
					sprites.spriteDst.X = e.x * gSize
					sprites.spriteDst.Y = e.y * gSize
					if e.dir == Up {
						sprites.spriteDst.Y -= e.precision * gSize / PrecisionMax
					} else if e.dir == Right {
						sprites.spriteDst.X += e.precision * gSize / PrecisionMax
					} else if e.dir == Down {
						sprites.spriteDst.Y += e.precision * gSize / PrecisionMax
					} else if e.dir == Left {
						sprites.spriteDst.X -= e.precision * gSize / PrecisionMax
					}
					sprites.spriteDst.W = gSize
					sprites.spriteDst.H = gSize
					t := (sprites.time / s.timeDiv) % int64(len(s.src))
					renderer.Copy(s.texture,
						s.src[t], sprites.spriteDst)
				}
			}
		}
	}
}

//TODO Figure out if it's better to save these variables in Stage
var rect8x8, rect6x6, rect4x4 sdl.Rect

func LoadTextures(renderer *sdl.Renderer, input *Input) *Stage {
	w, h := stageWidth, stageHeight
	rect8x8 = sdl.Rect{0, 0, gSize, gSize}
	rect6x6 = sdl.Rect{gSize/4 - 1, gSize/4 - 1, gSize/2 + 2, gSize/2 + 2}
	rect4x4 = sdl.Rect{gSize / 4, gSize / 4, gSize / 2, gSize / 2}
	stageRect := sdl.Rect{0, 0, w * gSize, h * gSize}
	offsetFromScreenX := (screenWidth - w*blockSize) / 2
	offsetFromScreenY := (screenHeight - h*(blockSize+2)) / 2
	stageScreenRect := sdl.Rect{offsetFromScreenX, blockSize + offsetFromScreenY,
		w * blockSize, h * blockSize}

	tileTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(w*gSize), int(h*gSize))
	e(err)
	tileInfo := TileInfo{&sdl.Rect{},
		make([]*sdl.Texture, NumTiles),
		make([]*sdl.Rect, NumTiles)}
	for i := 0; i < NumTiles; i++ {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
			sdl.TEXTUREACCESS_TARGET, int(gSize), int(gSize))
		e(err)
		tileInfo.textures[i] = texture
		tileInfo.src[i] = &rect8x8
	}

	tileInfo.Draw(renderer)
	tileStage := TileStage{false, &tileInfo, nil,
		tileTexture, &stageRect, &stageScreenRect,
		&sdl.Rect{}, w, h}

	spriteTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(w*gSize), int(h*gSize))
	e(err)
	spriteTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	spriteInfo := make(SpriteInfo, NumSprites)
	for i := 0; i < len(spriteInfo); i++ {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
			sdl.TEXTUREACCESS_TARGET, int(gSize), int(gSize))
		texture.SetBlendMode(sdl.BLENDMODE_BLEND)
		e(err)
		spriteInfo[i] = &Sprite{texture, []*sdl.Rect{&rect8x8}, 1, 0}
	}

	spriteInfo.Draw(renderer)
	spriteStage := SpriteStage{spriteInfo, nil, spriteTexture,
		&stageRect, &stageScreenRect, &sdl.Rect{}, 0}

	scoreField := ScoreField{&sdl.Rect{offsetFromScreenX, offsetFromScreenY,
		w * blockSize, blockSize}, &sdl.Rect{0, 4 + offsetFromScreenY,
		blockSize - 8, blockSize - 8}, 4 + offsetFromScreenX, blockSize}

	data, levels := GetPreStageDatas()
	return &Stage{input, &tileStage, &spriteStage,
		&scoreField, data, levels, -1, -1, false}
}

func (tileInfo *TileInfo) Draw(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(tileInfo.textures[Empty])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetRenderTarget(tileInfo.textures[Wall])
	renderer.SetDrawColor(25, 25, 112, 255)
	renderer.Clear()
	renderer.SetRenderTarget(tileInfo.textures[Point])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(240, 230, 140, 255)
	renderer.FillRect(&rect4x4)
	renderer.SetRenderTarget(tileInfo.textures[Powerup])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(255, 165, 0, 255)
	renderer.FillRect(&rect4x4)

	renderer.SetRenderTarget(tileInfo.textures[p200])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(255, 255, 51, 255)
	renderer.FillRect(&rect4x4)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.DrawRect(&rect6x6)

	renderer.SetRenderTarget(tileInfo.textures[p500])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(&rect4x4)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.DrawRect(&rect6x6)

	renderer.SetRenderTarget(tileInfo.textures[p1000])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(0, 0, 205, 255)
	renderer.FillRect(&rect4x4)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.DrawRect(&rect6x6)

	renderer.SetRenderTarget(tileInfo.textures[p2000])
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.FillRect(&rect6x6)

	renderer.SetRenderTarget(tileInfo.textures[SnakeWall])
	//renderer.SetDrawColor(45, 45, 132, 255)
	renderer.SetDrawColor(65, 105, 225, 255)
	renderer.Clear()
	renderer.SetDrawColor(80, 120, 240, 255)
	renderer.FillRect(&rect6x6)
	renderer.SetDrawColor(95, 135, 255, 255)
	renderer.FillRect(&rect4x4)

}

func (spriteInfo SpriteInfo) Draw(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(spriteInfo[Player1].texture)
	renderer.SetDrawColor(216, 75, 139, 255)
	renderer.Clear()
	renderer.SetDrawColor(255, 182, 193, 255)
	renderer.FillRect(&rect6x6)
	spriteInfo[Player1].priority = 1

	renderer.SetRenderTarget(spriteInfo[SnakeHead].texture)
	renderer.SetDrawColor(0, 95, 0, 255)
	renderer.Clear()
	spriteInfo[SnakeHead].priority = 2

	renderer.SetRenderTarget(spriteInfo[SnakeBody].texture)
	renderer.SetDrawColor(0, 127, 0, 255)
	renderer.Clear()
	spriteInfo[SnakeBody].priority = 0

	renderer.SetRenderTarget(spriteInfo[SnakeHeadSig].texture)
	renderer.SetDrawColor(127, 0, 0, 255)
	renderer.Clear()
	spriteInfo[SnakeHeadSig].priority = 2

	renderer.SetRenderTarget(spriteInfo[SnakeBodySig].texture)
	renderer.SetDrawColor(159, 0, 0, 255)
	renderer.Clear()
	spriteInfo[SnakeBodySig].priority = 0

	renderer.SetRenderTarget(spriteInfo[SnakeHeadGhost].texture)
	renderer.SetDrawColor(255, 255, 255, 235)
	renderer.Clear()
	spriteInfo[SnakeHeadGhost].priority = 2

	renderer.SetRenderTarget(spriteInfo[SnakeBodyGhost].texture)
	renderer.SetDrawColor(235, 235, 235, 220)
	renderer.Clear()
	spriteInfo[SnakeBodyGhost].priority = 0

}

func (stage *Stage) Free() {
	stage.scoreField = nil

	tileStage := stage.tiles
	tileStage.texture.Destroy()
	for i := 0; i < len(tileStage.tileInfo.textures); i++ {
		tileStage.tileInfo.textures[i].Destroy()
	}
	stage.tiles = nil

	spriteStage := stage.sprites
	spriteStage.texture.Destroy()
	for i := 0; i < len(spriteStage.sprites); i++ {
		spriteStage.sprites[i].texture.Destroy()
	}
	stage.sprites = nil
}

func SetWindowSize(w, h int32, stage *Stage) {
	maxSize := w * 5
	if maxSize > h*8 {
		maxSize = h * 8
	}
	screenWidth, screenHeight = w, h
	lcm := (1280 * 5 / gcd(maxSize, 1280*5)) * maxSize
	div := lcm / maxSize
	mult := lcm / (1280 * 5)
	blockSize = (48 * mult) / div
	blockSizeBigBoard = (24 * mult) / div

	offsetFromScreenX := (screenWidth - stageWidth*blockSize) / 2
	offsetFromScreenY := (screenHeight - stageHeight*(blockSize+2)) / 2
	stage.tiles.dst.X, stage.tiles.dst.Y, stage.tiles.dst.W, stage.tiles.dst.H =
		offsetFromScreenX, blockSize+offsetFromScreenY,
		stageWidth*blockSize, stageHeight*blockSize
}

//From rosettacode.org
func gcd(x, y int32) int32 {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
