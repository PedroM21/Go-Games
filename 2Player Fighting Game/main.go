package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1600, 900, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	idleSheet := rl.LoadTexture("sprites/idleChar.png")
	attackSheet := rl.LoadTexture("sprites/attackChar.png")
	jumpSheet := rl.LoadTexture("sprites/jumpChar.png")
	blockSheet := rl.LoadTexture("sprites/blockChar.png")

	defer rl.UnloadTexture(idleSheet)
	defer rl.UnloadTexture(attackSheet)
	defer rl.UnloadTexture(jumpSheet)
	defer rl.UnloadTexture(blockSheet)

	theme := NewColorTheme(rl.NewColor(255, 255, 255, 255), rl.NewColor(255, 128, 128, 255), rl.White)
	progressBarP1 := NewProgressBar(20, 30, 300, 50, &theme)
	progressBarP2 := NewProgressBar(1280, 30, 300, 50, &theme)

	fighter := NewFighter(rl.NewVector2(400, 600), idleSheet, attackSheet, jumpSheet, blockSheet, true)
	fighter2 := NewFighter(rl.NewVector2(1200, 600), idleSheet, attackSheet, jumpSheet, blockSheet, false)

	floor := NewBlocker(0, 650, 1600, 250, rl.DarkGray)
	floor2 := NewBlocker(0, 500, 300, 250, rl.DarkGray)
	floor3 := NewBlocker(1300, 500, 300, 250, rl.DarkGray)

	for !rl.WindowShouldClose() {
		// Handle Attack
		if rl.IsKeyPressed(rl.KeyF) {
			fighter.Attack(&fighter2)
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			fighter2.Attack(&fighter)
		}
		// Keep in bounds
		if fighter.Pos.X < 25 {
			fighter.Pos.X = 25
		}
		if fighter2.Pos.X < 25 {
			fighter2.Pos.X = 25
		}
		if fighter.Pos.X > 1575 {
			fighter.Pos.X = 1575
		}
		if fighter2.Pos.X > 1575 {
			fighter2.Pos.X = 1575
		}

		// Apply gravity
		fighter.ApplyGravity(980)
		fighter2.ApplyGravity(980)

		// Update the fighters
		fighter.Update()
		fighter2.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		fighter.DrawFighter()
		fighter2.DrawFighter()
		rl.DrawText("Player 1 Health", 20, 10, 20, rl.RayWhite)
		progressBarP1.SetProgress(fighter.Health / fighter.MaxHealth)
		progressBarP1.DrawBar()
		rl.DrawText("Player 2 Health", 1280, 10, 20, rl.RayWhite)
		progressBarP2.SetProgress(fighter2.Health / fighter2.MaxHealth)
		progressBarP2.DrawBar()

		floor.DrawBlocker()
		floor2.DrawBlocker()
		floor3.DrawBlocker()

		// Check for winner
		if fighter.Health <= 0 {
			rl.DrawRectangle(0, 0, 1600, 900, rl.Blue)
			rl.DrawText("Player 2 Wins!", 600, 450, 40, rl.RayWhite)
			rl.DrawText("Reset By Pressing R!", 600, 500, 20, rl.RayWhite)
		} else if fighter2.Health <= 0 {
			rl.DrawRectangle(0, 0, 1600, 900, rl.Blue)
			rl.DrawText("Player 1 Wins!", 600, 450, 40, rl.RayWhite)
			rl.DrawText("Reset By Pressing R!", 600, 500, 20, rl.RayWhite)
		}

		rl.EndDrawing()

		// Reset game
		if rl.IsKeyPressed(rl.KeyR) {
			fighter.Pos = rl.NewVector2(400, 600)
			fighter2.Pos = rl.NewVector2(1200, 600)

			fighter.Health = fighter.MaxHealth
			fighter2.Health = fighter2.MaxHealth
		}

	}
}
