package ray_cast

import (
	"fmt"
	"log"
	"math"
	"ray-casting/internal/field"
	"ray-casting/pkg/vec"
)

type intersection struct {
	exists bool
	point  vec.Vec2
}

func RayCast(field field.Field, pos vec.Vec2, rad float32) (vec.Vec2, error) {
	h, v := horizontalRayCast(field, pos, rad), verticalRayCast(field, pos, rad)

	switch {
	case h.exists && !v.exists:
		return h.point, nil
	case !h.exists && v.exists:
		return v.point, nil
	case h.exists && v.exists:
		if h.point.Sub(pos).Len2() < v.point.Sub(pos).Len2() {
			return h.point, nil
		} else {
			return v.point, nil
		}
	}

	return vec.Vec2{}, ErrUnreachableRay
}

func horizontalRayCast(f field.Field, pos vec.Vec2, rad float32) intersection {
	ray := vec.NewVec2(1, 0).Rot(rad)

	y := float32(0)

	if ray.Y > 0 {
		y = float32(int(pos.Y/f.BlockSize))*f.BlockSize + f.BlockSize
	} else {
		y = float32(int(pos.Y/f.BlockSize))*f.BlockSize - 1 // -1 is offset for detection correct block
	}

	x := pos.X + (y-pos.Y)/float32(math.Tan(float64(rad)))

	for {
		col, row, err := f.Coordinate(x, y)

		if err != nil {
			break
		}

		isWall, err2 := f.IsBlockType(col, row, field.Wall)

		if err2 != nil {
			fmt.Println(x, y, col, row, err)
			log.Fatalln("Error while ray-casting", col, row, err, err2)
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

	if ray.X > 0 {
		x = float32(int(pos.X/f.BlockSize))*f.BlockSize + f.BlockSize
	} else {
		x = float32(int(pos.X/f.BlockSize))*f.BlockSize - 1 // -1 is offset for detection correct block
	}

	y := pos.Y + (x-pos.X)*float32(math.Tan(float64(rad)))

	for {
		col, row, err := f.Coordinate(x, y)

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
