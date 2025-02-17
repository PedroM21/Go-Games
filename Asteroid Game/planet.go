package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Planet struct {
	Pos            rl.Vector2
	Health         float32
	Radius         float32
	Color          rl.Color
	isIntersection bool
}

func NewPlanet(newPos rl.Vector2, maxHealth, newRadius float32, newColor rl.Color, intersecting bool) Planet {
	newPlanet := Planet{Pos: newPos, Health: maxHealth, Radius: newRadius, Color: newColor, isIntersection: false}
	return newPlanet
}

func (p *Planet) CollisionWithAsteroid(a *Asteroid) bool {
	if rl.Vector2Distance(p.Pos, a.Pos) <= p.Radius+a.Radius {
		return true
	}
	return false
}

func (p Planet) DrawPlanet() {
	rl.DrawCircle(int32(p.Pos.X), int32(p.Pos.Y), p.Radius, p.Color)
}
