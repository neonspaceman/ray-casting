package ray_cast

import (
	"log"
	"math"
	"ray-casting/internal/field"
	"ray-casting/pkg/vec"
)

type intersection struct {
	exists bool
	point  vec.Vec2
}

func RayCast(field field.Field, pos vec.Vec2, rad float32) float32 {
	h, v := horizontalRayCast(field, pos, rad), verticalRayCast(field, pos, rad)

	switch {
	case h.exists && !v.exists:
		return h.point.Sub(pos).Len()
	case !h.exists && v.exists:
		return v.point.Sub(pos).Len()
	case h.exists && v.exists:
		hLen, vLen := h.point.Sub(pos).Len2(), v.point.Sub(pos).Len2()
		if hLen < vLen {
			return float32(math.Sqrt(float64(hLen)))
		} else {
			return float32(math.Sqrt(float64(vLen)))
		}
	}

	return math.MaxFloat32
}

func horizontalRayCast(f field.Field, pos vec.Vec2, rad float32) intersection {
	ray := vec.NewVec2(1, 0).Rot(rad)

	y := float32(0)
	correction := float32(0)

	if ray.Y > 0 {
		y = float32(int(pos.Y/f.BlockSize))*f.BlockSize + f.BlockSize
	} else {
		// Due to x, y are used to detect col and row, we need offset 1, to correct detection col and row
		y = float32(int(pos.Y/f.BlockSize)) * f.BlockSize
		correction = 1
	}

	x := pos.X + (y-pos.Y)/float32(math.Tan(float64(rad)))

	for {
		col, row, err := f.Coordinate(x, y-correction)

		if err != nil {
			break
		}

		isWall, err := f.IsBlockType(col, row, field.Wall)

		if err != nil {
			log.Fatalln("Error while ray-casting", err)
		}

		if isWall {
			return intersection{exists: true, point: vec.NewVec2(x, y)}
		}

		dy := float32(0)

		if ray.Y > 0 {
			dy = f.BlockSize
		} else {
			dy = -f.BlockSize
		}

		dx := dy / float32(math.Tan(float64(rad)))

		x = x + dx
		y = y + dy
	}

	return intersection{exists: false}
}

func verticalRayCast(f field.Field, pos vec.Vec2, rad float32) intersection {
	ray := vec.NewVec2(1, 0).Rot(rad)

	x := float32(0)
	correction := float32(0)

	if ray.X > 0 {
		x = float32(int(pos.X/f.BlockSize))*f.BlockSize + f.BlockSize
	} else {
		// Due to x, y are used to detect col and row, we need offset 1, to correct detection col and row
		x = float32(int(pos.X/f.BlockSize)) * f.BlockSize
		correction = 1
	}

	y := pos.Y + (x-pos.X)*float32(math.Tan(float64(rad)))

	for {
		col, row, err := f.Coordinate(x-correction, y)

		if err != nil {
			break
		}

		isWall, err := f.IsBlockType(col, row, field.Wall)

		if err != nil {
			log.Fatalln("Error while ray-casting", col, row, err)
		}

		if isWall {
			return intersection{exists: true, point: vec.NewVec2(x, y)}
		}

		dx := float32(0)

		if ray.X > 0 {
			dx = f.BlockSize
		} else {
			dx = -f.BlockSize
		}

		dy := dx * float32(math.Tan(float64(rad)))

		x = x + dx
		y = y + dy
	}

	return intersection{exists: false}
}
