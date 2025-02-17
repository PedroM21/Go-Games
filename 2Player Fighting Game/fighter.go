package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Fighter struct {
	StateMachine
	Animation
	PhysicsBody
	MoveSpeed   float32
	Jumping     bool
	FirstPlayer bool
	Health      float32
	MaxHealth   float32
}

func NewFighter(newPos rl.Vector2, idleSheet, attackSheet, jumpSheet, blockSheet rl.Texture2D, player bool) Fighter {
	newIdleAnim := NewAnimation(idleSheet, 0.25, IDLESTATE)
	newAttackAnim := NewAnimation(attackSheet, 0.15, ATTACKSTATE)
	newJumpAnim := NewAnimation(jumpSheet, 0.15, JUMPSTATE)
	newBlockAnim := NewAnimation(blockSheet, 0.15, BLOCKSTATE)

	stateMachine := NewStateMachine(&newIdleAnim)
	stateMachine.AddState(&newAttackAnim)
	stateMachine.AddState(&newJumpAnim)
	stateMachine.AddState(&newBlockAnim)

	return Fighter{
		StateMachine: stateMachine,
		Animation:    newIdleAnim,
		PhysicsBody:  PhysicsBody{Pos: newPos, Scale: rl.NewVector2(6.0, 6.0)},
		MoveSpeed:    300,
		FirstPlayer:  player,
		Jumping:      false,
		Health:       100,
		MaxHealth:    100,
	}
}

func (f *Fighter) Update() {
	if f.FirstPlayer {
		// Handle left and right movement
		if rl.IsKeyDown(rl.KeyD) {
			if f.Pos.X > 1275 && f.Pos.Y >= 600 {
				f.Pos.X = 1275
			}
			f.Pos.X += f.MoveSpeed * rl.GetFrameTime()
		} else if rl.IsKeyDown(rl.KeyA) {
			if f.Pos.X < 325 && f.Pos.Y >= 600 {
				f.Pos.X = 325
			}
			f.Pos.X -= f.MoveSpeed * rl.GetFrameTime()
		}
		// Handle Block
		if rl.IsKeyDown(rl.KeyLeftShift) {
			f.Block()
		} else {
			if f.StateMachine.GetStateName() == BLOCKSTATE {
				f.ChangeState(IDLESTATE)
			}
		}
		// Handle Jump
		if rl.IsKeyPressed(rl.KeySpace) && !f.Jumping {
			f.Jumping = true
			f.Vel.Y = -900
			f.ChangeState(JUMPSTATE)
		}
		// Apply gravity
		if f.Jumping {
			f.ApplyGravity(980)
			f.ApplyVel()
			// Check for whether player is jumping on normal floor or elevated floor
			if f.Pos.Y >= 600 {
				f.Pos.Y = 600
				f.Vel.Y = 0
				f.Jumping = false
				f.ChangeState(IDLESTATE)
			} else if f.Pos.X <= 300 && f.Pos.Y >= 450 {
				f.Pos.Y = 450
				f.Vel.Y = 0
				f.Jumping = false
				f.ChangeState(IDLESTATE)
			} else if f.Pos.X >= 1300 && f.Pos.Y >= 450 {
				f.Pos.Y = 450
				f.Vel.Y = 0
				f.Jumping = false
				f.ChangeState(IDLESTATE)
			}
		}
	} else {
		// Handle Movement for player 2
		if rl.IsKeyDown(rl.KeyRight) {
			if f.Pos.X > 1275 && f.Pos.Y >= 600 {
				f.Pos.X = 1275
			}
			f.Pos.X += f.MoveSpeed * rl.GetFrameTime()
		} else if rl.IsKeyDown(rl.KeyLeft) {
			if f.Pos.X < 325 && f.Pos.Y >= 600 {
				f.Pos.X = 325
			}
			f.Pos.X -= f.MoveSpeed * rl.GetFrameTime()
		}
		// Handle Block for player 2
		if rl.IsKeyDown(rl.KeyRightShift) {
			f.Block()
		} else {
			if f.StateMachine.GetStateName() == BLOCKSTATE {
				f.ChangeState(IDLESTATE)
			}
		}
		//Handle jump for player 2
		if rl.IsKeyPressed(rl.KeyRightControl) && !f.Jumping {
			f.Jumping = true
			f.Vel.Y = -900
			f.ChangeState(JUMPSTATE)
		}
		// Apply gravity for player 2
		if f.Jumping {
			f.ApplyGravity(980)
			f.ApplyVel()
			// Check for whether player 2 is jumping on normal floor or elevated floor
			if f.Pos.Y >= 600 {
				f.Pos.Y = 600
				f.Vel.Y = 0
				f.Jumping = false
				f.ChangeState(IDLESTATE)
			} else if f.Pos.X <= 300 && f.Pos.Y >= 450 {
				f.Pos.Y = 450
				f.Vel.Y = 0
				f.Jumping = false
				f.ChangeState(IDLESTATE)
			} else if f.Pos.X >= 1300 && f.Pos.Y >= 450 {
				f.Pos.Y = 450
				f.Vel.Y = 0
				f.Jumping = false
				f.ChangeState(IDLESTATE)
			}
		}
	}
	f.Tick()
}

func (f *Fighter) Attack(opponent *Fighter) {
	f.ChangeState(ATTACKSTATE)

	var hitbox rl.Rectangle
	attackWidth := float32(50)
	attackHeight := f.Scale.Y * 20

	if f.Pos.X < opponent.Pos.X {
		hitbox = rl.NewRectangle(f.Pos.X+f.Scale.X, f.Pos.Y-attackHeight/2, attackWidth, attackHeight)
	} else {
		hitbox = rl.NewRectangle(f.Pos.X-attackWidth, f.Pos.Y-attackHeight/2, attackWidth, attackHeight)
	}

	// Check for collision between the attack hitbox and the opponent
	if rl.CheckCollisionRecs(hitbox, opponent.GetRectangle()) {
		if opponent.StateMachine.GetStateName() == BLOCKSTATE {
			fmt.Println("Blocked!")
		} else {
			fmt.Println("Hit!")
			opponent.Health -= 10
		}
	} else {
		fmt.Println("Miss!")
	}
}

func (f *Fighter) Block() {
	f.ChangeState(BLOCKSTATE)
	fmt.Println("Blocking...")

}

func (f *Fighter) DrawFighter() {
	frameWidth := float32(f.SpriteSheet.Height)
	sourceRect := rl.NewRectangle(float32(f.CurrentFrame)*frameWidth, 0, frameWidth, frameWidth)
	destRect := rl.NewRectangle(f.Pos.X, f.Pos.Y, frameWidth*f.Scale.X, frameWidth*f.Scale.Y)
	origin := rl.NewVector2(destRect.Width/2, destRect.Height/2)
	rl.DrawTexturePro(f.SpriteSheet, sourceRect, destRect, origin, 0, rl.White)
}
