package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PhysicsBody struct {
	Pos   rl.Vector2
	Vel   rl.Vector2
	Scale rl.Vector2
}

func (b PhysicsBody) GetRectangle() rl.Rectangle {
	return rl.NewRectangle(b.Pos.X-b.Scale.X/2, b.Pos.Y-b.Scale.Y/2, b.Scale.X, b.Scale.Y)
}

func (b *PhysicsBody) ApplyGravity(gravity float32) {
	b.Vel.Y += gravity * rl.GetFrameTime()
}

func (b *PhysicsBody) ApplyVel() {
	b.Pos = rl.Vector2Add(b.Pos, rl.Vector2Scale(b.Vel, rl.GetFrameTime()))
}
