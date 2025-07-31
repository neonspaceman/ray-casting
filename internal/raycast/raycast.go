package raycast

import (
	"math"
	scenePkg "ray-casting/internal/scene"
	"ray-casting/pkg/vec"
)

func Cast(scene scenePkg.Scene, pos vec.Vec2, rad float32) (float32, int, int) {
	hPoint, hWallOffset, hWallType, hOk := horizontalRayCast(scene, pos, rad)
	vPoint, vWallOffset, vWallType, vOk := verticalRayCast(scene, pos, rad)

	switch {
	case hOk && !vOk:
		return hPoint.Sub(pos).Len(), hWallOffset, hWallType
	case !hOk && vOk:
		return vPoint.Sub(pos).Len(), vWallOffset, vWallType
	case hOk && vOk:
		hLen, vLen := hPoint.Sub(pos).Len(), vPoint.Sub(pos).Len()
		if hLen <= vLen {
			return hLen, hWallOffset, hWallType
		} else {
			return vLen, vWallOffset, vWallType
		}
	}

	return math.MaxFloat32, 0, 0
}

func horizontalRayCast(scene scenePkg.Scene, pos vec.Vec2, rad float32) (vec.Vec2, int, int, bool) {
	ray := vec.NewVec2(1, 0).Rot(rad)

	y := float32(0)
	correction := float32(0)

	if ray.Y > 0 {
		y = float32(int(pos.Y/scene.BlockSize))*scene.BlockSize + scene.BlockSize
	} else {
		// Due to x, y are used to detect col and row, we need offset 1, to correct detection col and row
		y = float32(int(pos.Y/scene.BlockSize)) * scene.BlockSize
		correction = 1
	}

	x := pos.X + (y-pos.Y)/float32(math.Tan(float64(rad)))

	for {
		wall, err := scene.WallType(x, y-correction)

		if err != nil {
			break
		}

		if wall != scenePkg.None {
			return vec.NewVec2(x, y), int(x) % int(scene.BlockSize), wall, true
		}

		dy := float32(0)

		if ray.Y > 0 {
			dy = scene.BlockSize
		} else {
			dy = -scene.BlockSize
		}

		dx := dy / float32(math.Tan(float64(rad)))

		x = x + dx
		y = y + dy
	}

	return vec.Vec2{}, 0, 0, false
}

func verticalRayCast(scene scenePkg.Scene, pos vec.Vec2, rad float32) (vec.Vec2, int, int, bool) {
	ray := vec.NewVec2(1, 0).Rot(rad)

	x := float32(0)
	correction := float32(0)

	if ray.X > 0 {
		x = float32(int(pos.X/scene.BlockSize))*scene.BlockSize + scene.BlockSize
	} else {
		// Due to x, y are used to detect col and row, we need offset 1, to correct detection col and row
		x = float32(int(pos.X/scene.BlockSize)) * scene.BlockSize
		correction = 1
	}

	y := pos.Y + (x-pos.X)*float32(math.Tan(float64(rad)))

	for {
		wall, err := scene.WallType(x-correction, y)

		if err != nil {
			break
		}

		if wall != scenePkg.None {
			return vec.NewVec2(x, y), int(scene.BlockSize) - int(y)%int(scene.BlockSize) - 1, wall, true
		}

		dx := float32(0)

		if ray.X > 0 {
			dx = scene.BlockSize
		} else {
			dx = -scene.BlockSize
		}

		dy := dx * float32(math.Tan(float64(rad)))

		x = x + dx
		y = y + dy
	}

	return vec.Vec2{}, 0, 0, false
}
