package scene

import "ray-casting/pkg/vec"

type SpriteType int

type WallType int

const (
	None WallType = iota
)

type Sprite struct {
	pos vec.Vec2
	t   SpriteType
}

type Scene struct {
	Walls     [][]WallType
	Sprites   []Sprite
	Cols      int
	Rows      int
	Width     float32
	Height    float32
	BlockSize float32
}

func NewScene(walls [][]WallType, cols, rows int, blockSize float32, sprites []Sprite) Scene {
	return Scene{
		Walls:     walls,
		Sprites:   sprites,
		Cols:      cols,
		Rows:      rows,
		Width:     float32(cols) * blockSize,
		Height:    float32(rows) * blockSize,
		BlockSize: blockSize,
	}
}

func (f *Scene) WallType(x, y float32) (WallType, error) {
	if x < 0 || y < 0 || x >= f.Height || y >= f.Width {
		return 0, ErrOutOfRange
	}

	return f.Walls[int(y/f.BlockSize)][int(x/f.BlockSize)], nil
}
