package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	Pos      rl.Vector2
	Velocity rl.Vector2
	Radius   float32
}

func NewProjectile(newPos rl.Vector2, vel rl.Vector2, newRadius float32) Projectile {
	newProjectile := Projectile{Pos: newPos, Velocity: vel, Radius: newRadius}
	return newProjectile
}

func (p *Projectile) UpdateProjectile() {
	p.Pos.X += p.Velocity.X * rl.GetFrameTime()
	p.Pos.Y += p.Velocity.Y * rl.GetFrameTime()
}

func (p *Projectile) DestroyAsteroid(a *Asteroid) bool {
	if rl.Vector2Distance(p.Pos, a.Pos) <= p.Radius+a.Radius {
		return true
	}
	return false
}

func (p Projectile) DrawProjectile() {
	rl.DrawCircle(int32(p.Pos.X), int32(p.Pos.Y), p.Radius, rl.White)
}

func DestroyProjectile(projectiles *[]Projectile, maxDistance float32) {
	for i := len(*projectiles) - 1; i >= 0; i-- {
		distance := rl.Vector2Distance((*projectiles)[i].Pos, rl.NewVector2(0, 0))
		if distance > maxDistance {
			*projectiles = append((*projectiles)[:i], (*projectiles)[i+1:]...)
		}
	}
}
