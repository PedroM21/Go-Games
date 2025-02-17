package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Ship struct {
	Sprite    rl.Texture2D
	Pos       rl.Vector2
	Rotation  float32
	Scale     float32
	MoveSpeed float32
	Color     rl.Color
}

func NewShip(sprite rl.Texture2D, pos rl.Vector2, rotation, size, speed float32, color rl.Color) Ship {
	newShip := Ship{Sprite: sprite, Pos: pos, Rotation: rotation, Scale: size, MoveSpeed: speed, Color: color}
	return newShip
}

func (s *Ship) CollectCargo(a *Asteroid) bool {
	distance := rl.Vector2Distance(s.Pos, a.Pos)

	shipRadius := float32(s.Sprite.Width) * s.Scale / 2
	if a.Radius < 20 {
		if distance <= shipRadius+a.Radius {
			return true
		}
	}
	return false
}

func (s *Ship) DepositeCargo(p *Planet) bool {
	distance := rl.Vector2Distance(s.Pos, p.Pos)

	shipRadius := float32(s.Sprite.Width) * s.Scale / 2
	if distance <= shipRadius+p.Radius {
		return true
	}
	return false
}

func (s Ship) DrawShip() {
	sourceRect := rl.NewRectangle(0, 0, float32(s.Sprite.Width), float32(s.Sprite.Height))
	destRect := rl.NewRectangle(s.Pos.X, s.Pos.Y, float32(s.Sprite.Width)*s.Scale, float32(s.Sprite.Height)*s.Scale)
	origin := rl.Vector2Scale(rl.NewVector2(float32(s.Sprite.Width)/2, float32(s.Sprite.Height)/2), s.Scale)
	rl.DrawTexturePro(s.Sprite, sourceRect,
		destRect,
		origin, s.Rotation, s.Color)
}
