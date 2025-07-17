package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
	"ray-casting/internal/camera"
	"ray-casting/internal/player"
	"ray-casting/internal/raycast"
	"ray-casting/internal/scene"
	"ray-casting/pkg/helpers"
	"ray-casting/pkg/vec"
)

const sceneRows int = 24
const sceneCols int = 24
const blockSize float32 = 64
const projectionPlaneWidth float32 = 320
const projectionPlaneHeight float32 = 240
const playerFOV int = 60
const playerHeight float32 = 32

var wallTextures = map[scene.WallType]string{
	1: "assets/Bark.png",
	2: "assets/WalkStone.png",
	3: "assets/WalkStone.png",
	4: "assets/WalkStone.png",
}

var spriteTextures = map[int]string{
	1: "assets/Coin.png",
}

type Game struct {
	scene        scene.Scene
	player       player.Player
	camera       camera.Camera
	wallTextures map[scene.WallType]*ebiten.Image
	pixels       []byte
	zBuffer      []int
}

func NewGame() *Game {
	return &Game{
		scene: scene.NewScene(
			[][]scene.WallType{
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			},
			sceneRows,
			sceneCols,
			blockSize,
			[]scene.Sprite{},
		),
		player: player.NewPlayer(
			vec.NewVec2((float32(sceneRows-1)*blockSize)*.5, (float32(sceneCols-1)*blockSize)*.5),
			0,
		),
		camera: camera.NewCamera(
			projectionPlaneWidth,
			projectionPlaneHeight,
			playerFOV,
			playerHeight,
		),
		wallTextures: make(map[scene.WallType]*ebiten.Image),
		pixels:       make([]byte, int(projectionPlaneWidth)*int(projectionPlaneHeight)*4),
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Up(2)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Left(math.Pi / 40)
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Right(math.Pi / 40)
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Down(2)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.clearPixels()
	g.renderFloor()
	g.renderCell()
	g.renderWall()

	screen.WritePixels(g.pixels)

	//opts := ebiten.DrawImageOptions{}
	//screen.DrawImage(g.wallTextures[scene.WallType(5)], &opts)

	//scale := float32(.1)
	//offset := vec.NewVec2(g.camera.ProjectionPlaneWidth, g.camera.ProjectionPlaneHeight).Sub(vec.NewVec2(g.scene.Width, g.scene.Height).MulValue(scale)).Sub(vec.NewVec2(20, 20))
	//g.renderMinimap(screen, offset, scale)
	//
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%#v, %f\nFPS: %f", g.player.Pos, g.player.Angel, ebiten.ActualFPS()))
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(g.camera.ProjectionPlaneWidth), int(g.camera.ProjectionPlaneHeight)
}

func (g *Game) LoadWallTexture(wallType scene.WallType, path string) error {
	t, _, err := ebitenutil.NewImageFromFile(path)

	if err != nil {
		return err
	}

	g.wallTextures[wallType] = t

	return nil
}

func (g *Game) renderWall() {
	h := g.camera.ProjectionPlaneHeight * .5

	a := -g.camera.FOV * .5

	for x := float32(0); x < g.camera.ProjectionPlaneWidth; x++ {
		distance, wallOffset, wall := raycast.Cast(g.scene, g.player.Pos, g.player.Angel+a)

		brightness := float32(160 / distance)
		sliceHeight := g.scene.BlockSize / (distance * float32(math.Cos(float64(a)))) * g.camera.ProjectionPlaneDistance

		if texture, ok := g.wallTextures[wall]; ok {
			ratio := g.scene.BlockSize / sliceHeight
			ty := float32(0)

			for y := h - sliceHeight*.5; y <= h+sliceHeight*.5; y += 1 {
				clr := texture.At(wallOffset, int(ty)).(color.RGBA)

				g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))

				ty += ratio
			}
		} else {
			for y := h - sliceHeight*.5; y <= h+sliceHeight*.5; y += 1 {
				g.setPixel(x, y, 255*brightness, 255*brightness, 255*brightness, 255)
			}
		}

		a += g.camera.FOV / g.camera.ProjectionPlaneWidth
	}

}

func (g *Game) renderMinimap(screen *ebiten.Image, offset vec.Vec2, scale float32) {
	// MAP
	{
		for row := range g.scene.Rows {
			for col := range g.scene.Cols {
				p := vec.NewVec2(float32(col), float32(row)).MulValue(blockSize).MulValue(scale).Add(offset)
				d := vec.NewVec2(blockSize, blockSize).MulValue(scale)

				wallType := g.scene.Walls[row][col]

				if wallType != scene.None {
					vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, helpers.Color(0xFFFFFFFF), true)
				} else {
					vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, helpers.Color(0x000000FF), true)
				}
			}
		}
	}

	//Player
	{
		from := g.player.Pos.MulValue(scale).Add(offset)

		to1 := g.player.Pos.Add(vec.NewRotated(g.player.Angel - g.camera.FOV*.5).MulValue(300)).MulValue(scale).Add(offset)
		to2 := g.player.Pos.Add(vec.NewRotated(g.player.Angel + g.camera.FOV*.5).MulValue(300)).MulValue(scale).Add(offset)

		vector.DrawFilledCircle(screen, from.X, from.Y, 3, helpers.Color(0xFF0000FF), true)
		vector.StrokeLine(screen, from.X, from.Y, to1.X, to1.Y, 2, helpers.Color(0xFF0000FF), true)
		vector.StrokeLine(screen, from.X, from.Y, to2.X, to2.Y, 2, helpers.Color(0xFF0000FF), true)
	}
}

