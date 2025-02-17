package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const animationSize = 6

type Animation struct {
	Pos          rl.Vector2
	SpriteSheet  rl.Texture2D
	MaxIndex     int
	CurrentIndex int
	Timer        float32
	SwitchTime   float32
}

func NewAnimation(newPos rl.Vector2, newSheet rl.Texture2D, spriteNum int, newTime float32) Animation {
	newAnimation := Animation{Pos: newPos, SpriteSheet: newSheet, MaxIndex: spriteNum, CurrentIndex: 0, Timer: 0, SwitchTime: newTime}
	return newAnimation
}

func (a *Animation) UpdateTime() {
	a.Timer += rl.GetFrameTime()
	if a.Timer > a.SwitchTime {
		a.Timer = 0
		a.CurrentIndex++
	}

	if a.CurrentIndex > a.MaxIndex {
		a.CurrentIndex = 0
	}
}

func (a Animation) DrawAnimation() {
	sourceRect := rl.NewRectangle(float32(32*a.CurrentIndex), 0, 32, 32)
	destRect := rl.NewRectangle(a.Pos.X, a.Pos.Y, 32*animationSize, 32*animationSize)
	rl.DrawTexturePro(a.SpriteSheet, sourceRect, destRect, rl.Vector2Zero(), 0, rl.White)
}
