package camera

import "math"

type Camera struct {
	ProjectionPlaneWidth    float32
	ProjectionPlaneHeight   float32
	ProjectionPlaneDistance float32
	FOV                     float32
	Z                       float32
}

func NewCamera(projectionPlaneWidth, projectionPlaneHeight float32, fovInDegrees int, z float32) Camera {
	fov := float32(fovInDegrees) * math.Pi / 180

	return Camera{
		ProjectionPlaneWidth:    projectionPlaneWidth,
		ProjectionPlaneHeight:   projectionPlaneHeight,
		ProjectionPlaneDistance: projectionPlaneWidth * .5 * float32(math.Tan(float64(fov))),
		FOV:                     fov,
		Z:                       z,
	}
}
