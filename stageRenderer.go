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
)

type SpriteID uint8

const (
	Player1   SpriteID = iota
	Player2   SpriteID = iota
	SnakeHead SpriteID = iota
	SnakeBody SpriteID = iota
	SnakeTail SpriteID = iota
)

type Stage struct {
	tiles          *TileStage
	sprites        *SpriteStage
	scoreField     *ScoreField
	pointsLeft, ID int
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
	sprite    *Sprite
	x, y      int32
	precision int32 //Consider making uint8
	display   bool
	dir       Direction
}

type Sprite struct {
	texture  *sdl.Texture
	src      []*sdl.Rect
	timeDiv  int64
	priority int
}

type ScoreField struct {
	rect           *sdl.Rect
	lives          *sdl.Rect
	xOffset, xMult int32
}

func (spriteStage *SpriteStage) GetEntity(x, y int32, id SpriteID) *Entity {
	entity := &Entity{
		spriteStage.sprites[id],
		x, y, 0, true, Right}
	spriteStage.entities = append(spriteStage.entities, entity)
	return entity
}

func (spriteStage *SpriteStage) GetSnake(x, y int32, length int, ai AI,
	moveTimerMax, speedUpTimerMax, growTimerMax,
	minLength, maxLength int) *Snake {
	length += 2
	entities := make([]*Entity, length)
	entities[0] = spriteStage.GetEntity(x, y, SnakeHead)
	for i := 1; i < length; i++ {
		entities[i] = spriteStage.GetEntity(x, y, SnakeBody)
	}
	return &Snake{
		entities[0], entities[1 : length-1], entities[length-1], ai, false,
		moveTimerMax / 2, moveTimerMax, speedUpTimerMax, speedUpTimerMax,
		growTimerMax / 3, growTimerMax, 0, minLength, length - 2, maxLength}
}

func (stage *Stage) Render(renderer *sdl.Renderer,
	lives, score int32) {
	if !stage.tiles.renderedOnce {
		stage.tiles.Render(renderer)
	}
	stage.sprites.Render(renderer)
	renderer.SetRenderTarget(nil)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.Clear()
	renderer.Copy(stage.tiles.texture, stage.tiles.src, stage.tiles.dst)
	renderer.Copy(stage.sprites.texture, stage.sprites.src, stage.sprites.dst)

	for i := int32(0); i < lives; i++ {
		renderer.SetDrawColor(255, 182, 193, 255)
		stage.scoreField.lives.X = stage.scoreField.xOffset +
			stage.scoreField.xMult*i
		renderer.FillRect(stage.scoreField.lives)
	}

	renderer.Present()
}

func (tiles *TileStage) Render(renderer *sdl.Renderer) {
	renderer.SetRenderTarget(tiles.texture)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	if tiles.tiles != nil {
		for x := int32(0); x < tiles.w; x++ {
			tiles.tileDst.X = x * blockSize
			for y := int32(0); y < tiles.h; y++ {
				tiles.tileDst.Y = y * blockSize
				tiles.tileDst.W = blockSize
				tiles.tileDst.H = blockSize
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
				if s.priority == priority && e.display {
					sprites.spriteDst.X = e.x * blockSize
					sprites.spriteDst.Y = e.y * blockSize
					if e.dir == Up {
						sprites.spriteDst.Y -= e.precision * blockSize / PrecisionMax
					} else if e.dir == Right {
						sprites.spriteDst.X += e.precision * blockSize / PrecisionMax
					} else if e.dir == Down {
						sprites.spriteDst.Y += e.precision * blockSize / PrecisionMax
					} else if e.dir == Left {
						sprites.spriteDst.X -= e.precision * blockSize / PrecisionMax
					}
					sprites.spriteDst.W = blockSize
					sprites.spriteDst.H = blockSize
					t := (sprites.time / s.timeDiv) % int64(len(s.src))
					renderer.Copy(s.texture,
						s.src[t], sprites.spriteDst)
				}
			}
		}
	}
}

func LoadTextures(width, height int32, renderer *sdl.Renderer) *Stage {
	gSize := int32(12)
	rect8x8 := sdl.Rect{0, 0, gSize, gSize}
	rect6x6 := sdl.Rect{gSize/4 - 1, gSize/4 - 1, gSize/2 + 2, gSize/2 + 2}
	rect4x4 := sdl.Rect{gSize / 4, gSize / 4, gSize / 2, gSize / 2}
	stageRect := sdl.Rect{0, 0, width * blockSize, height * blockSize}
	offsetFromScreenX := (screenWidth - width*blockSize) / 2
	offsetFromScreenY := (screenHeight - height*(blockSize+2)) / 2
	stageScreenRect := sdl.Rect{offsetFromScreenX, blockSize + offsetFromScreenY,
		width * blockSize, height * blockSize}

	tileTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(width*blockSize), int(height*blockSize))
	e(err)
	tileInfo := TileInfo{&sdl.Rect{},
		make([]*sdl.Texture, SnakeWall+1),
		make([]*sdl.Rect, SnakeWall+1)}
	for i := Empty; i < SnakeWall; i++ {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
			sdl.TEXTUREACCESS_TARGET, 16, 16)
		e(err)
		tileInfo.textures[i] = texture
		tileInfo.src[i] = &rect8x8
	}
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

	tileInfo.textures[SnakeWall] = tileInfo.textures[Empty]

	tileStage := TileStage{false, &tileInfo, nil,
		tileTexture, &stageRect, &stageScreenRect,
		&sdl.Rect{}, width, height}

	spriteTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
		sdl.TEXTUREACCESS_TARGET, int(width*blockSize), int(height*blockSize))
	e(err)
	spriteTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	spriteDatas := make([]*Sprite, 8)
	for i := 0; i < len(spriteDatas); i++ {
		texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB565,
			sdl.TEXTUREACCESS_TARGET, 4, 4)
		texture.SetBlendMode(sdl.BLENDMODE_BLEND)
		e(err)
		spriteDatas[i] = &Sprite{texture, []*sdl.Rect{&rect8x8}, 1, 0}
	}

	renderer.SetRenderTarget(spriteDatas[Player1].texture)
	renderer.SetDrawColor(255, 182, 193, 255)
	renderer.Clear()
	spriteDatas[Player1].priority = 5

	renderer.SetRenderTarget(spriteDatas[SnakeHead].texture)
	renderer.SetDrawColor(0, 95, 0, 255)
	renderer.Clear()
	spriteDatas[SnakeHead].priority = 6

	renderer.SetRenderTarget(spriteDatas[SnakeBody].texture)
	renderer.SetDrawColor(0, 127, 0, 255)
	renderer.Clear()
	spriteDatas[SnakeBody].priority = 4

	spriteStage := SpriteStage{spriteDatas, nil, spriteTexture,
		&stageRect, &stageScreenRect, &sdl.Rect{}, 0}

	scoreField := ScoreField{&sdl.Rect{offsetFromScreenX, offsetFromScreenY,
		width * blockSize, blockSize}, &sdl.Rect{0, 4 + offsetFromScreenY,
		blockSize - 8, blockSize - 8}, 4 + offsetFromScreenX, blockSize}

	return &Stage{&tileStage, &spriteStage, &scoreField, -1, -1}
}
