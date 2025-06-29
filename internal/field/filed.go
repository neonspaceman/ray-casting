package field

type BlockTypes int

const Wall BlockTypes = 1

type Field struct {
	Field     [][]BlockTypes
	Width     int
	Height    int
	BlockSize float32
}

func NewField(field [][]BlockTypes, width, height int, blockSize float32) Field {
	return Field{
		Field:     field,
		Width:     width,
		Height:    height,
		BlockSize: blockSize,
	}
}

func (f *Field) Coordinate(x, y float32) (col int, row int, err error) {
	if x < 0 || y < 0 || x >= float32(f.Width)*f.BlockSize || y >= float32(f.Height)*f.BlockSize {
		return 0, 0, ErrOutOfRange
	}

	return int(x / f.BlockSize), int(y / f.BlockSize), nil
}

func (f *Field) IsBlockType(col, row int, blockType BlockTypes) (bool, error) {
	if row < 0 || col < 0 || row >= f.Height || col >= f.Width {
		return false, ErrOutOfRange
	}

	return f.Field[row][col] == blockType, nil
}
