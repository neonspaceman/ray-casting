package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
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
const projectionPlaneWidth float32 = 640
const projectionPlaneHeight float32 = 480
const playerFOV int = 60
const playerHeight float32 = 32

var wallTextures = map[scene.WallType]string{
	1: "assets/Bark.png",
	2: "assets/WalkStone.png",
	3: "assets/WalkStone.png",
	4: "assets/WalkStone.png",
}

type Game struct {
	scene        scene.Scene
	player       player.Player
	camera       camera.Camera
	wallTextures map[scene.WallType]*ebiten.Image
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%#v, %f", g.player.Pos, g.player.Angel))

	g.renderWall(screen)

	scale := float32(.1)
	offset := vec.NewVec2(g.camera.ProjectionPlaneWidth, g.camera.ProjectionPlaneHeight).Sub(vec.NewVec2(g.scene.Width, g.scene.Height).MulValue(scale)).Sub(vec.NewVec2(20, 20))
	g.renderMinimap(screen, offset, scale)
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

func (g *Game) renderWall(screen *ebiten.Image) {
	h := vec.NewVec2(0, float32(g.camera.ProjectionPlaneHeight)*.5)

	a := -g.camera.FOV * .5

	for i := range int(g.camera.ProjectionPlaneWidth) {
		distance, wallOffset, wall := raycast.Cast(g.scene, g.player.Pos, g.player.Angel+a)

		sliceHeight := g.scene.BlockSize / (distance * float32(math.Cos(float64(a)))) * g.camera.ProjectionPlaneDistance

		from := h.Add(vec.NewVec2(float32(i), -sliceHeight*.5))
		to := from.Add(vec.NewVec2(0, sliceHeight))

		if texture, ok := g.wallTextures[wall]; ok {
			crop := texture.SubImage(image.Rect(wallOffset, 0, wallOffset+1, texture.Bounds().Dy())).(*ebiten.Image)

			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Scale(1, float64(sliceHeight/g.scene.BlockSize))
			opts.GeoM.Translate(float64(from.X), float64(from.Y))
			opts.ColorScale.Scale(160/distance, 160/distance, 160/distance, 1)

			screen.DrawImage(crop, &opts)
		} else {
			vector.StrokeLine(screen, from.X, from.Y, to.X, to.Y, 1, helpers.Color(0xFFFFFFFF), true)
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

func (g *Game) renderFloor() {

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
