package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Asteroid struct {
	Pos            rl.Vector2
	Velocity       rl.Vector2
	Radius         float32
	MoveSpeed      float32
	Color          rl.Color
	isIntersection bool
}

func NewAsteroid(newPos rl.Vector2, vel rl.Vector2, newRadius, speed float32, newColor rl.Color) Asteroid {
	newAsteroid := Asteroid{Pos: newPos, Velocity: vel, Radius: newRadius, MoveSpeed: speed, Color: newColor}
	return newAsteroid
}

func (a *Asteroid) MoveAsteroid(dest rl.Vector2) {
	direction := rl.Vector2Normalize(rl.Vector2Subtract(dest, a.Pos))

	if a.Radius > 20 {
		a.Pos.X += direction.X * a.MoveSpeed * rl.GetFrameTime()
		a.Pos.Y += direction.Y * a.MoveSpeed * rl.GetFrameTime()
	}
}

func SplitAsteroids(a Asteroid) (Asteroid, Asteroid) {
	newRadius := a.Radius / 2

	offset := float32(100)

	asteroid1 := Asteroid{
		Pos:       rl.Vector2{X: a.Pos.X + offset, Y: a.Pos.Y + offset},
		Velocity:  a.Velocity,
		Radius:    newRadius,
		MoveSpeed: a.MoveSpeed,
		Color:     a.Color,
	}

	asteroid2 := Asteroid{
		Pos:       rl.Vector2{X: a.Pos.X - offset, Y: a.Pos.Y - offset},
		Velocity:  a.Velocity,
		Radius:    newRadius,
		MoveSpeed: a.MoveSpeed,
		Color:     a.Color,
	}

	if newRadius < 25 {
		asteroid1.Color = rl.NewColor(72, 255, 0, 255)
		asteroid2.Color = rl.NewColor(72, 255, 0, 255)
	}

	return asteroid1, asteroid2
}

func (a Asteroid) DrawAsteroid() {
	rl.DrawCircle(int32(a.Pos.X), int32(a.Pos.Y), a.Radius, a.Color)
}
