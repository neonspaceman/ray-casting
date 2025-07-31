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
	"slices"
)

const sceneRows int = 15
const sceneCols int = 15
const blockSize float32 = 16
const projectionPlaneWidth float32 = 320
const projectionPlaneHeight float32 = 240
const playerFOV int = 60
const playerHeight float32 = 8
const texturesPath = "assets/textures.png"
const floorTexture = 4
const ceilTexture = 3

type Game struct {
	scene    scene.Scene
	player   player.Player
	camera   camera.Camera
	textures *ebiten.Image
	pixels   []byte
	zBuffer  []float32
}

func NewGame() *Game {
	return &Game{
		scene: scene.NewScene(
			[][]int{
				{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
			},
			sceneRows,
			sceneCols,
			blockSize,
			[]scene.Sprite{
				{vec.NewVec2(blockSize*5.5, blockSize*8.5), 5},
				{vec.NewVec2(blockSize*4.5, blockSize*8.5), 5},
			},
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
		pixels:  make([]byte, int(projectionPlaneWidth)*int(projectionPlaneHeight)*4),
		zBuffer: make([]float32, int(projectionPlaneWidth)),
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
	g.renderFloor(floorTexture)
	g.renderCell(ceilTexture)
	g.renderWall()
	g.renderSprites()

	screen.WritePixels(g.pixels)

	scale := float32(.25)
	offset := vec.NewVec2(g.camera.ProjectionPlaneWidth, g.camera.ProjectionPlaneHeight).Sub(vec.NewVec2(g.scene.Width, g.scene.Height).MulValue(scale)).Sub(vec.NewVec2(20, 20))
	g.renderMinimap(screen, offset, scale)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%#v, %f\nFPS: %f", g.player.Pos, g.player.Angel, ebiten.ActualFPS()))
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(g.camera.ProjectionPlaneWidth), int(g.camera.ProjectionPlaneHeight)
}

func (g *Game) LoadTextures(path string) error {
	t, _, err := ebitenutil.NewImageFromFile(path)

	if err != nil {
		return err
	}

	g.textures = t

	return nil
}

func (g *Game) renderWall() {
	h := g.camera.ProjectionPlaneHeight * .5

	a := -g.camera.FOV * .5

	for x := float32(0); x < g.camera.ProjectionPlaneWidth; x += 1 {
		distance, wallOffset, texture := raycast.Cast(g.scene, g.player.Pos, g.player.Angel+a)

		brightness := 32 / distance
		sliceHeight := g.scene.BlockSize / (distance * float32(math.Cos(float64(a)))) * g.camera.ProjectionPlaneDistance

		ratio := g.scene.BlockSize / sliceHeight
		ty := float32(0)

		for y := h - sliceHeight*.5; y <= h+sliceHeight*.5; y += 1 {
			clr := g.textures.At((texture-1)*int(g.scene.BlockSize)+wallOffset, int(ty)).(color.RGBA)

			g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))

			ty += ratio
		}

		g.zBuffer[int(x)] = distance

		a += g.camera.FOV / g.camera.ProjectionPlaneWidth
	}
}

func (g *Game) renderMinimap(screen *ebiten.Image, offset vec.Vec2, scale float32) {
	// MAP
	{
		for row := range g.scene.Rows {
			for col := range g.scene.Cols {
				p := vec.NewVec2(float32(col), float32(row)).MulValue(g.scene.BlockSize).MulValue(scale).Add(offset)
				d := vec.NewVec2(g.scene.BlockSize, g.scene.BlockSize).MulValue(scale)

				wallType := g.scene.Walls[row][col]

				if wallType != scene.None {
					vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, color.RGBA{255, 255, 255, 255}, true)
				} else {
					vector.DrawFilledRect(screen, p.X, p.Y, d.X, d.Y, color.RGBA{0, 0, 0, 255}, true)
				}
			}
		}
	}

	//Player
	{
		from := g.player.Pos.MulValue(scale).Add(offset)

		to1 := g.player.Pos.Add(vec.NewRotated(g.player.Angel - g.camera.FOV*.5).MulValue(50)).MulValue(scale).Add(offset)
		to2 := g.player.Pos.Add(vec.NewRotated(g.player.Angel + g.camera.FOV*.5).MulValue(50)).MulValue(scale).Add(offset)

		vector.DrawFilledCircle(screen, from.X, from.Y, 2, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeLine(screen, from.X, from.Y, to1.X, to1.Y, 1, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeLine(screen, from.X, from.Y, to2.X, to2.Y, 1, color.RGBA{255, 0, 0, 255}, true)
	}

	//Sprites
	{
		startAngel := g.player.Angel - g.camera.FOV*.5
		if startAngel < 0 {
			startAngel += 2 * math.Pi
		}

		endAngel := g.player.Angel + g.camera.FOV*.5
		if endAngel >= 2*math.Pi {
			endAngel -= 2 * math.Pi
		}

		for _, sprite := range g.scene.Sprites {
			pos := sprite.Pos.MulValue(scale).Add(offset)
			vector.DrawFilledCircle(screen, pos.X, pos.Y, 1, color.RGBA{255, 0, 0, 255}, true)
		}
	}
}

func (g *Game) renderFloor(textureIndex int) {
	for y := g.camera.ProjectionPlaneHeight*.5 + 1; y < g.camera.ProjectionPlaneHeight; y += 1 {
		a := -g.camera.FOV * .5
		ratio := g.camera.Z / (y - g.camera.ProjectionPlaneHeight*.5)

		for x := float32(0); x < g.camera.ProjectionPlaneWidth; x += 1 {
			distance := g.camera.ProjectionPlaneDistance * ratio / float32(math.Cos(float64(a)))

			t := g.player.Pos.Add(vec.NewRotated(g.player.Angel + a).MulValue(distance))

			brightness := 32 / distance

			textureX := (textureIndex-1)*int(g.scene.BlockSize) + int(t.X)%int(g.scene.BlockSize)
			textureY := int(t.Y) % int(g.scene.BlockSize)

			clr := g.textures.At(textureX, textureY).(color.RGBA)
			g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))

			// another procedural floor texture - red-blue floor
			//cellX := t.X / g.scene.BlockSize
			//cellY := t.Y / g.scene.BlockSize
			//if (int(cellX)+int(cellY))%2 == 0 {
			//	g.setPixel(x, y, 255*brightness, 0, 0, 255)
			//} else {
			//	g.setPixel(x, y, 0, 0, 255*brightness, 255)
			//}

			a += g.camera.FOV / g.camera.ProjectionPlaneWidth
		}
	}
}

