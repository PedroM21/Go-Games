package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Blocker struct {
	Pos   rl.Vector2
	Size  rl.Vector2
	Color rl.Color
}

func NewBlocker(px, py, sx, sy float32, Color rl.Color) Blocker {
	return Blocker{Pos: rl.NewVector2(px, py), Size: rl.NewVector2(sx, sy), Color: Color}
}

func (bl Blocker) DrawBlocker() {
	rl.DrawRectangle(int32(bl.Pos.X), int32(bl.Pos.Y), int32(bl.Size.X), int32(bl.Size.Y), bl.Color)
}
