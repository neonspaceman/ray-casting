package raycast

import (
	"math"
	"ray-casting/internal/scene"
	"ray-casting/pkg/vec"
)

func Cast(scene scene.Scene, pos vec.Vec2, rad float32) float32 {
	hPoint, hOk := horizontalRayCast(scene, pos, rad)
	vPoint, vOk := verticalRayCast(scene, pos, rad)

	switch {
	case hOk && !vOk:
		return hPoint.Sub(pos).Len()
	case !hOk && vOk:
		return vPoint.Sub(pos).Len()
	case hOk && vOk:
		return float32(math.Min(float64(hPoint.Sub(pos).Len()), float64(vPoint.Sub(pos).Len())))
	}

	return math.MaxFloat32
}

func horizontalRayCast(f scene.Scene, pos vec.Vec2, rad float32) (vec.Vec2, bool) {
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
		wall, err := f.WallType(x, y-correction)

		if err != nil {
			break
		}

		if wall != scene.None {
			return vec.NewVec2(x, y), true
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

	return vec.Vec2{}, false
}

func verticalRayCast(f scene.Scene, pos vec.Vec2, rad float32) (vec.Vec2, bool) {
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
		wall, err := f.WallType(x-correction, y)

		if err != nil {
			break
		}

		if wall != scene.None {
			return vec.NewVec2(x, y), true
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

	return vec.Vec2{}, false
}
