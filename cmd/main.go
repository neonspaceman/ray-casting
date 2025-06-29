package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
	"ray-casting/internal/field"
	"ray-casting/internal/ray_cast"
	"ray-casting/pkg/helpers"
	"ray-casting/pkg/vec"
)

const Width = 12
const Height = 12
const BlockSize = 64

type intersection struct {
	exists bool
	point  vec.Vec2
}

type Game struct {
	f      field.Field
	player vec.Vec2
	pov    vec.Vec2
	angel  float32
	mouse  vec.Vec2
	fov    float32
	h      intersection
	v      intersection
	rays   []vec.Vec2
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
		f:      f,
		player: vec.NewVec2(float32((Width-1)*BlockSize)*.5, float32((Height-1)*BlockSize)*.5),
		fov:    math.Pi * 60 / 180, // 60 degrees
	}
}

func (g *Game) Update() error {
	//if ebiten.IsKeyPressed(ebiten.KeyUp) {
	//g.player = g.player.Add(vec.NewVec2(float32(math.Cos(g.angel)*2), float32(math.Sin(g.angel)*2)))
	//}

	//if ebiten.IsKeyPressed(ebiten.KeyDown) {
	//	g.player = g.player.Sub(vec.NewVec2(float32(math.Cos(g.angel)*2), float32(math.Sin(g.angel)*2)))
	//	//g.player = g.player.Add(vec.NewVec2(float32(math.Cos(g.angel+math.Pi)*2), float32(math.Sin(g.angel+math.Pi)*2)))
	//}

	mouseX, mouseY := ebiten.CursorPosition()

	g.mouse = vec.NewVec2(float32(mouseX), float32(mouseY))
	g.angel = g.mouse.Sub(g.player).Rad()

	g.pov = vec.NewVec2(1, 0).Rot(g.angel)

	// Calculate rays
	g.rays = g.rays[:0]

	a := -g.fov * 0.5
	da := g.fov / 640

	for a <= g.fov*0.5 {
		ray, _ := ray_cast.RayCast(g.f, g.player, g.angel+a)
		a += da
		g.rays = append(g.rays, ray)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scale := vec.Vec2{0.5, 0.5}

	//ebitenutil.DebugPrint(screen, fmt.Sprintf("%#v", g.pov))

	// MAP
	for row := range Height {
		for col := range Width {
			p := vec.NewVec2(float32(col), float32(row)).MulValue(BlockSize).Mul(scale)
			d := vec.NewVec2(BlockSize, BlockSize).Mul(scale)

			isWall, _ := g.f.IsBlockType(col, row, field.Wall)

			if isWall {
				vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, helpers.Color(0xFFFFFFFF), true)
			} else {
				vector.StrokeRect(screen, p.X, p.Y, d.X, d.Y, 1, helpers.Color(0xFFFFFFFF), true)
			}
		}
	}

	// Player
	{
		p := g.player.Mul(scale)
		vector.DrawFilledCircle(screen, p.X, p.Y, 5, helpers.Color(0xFF0000FF), true)
		pov := g.player.Add(g.pov.MulValue(100)).Mul(scale)
		vector.StrokeLine(screen, p.X, p.Y, pov.X, pov.Y, 2, color.RGBA{255, 0, 0, 255}, true)
	}

	// Intersection
	if g.h.exists {
		p := g.h.point.Mul(scale)
		vector.DrawFilledCircle(screen, p.X, p.Y, 5, helpers.Color(0xFF0000FF), true)
	}
	if g.v.exists {
		p := g.v.point.Mul(scale)
		vector.DrawFilledCircle(screen, p.X, p.Y, 5, helpers.Color(0x00FF00FF), true)
	}

	// Ray
	for _, r := range g.rays {
		p := g.player.Mul(scale)
		r = r.Mul(scale)
		vector.StrokeLine(screen, p.X, p.Y, r.X, r.Y, 1, color.RGBA{255, 0, 0, 255}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ray-casting")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
