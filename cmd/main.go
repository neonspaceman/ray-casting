package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
	"ray-casting/pkg/helpers"
	"ray-casting/pkg/vec"
)

const Wall = 1
const Width = 12
const Height = 12
const BlockSize = 64

type intersection struct {
	exists bool
	point  vec.Vec2
}

type Game struct {
	filed  [Height][Width]int
	player vec.Vec2
	pov    vec.Vec2
	angel  float32
	mouse  vec.Vec2
	fov    float64
	h      intersection
	v      intersection
	ray    vec.Vec2
}

func NewGame() *Game {
	return &Game{
		filed: [Height][Width]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		player: vec.NewVec2(float32((Width-1)*BlockSize)/2, float32((Height-1)*BlockSize)/2),
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

	g.h = intersection{}
	g.v = intersection{}

	// Horizontal intersection (iterable)
	{
		y := float32(0)

		if g.pov.Y > 0 {
			y = float32(int(g.player.Y/BlockSize))*BlockSize + BlockSize
		} else {
			y = float32(int(g.player.Y/BlockSize))*BlockSize - 1 // -1 is offset for detection correct block
		}

		x := g.player.X + (y-g.player.Y)/float32(math.Tan(float64(g.angel)))

		for x >= 0 && y >= 0 && x <= Width*BlockSize && y <= Height*BlockSize {
			col, row := int(x/BlockSize), int(y/BlockSize)
			if g.filed[row][col] == Wall {
				g.h.exists = true
				g.h.point = vec.NewVec2(x, y)
				break
			}

			dy := float32(0)

			if g.pov.Y > 0 {
				dy = BlockSize
			} else {
				dy = -BlockSize
			}

			dx := dy / float32(math.Tan(float64(g.angel)))

			x = x + dx
			y = y + dy
		}
	}

	// Vertical intersection
	{
		x := float32(0)

		if g.pov.X > 0 {
			x = float32(int(g.player.X/BlockSize))*BlockSize + BlockSize
		} else {
			x = float32(int(g.player.X/BlockSize))*BlockSize - 1 // -1 is offset for detection correct block
		}

		y := g.player.Y + (x-g.player.X)*float32(math.Tan(float64(g.angel)))

		for x >= 0 && y >= 0 && x <= Width*BlockSize && y <= Height*BlockSize {
			col, row := int(x/BlockSize), int(y/BlockSize)
			if g.filed[row][col] == Wall {
				g.v.exists = true
				g.v.point = vec.NewVec2(x, y)
				break
			}

			dx := float32(0)

			if g.pov.X > 0 {
				dx = BlockSize
			} else {
				dx = -BlockSize
			}

			dy := dx * float32(math.Tan(float64(g.angel)))

			x = x + dx
			y = y + dy

		}
	}

	switch {
	case g.h.exists && !g.v.exists:
		g.ray = g.h.point
	case !g.h.exists && g.v.exists:
		g.ray = g.v.point
	case g.h.exists && g.v.exists:
		if g.h.point.Sub(g.player).Len() < g.v.point.Sub(g.player).Len() {
			g.ray = g.h.point
		} else {
			g.ray = g.v.point
		}
	default:
		g.ray = vec.NewVec2(0, 0)
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
			if g.filed[row][col] == Wall {
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
	{
		p := g.ray.Mul(scale)
		vector.StrokeCircle(screen, p.X, p.Y, 5, 2, helpers.Color(0x0000FFFF), true)
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
