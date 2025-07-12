package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"log"
	"math"
	"ray-casting/internal/camera"
	"ray-casting/internal/player"
	"ray-casting/internal/raycast"
	"ray-casting/internal/scene"
	"ray-casting/pkg/helpers"
	"ray-casting/pkg/vec"
)

const SceneRows int = 24
const SceneCols int = 24
const BlockSize float32 = 64
const ProjectionPlaneWidth float32 = 640
const ProjectionPlaneHeight float32 = 480
const PlayerFOV int = 60
const PlayerHeight float32 = 32

type Game struct {
	scene  scene.Scene
	player player.Player
	camera camera.Camera
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
			SceneRows,
			SceneCols,
			BlockSize,
		),
		player: player.NewPlayer(
			vec.NewVec2((float32(SceneRows-1)*BlockSize)*.5, (float32(SceneCols-1)*BlockSize)*.5),
			0,
		),
		camera: camera.NewCamera(
			ProjectionPlaneWidth,
			ProjectionPlaneHeight,
			PlayerFOV,
			PlayerHeight,
		),
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
	offset := vec.NewVec2(float32(ProjectionPlaneWidth), float32(ProjectionPlaneHeight)).Sub(vec.NewVec2(BlockSize*float32(SceneRows), BlockSize*float32(SceneCols)).MulValue(scale)).Sub(vec.NewVec2(20, 20))
	g.renderMinimap(screen, offset, scale)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return int(g.camera.ProjectionPlaneWidth), int(g.camera.ProjectionPlaneHeight)
}

func (g *Game) renderWall(screen *ebiten.Image) {
	// 2.5D
	h := vec.NewVec2(0, float32(ProjectionPlaneHeight)*.5)

	a := -g.camera.FOV * .5

	for i := range int(g.camera.ProjectionPlaneWidth) {
		distance := raycast.Cast(g.scene, g.player.Pos, g.player.Angel+a)

		sliceHeight := BlockSize / (distance * float32(math.Cos(float64(a)))) * g.camera.ProjectionPlaneDistance

		from := h.Add(vec.NewVec2(float32(i), -sliceHeight*.5))
		to := from.Add(vec.NewVec2(0, sliceHeight))

		vector.StrokeLine(screen, from.X, from.Y, to.X, to.Y, 1, helpers.Color(0xFFFFFFFF), true)

		a += g.camera.FOV / g.camera.ProjectionPlaneWidth
	}
}

func (g *Game) renderMinimap(screen *ebiten.Image, offset vec.Vec2, scale float32) {
	{
		for row := range g.scene.Rows {
			for col := range g.scene.Cols {
				p := vec.NewVec2(float32(col), float32(row)).MulValue(BlockSize).MulValue(scale).Add(offset)
				d := vec.NewVec2(BlockSize, BlockSize).MulValue(scale)

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

	// Rays
	//a := -g.camera.FOV * .5
	//
	//for i := range int(g.camera.ProjectionPlaneWidth) {
	//	p := g.position.Mul(scale).Add(offset)
	//	r := g.position.Add(vec.NewVec2(1, 0).Rot(ray.angel).MulValue(ray.distance)).Mul(scale).Add(offset)
	//	vector.StrokeLine(screen, p.X, p.Y, r.X, r.Y, 1, helpers.Color(0xFF0000FF), true)
	//}
}

func (g *Game) renderFloor() {

}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ray-casting")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
