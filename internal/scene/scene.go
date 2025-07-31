package scene

import "ray-casting/pkg/vec"

const (
	None int = iota
)

type Sprite struct {
	Pos  vec.Vec2
	Type int
}

type Scene struct {
	Walls     [][]int
	Sprites   []Sprite
	Cols      int
	Rows      int
	Width     float32
	Height    float32
	BlockSize float32
}

func NewScene(walls [][]int, cols, rows int, blockSize float32, sprites []Sprite) Scene {
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

func (f *Scene) WallType(x, y float32) (int, error) {
	if x < 0 || y < 0 || x >= f.Height || y >= f.Width {
		return 0, ErrOutOfRange
	}

	return f.Walls[int(y/f.BlockSize)][int(x/f.BlockSize)], nil
}
