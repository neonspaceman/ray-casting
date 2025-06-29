package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"log"
	"math"
	"ray-casting/internal/field"
	"ray-casting/internal/ray_cast"
	"ray-casting/pkg/helpers"
	"ray-casting/pkg/vec"
)

const Width int = 12
const Height int = 12
const BlockSize float32 = 64
const ProjectionPlaneWidth int = 640
const ProjectionPlaceHeight int = 480
const FOV float32 = 60 * math.Pi / 180 // 60 degree in rad
var ProjectionDistance float32 = float32(ProjectionPlaneWidth) * .5 * float32(math.Tan(float64(FOV*.5)))

type rayCastResult struct {
	angel    float32
	distance float32
}

type Game struct {
	f        field.Field
	position vec.Vec2
	pov      vec.Vec2
	angel    float32
	mouse    vec.Vec2
	rays     []rayCastResult
}

func NewGame() *Game {
	f := field.NewField(
		[][]field.BlockTypes{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1},
			{1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		Width,
		Height,
		BlockSize,
	)

	return &Game{
		f:        f,
		position: vec.NewVec2((float32(Width-1)*BlockSize)*.5, (float32(Height-1)*BlockSize)*.5),
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.position = g.position.Add(vec.NewVec2(1, 0).Rot(g.angel).MulValue(2))
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.angel -= math.Pi / 40
		if g.angel < 0 {
			g.angel = 2 * math.Pi
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.angel += math.Pi / 40
		if g.angel > 2*math.Pi {
			g.angel = 0
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.position = g.position.Add(vec.NewVec2(1, 0).Rot(g.angel + math.Pi).MulValue(2))
	}

	//mouseX, mouseY := ebiten.CursorPosition()
	//g.mouse = vec.NewVec2(float32(mouseX), float32(mouseY))
	//g.angel = g.mouse.Sub(g.position).Rad()

	// Calculate rays
	g.rays = g.rays[:0]

	for a := -FOV * .5; a <= FOV*.5; a += FOV / float32(ProjectionPlaneWidth) {
		angel := g.angel + float32(a)
		distance := ray_cast.RayCast(g.f, g.position, angel)
		g.rays = append(g.rays, rayCastResult{
			angel:    angel,
			distance: distance,
		})
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%#v, %f", g.position, g.angel))

	// 2.5D
	h := vec.NewVec2(0, float32(ProjectionPlaceHeight)*.5)

	for i := range ProjectionPlaneWidth {
		ray := g.rays[i]
		projectHeight := BlockSize / ray.distance * ProjectionDistance

		from := h.Add(vec.NewVec2(float32(i), -projectHeight*.5))
		to := from.Add(vec.NewVec2(0, projectHeight))

		vector.StrokeLine(screen, from.X, from.Y, to.X, to.Y, 1, helpers.Color(0xFFFFFFFF), true)
		//from := g.position.Mul(scale)
		//to := g.position.Add(vec.NewVec2(1, 0).Rot(r.angel).MulValue(r.distance)).Mul(scale)
		//vector.StrokeLine(screen, from.X, from.Y, to.X, to.Y, 1, color.RGBA{255, 0, 0, 255}, true)
	}

	scale := vec.Vec2{.25, .25}
	offset := vec.NewVec2(float32(ProjectionPlaneWidth), float32(ProjectionPlaceHeight)).Sub(vec.NewVec2(BlockSize*float32(Width), BlockSize*float32(Height)).Mul(scale)).Sub(vec.NewVec2(20, 20))

	// MAP
	for row := range Height {
		for col := range Width {
			p := vec.NewVec2(float32(col), float32(row)).MulValue(BlockSize).Mul(scale).Add(offset)
			d := vec.NewVec2(BlockSize, BlockSize).Mul(scale)

			isWall, _ := g.f.IsBlockType(col, row, field.Wall)

			if isWall {
				vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, helpers.Color(0xFFFFFFFF), true)
			} else {
				vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, helpers.Color(0x000000FF), true)
				vector.StrokeRect(screen, p.X, p.Y, d.X, d.Y, 1, helpers.Color(0xFFFFFFFF), true)
			}
		}
	}

	// Player
	{
		p := g.position.Mul(scale).Add(offset)
		vector.DrawFilledCircle(screen, p.X, p.Y, 1, helpers.Color(0xFF0000FF), true)
		//pov := g.position.Add(g.pov.MulValue(100)).Mul(scale)
		//vector.StrokeLine(screen, p.X, p.Y, pov.X, pov.Y, 2, helpers.Color(0xFF0000FF), true)
	}

	// Rays
	for _, ray := range g.rays {
		p := g.position.Mul(scale).Add(offset)
		r := g.position.Add(vec.NewVec2(1, 0).Rot(ray.angel).MulValue(ray.distance)).Mul(scale).Add(offset)
		vector.StrokeLine(screen, p.X, p.Y, r.X, r.Y, 1, helpers.Color(0xFF0000FF), true)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ProjectionPlaneWidth, ProjectionPlaceHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ray-casting")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