func (g *Game) renderCell(textureIndex int) {
	for y := float32(0); y <= g.camera.ProjectionPlaneHeight*.5; y += 1 {
		a := -g.camera.FOV * .5
		ratio := g.camera.Z / (g.camera.ProjectionPlaneHeight*.5 - y)

		for x := float32(0); x < g.camera.ProjectionPlaneWidth; x += 1 {
			distance := g.camera.ProjectionPlaneDistance * ratio / float32(math.Cos(float64(a)))

			t := g.player.Pos.Add(vec.NewRotated(g.player.Angel + a).MulValue(distance))

			brightness := 16 / distance

			textureX := (textureIndex-1)*int(g.scene.BlockSize) + int(t.X)%int(g.scene.BlockSize)
			textureY := int(t.Y) % int(g.scene.BlockSize)

			clr := g.textures.At(textureX, textureY).(color.RGBA)
			g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))

			// another procedural floor texture - red-blue floor
			//cellX := t.X / g.scene.BlockSize
			//cellY := t.Y / g.scene.BlockSize
			//if (int(cellX)+int(cellY))%2 == 0 {
			//	g.setPixel(x, y, 255*brightness, 0, 0, 255)
			//} else {
			//	g.setPixel(x, y, 0, 0, 255*brightness, 255)
			//}

			a += g.camera.FOV / g.camera.ProjectionPlaneWidth
		}
	}
}

func (g *Game) renderSprites() {
	slices.SortFunc(g.scene.Sprites, func(a, b scene.Sprite) int {
		aLen, bLen := a.Pos.Sub(g.player.Pos).Len2(), b.Pos.Sub(g.player.Pos).Len2()
		if aLen > bLen {
			return -1
		} else {
			return 1
		}
	})

	startAngel := g.player.Angel - g.camera.FOV*.5

	h := g.camera.ProjectionPlaneHeight * .5

	for _, sprite := range g.scene.Sprites {

		angel := sprite.Pos.Sub(g.player.Pos).Rad()

		if g.player.Angel >= math.Pi*3/2 && g.player.Angel <= math.Pi*2 && angel >= 0 && angel <= math.Pi*1/2 {
			angel += 2 * math.Pi
		}
		if g.player.Angel >= 0 && g.player.Angel <= math.Pi*1/2 && angel >= math.Pi*3/2 && angel <= math.Pi*2 {
			angel -= 2 * math.Pi
		}

		spriteDistance := sprite.Pos.Sub(g.player.Pos).Len()
		spriteHeight := (g.scene.BlockSize * g.camera.ProjectionPlaneDistance) / spriteDistance
		spriteX := (angel - startAngel) * g.camera.ProjectionPlaneWidth / g.camera.FOV
		//spriteY := h - spriteHeight*.5 //g.camera.ProjectionPlaneHeight * .5

		brightness := 16 / spriteDistance
		ratio := g.scene.BlockSize / spriteHeight

		tx := float32(sprite.Type-1) * g.scene.BlockSize

		for x := spriteX - spriteHeight*.5; x < spriteX+spriteHeight*.5; x += 1 {
			ty := float32(0)

			if x > 0 && x < g.camera.ProjectionPlaneWidth && spriteDistance < g.zBuffer[int(x)] {
				for y := h - spriteHeight*.5; y <= h+spriteHeight*.5; y += 1 {
					clr := g.textures.At(int(tx), int(ty)).(color.RGBA)

					// Skip purple color
					if clr.R != 152 || clr.G != 0 || clr.B != 136 {
						g.setPixel(x, y, float32(clr.R)*brightness, float32(clr.G)*brightness, float32(clr.B)*brightness, float32(clr.A))
					}

					ty += ratio
				}
			}

			tx += ratio
		}
	}
}

func (g *Game) setPixel(x, y float32, rc, gc, bc, ac float32) {
	if x < 0 || x >= g.camera.ProjectionPlaneWidth || y < 0 || y >= g.camera.ProjectionPlaneHeight {
		return
	}

	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4] = helpers.ClampColor(rc)
	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4+1] = helpers.ClampColor(gc)
	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4+2] = helpers.ClampColor(bc)
	g.pixels[int(y)*int(g.camera.ProjectionPlaneWidth)*4+int(x)*4+3] = helpers.ClampColor(ac)
}

func (g *Game) clearPixels() {
	for i := range g.pixels {
		g.pixels[i] = 0
	}
}

func main() {
	game := NewGame()
	err := game.LoadTextures(texturesPath)

	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ray-casting")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
