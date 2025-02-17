package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	IDLESTATE   = "idle"
	ATTACKSTATE = "attack"
	JUMPSTATE   = "jump"
	BLOCKSTATE  = "block"
)

type Animation struct {
	SpriteSheet  rl.Texture2D
	MaxFrames    int
	CurrentFrame int
	FrameTime    float32
	Timer        float32
	Name         string
}

func (a *Animation) TickState() {
	a.Timer += rl.GetFrameTime()
	if a.Timer >= a.FrameTime {
		a.Timer = 0
		a.CurrentFrame++
		if a.CurrentFrame >= a.MaxFrames {
			a.CurrentFrame = 0
		}
	}
}

func (a *Animation) GetName() string {
	return a.Name
}

func (a *Animation) ResetTime() {
	a.Timer = 0
}

func NewAnimation(spriteSheet rl.Texture2D, frameTime float32, name string) Animation {
	frameWidth := spriteSheet.Height
	totalFrames := int(spriteSheet.Width / frameWidth)

	return Animation{
		SpriteSheet:  spriteSheet,
		MaxFrames:    totalFrames,
		CurrentFrame: 0,
		FrameTime:    frameTime,
		Timer:        0,
		Name:         name,
	}
}