// clampColor ensures a color value stays within 0-255 range
func clampColor(value float32) uint8 {
	return uint8(math.Max(0, math.Min(255, math.Round(float64(value)))))
}

func (g *Game) renderFloor() {
	for y := g.camera.ProjectionPlaneHeight*.5 + 1; y < g.camera.ProjectionPlaneHeight; y += 1 {
		a := -g.camera.FOV * .5
		ratio := g.camera.Z / (y - g.camera.ProjectionPlaneHeight*.5)

		for x := float32(0); x < g.camera.ProjectionPlaneWidth; x += 1 {
			distance := g.camera.ProjectionPlaneDistance * ratio / float32(math.Cos(float64(a)))

			t := g.player.Pos.Add(vec.NewRotated(g.player.Angel + a).MulValue(distance))

			cellX := t.X / g.scene.BlockSize
			cellY := t.Y / g.scene.BlockSize
			brightness := float32(100 / distance)

			//clr := g.wallTextures[scene.WallType(1)].At(int(t.X)%int(g.scene.BlockSize), int(t.Y)%int(g.scene.BlockSize)).(color.RGBA)
			//g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))

			if (int(cellX)+int(cellY))%2 == 0 {
				g.setPixel(x, y, 255*brightness, 0, 0, 255)
				//screen.DrawImage()
				//vector.DrawFilledRect(screen, x, y, 1, 1, ScaleRGBA(color.RGBA{0, 0, 255, 255}, brightness, brightness, brightness, 1), true)
			} else {
				g.setPixel(x, y, 0, 0, 255*brightness, 255)
				//vector.DrawFilledRect(screen, x, y, 1, 1, ScaleRGBA(color.RGBA{255, 0, 0, 255}, brightness, brightness, brightness, 1), true)
			}

			a += g.camera.FOV / g.camera.ProjectionPlaneWidth
		}
	}
	// for y := g.camera.ProjectionPlaneHeight/2 + 1; y < g.camera.ProjectionPlaneHeight; y += 1 {
}

func (g *Game) renderCell() {
	for y := float32(0); y <= g.camera.ProjectionPlaneHeight*.5; y += 1 {
		a := -g.camera.FOV * .5
		ratio := g.camera.Z / (g.camera.ProjectionPlaneHeight*.5 - y)

		for x := float32(0); x < g.camera.ProjectionPlaneWidth; x += 1 {
			distance := g.camera.ProjectionPlaneDistance * ratio / float32(math.Cos(float64(a)))

			t := g.player.Pos.Add(vec.NewRotated(g.player.Angel + a).MulValue(distance))

			cellX := t.X / g.scene.BlockSize
			cellY := t.Y / g.scene.BlockSize
			brightness := float32(100 / distance)

			//clr := g.wallTextures[scene.WallType(1)].At(int(t.X)%int(g.scene.BlockSize), int(t.Y)%int(g.scene.BlockSize)).(color.RGBA)
			//g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))

			if (int(cellX)+int(cellY))%2 == 0 {
				g.setPixel(x, y, 255*brightness, 0, 0, 255)
				//screen.DrawImage()
				//vector.DrawFilledRect(screen, x, y, 1, 1, ScaleRGBA(color.RGBA{0, 0, 255, 255}, brightness, brightness, brightness, 1), true)
			} else {
				g.setPixel(x, y, 0, 0, 255*brightness, 255)
				//vector.DrawFilledRect(screen, x, y, 1, 1, ScaleRGBA(color.RGBA{255, 0, 0, 255}, brightness, brightness, brightness, 1), true)
			}

			a += g.camera.FOV / g.camera.ProjectionPlaneWidth
		}
	}
	// for y := g.camera.ProjectionPlaneHeight/2 + 1; y < g.camera.ProjectionPlaneHeight; y += 1 {
}

func (g *Game) setPixel(x, y float32, rc, gc, bc, ac float32) {
	if x < 0 || x >= g.camera.ProjectionPlaneWidth || y < 0 || y >= g.camera.ProjectionPlaneHeight {
		return
	}

	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4] = clampColor(rc)
	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4+1] = clampColor(gc)
	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4+2] = clampColor(bc)
	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4+3] = clampColor(ac)
}

func (g *Game) clearPixels() {
	for i := range g.pixels {
		g.pixels[i] = 0
	}
}

func main() {
	game := NewGame()

	for wallType, path := range wallTextures {
		err := game.LoadWallTexture(wallType, path)
		if err != nil {
			log.Fatal("unable to load texture: ", err)
		}
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ray-casting")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
