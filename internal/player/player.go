package player

import (
	"math"
	"ray-casting/pkg/vec"
)

type Player struct {
	Pos   vec.Vec2
	Angel float32
}

func NewPlayer(pos vec.Vec2, angel float32) Player {
	return Player{
		Pos:   pos,
		Angel: angel,
	}
}

func (p *Player) Up(speed float32) {
	p.Pos = p.Pos.Add(vec.NewRotated(p.Angel).MulValue(speed))
}

func (p *Player) Left(speed float32) {
	p.Angel -= speed
	if p.Angel < 0 {
		p.Angel = 2 * math.Pi
	}
}

func (p *Player) Right(speed float32) {
	p.Angel += speed
	if p.Angel > 2*math.Pi {
		p.Angel = 0
	}
}

func (p *Player) Down(speed float32) {
	p.Pos = p.Pos.Add(vec.NewRotated(p.Angel + math.Pi).MulValue(speed))
}
