package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ProgressBar struct {
	X             int32
	Y             int32
	Width         int32
	Height        int32
	paddingWidth  int32
	paddingHeight int32
	progress      float32
	Color         rl.Color
}

type DamageText struct {
	X        float32
	Y        float32
	Text     string
	Duration float32
}

func (pb *ProgressBar) SetProgress(newProgress float32) {
	pb.progress = newProgress
	if pb.progress < 0 {
		pb.progress = 0
	}
	if pb.progress > 1 {
		pb.progress = 1
	}
}

func (pb ProgressBar) DrawBar() {
	rl.DrawRectangle(pb.X, pb.Y, pb.Width, pb.Height, rl.White)
	rl.DrawRectangle(pb.X, pb.Y, int32(pb.progress*float32(pb.Width)), pb.Height, pb.Color)
}

func NewProgressBar(newX, newY, newWidth, newHeight, pWidth, pHeight int32, color rl.Color) ProgressBar {
	pb := ProgressBar{X: newX, Y: newY, Width: newWidth, Height: newHeight, paddingWidth: pWidth, paddingHeight: pHeight, Color: color}
	pb.progress = 0
	return pb
}

func (d *DamageText) Draw() {
	rl.DrawText(d.Text, int32(d.X), int32(d.Y), 20, rl.White)
}

func (d *DamageText) Update() {
	d.Duration -= rl.GetFrameTime()
}

func NewDamageText(x, y float32, source string, damage int32) DamageText {
	return DamageText{
		X:        x,
		Y:        y,
		Text:     fmt.Sprintf("%s took %d damage.", source, damage),
		Duration: 1.0,
	}
}
