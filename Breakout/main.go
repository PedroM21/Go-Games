package main

import (
	"encoding/json"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PhysicsBody struct {
	Pos        rl.Vector2
	Vel        rl.Vector2
	Radius     float32
	isLaunched bool
}

type Object struct {
	pos          rl.Vector2
	objectWidth  float32
	objectHeight float32
}

func NewObject(newPos rl.Vector2, objWidth float32, objHeight float32) *Object {
	newObj := Object{pos: newPos, objectWidth: objWidth, objectHeight: objHeight}
	return &newObj
}

func (o *Object) DrawObject() {
	rl.DrawRectangle(int32(o.pos.X), int32(o.pos.Y), int32(o.objectWidth), int32(o.objectHeight), rl.Orange)
}

func (pb *PhysicsBody) VelocityTick() {
	if pb.isLaunched {
		adjustedVel := rl.Vector2Scale(pb.Vel, rl.GetFrameTime())
		pb.Pos = rl.Vector2Add(pb.Pos, adjustedVel)
	}
}

func NewPhysicsBody(newPos rl.Vector2, newVel rl.Vector2, newRadius float32) PhysicsBody {
	pb := PhysicsBody{Pos: newPos, Vel: newVel, Radius: newRadius, isLaunched: false}
	return pb
}

func (pb *PhysicsBody) BounceOffPaddle(paddlePos rl.Vector2, paddleWidth float32, paddleHeight float32) {
	// check for intersection
	if pb.Pos.X+pb.Radius > paddlePos.X && pb.Pos.Y+pb.Radius > paddlePos.Y && pb.Pos.X-pb.Radius < paddlePos.X+paddleWidth && pb.Pos.Y-pb.Radius < paddlePos.Y+paddleHeight {

		pb.Vel.Y *= -1
		relativePos := (pb.Pos.X - paddlePos.X) / paddleWidth
		pb.Vel.X = (relativePos - 0.5) * 400
	}
}

func (pb *PhysicsBody) BounceOffObjects(obj *Object) bool {
	// check for intersection
	if pb.Pos.X+pb.Radius > obj.pos.X && pb.Pos.Y+pb.Radius > obj.pos.Y && pb.Pos.X-pb.Radius < obj.pos.X+obj.objectWidth && pb.Pos.Y-pb.Radius < obj.pos.Y+obj.objectHeight {
		if pb.Pos.X < obj.pos.X || pb.Pos.X > obj.pos.X+obj.objectWidth {
			pb.Vel.X *= -1
		} else {
			pb.Vel.Y *= -1
		}
		return true
	}
	return false
}

func (p *PhysicsBody) Save(filename string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	fmt.Println("SAVING...")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (p *PhysicsBody) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, p)
}

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.InitAudioDevice()
	defer rl.CloseWindow()

	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	var paddleWidth float32 = 100
	var paddleHeight float32 = 5
	var ballSize float32 = 10

	// Creates the paddle and ball
	paddlePos := rl.NewVector2(float32(screenWidth)/2-paddleWidth/2, float32(screenHeight)-50)
	ball := NewPhysicsBody(rl.NewVector2(paddlePos.X+paddleWidth/2, paddlePos.Y-ballSize), rl.NewVector2(20, 0), ballSize)

	// Creates the objects
	startingPosX := 100
	startingPosY := 50
	spacingX := 60
	spacingY := 60
	rows := 3
	cols := 10
	objects := make([]*Object, 0, rows*cols)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			posX := startingPosX + col*spacingX
			posY := startingPosY + row*spacingY
			newObject := NewObject(rl.NewVector2(float32(posX), float32(posY)), 55, 55)
			objects = append(objects, newObject)
		}
	}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		// Draws paddle and ball
		rl.DrawRectangle(int32(paddlePos.X), int32(paddlePos.Y), int32(paddleWidth), int32(paddleHeight), rl.Gray)
		rl.DrawCircle(int32(ball.Pos.X), int32(ball.Pos.Y), ballSize, rl.White)

		// Draws the objects
		for _, object := range objects {
			object.DrawObject()
		}

		// Keep Paddle from leaving screen
		if paddlePos.X > float32(screenWidth)-paddleWidth {
			paddlePos.X = float32(screenWidth) - paddleWidth
		} else if paddlePos.X < 0 {
			paddlePos.X = 0
		}

		// Bounce ball off walls
		if ball.Pos.X < 0 {
			ball.Vel.X *= -1
		} else if ball.Pos.X > float32(screenWidth)-ballSize {
			ball.Vel.X *= -1
		}
		if ball.Pos.Y < 0 {
			ball.Vel.Y *= -1
		}

		// Handle paddle and ball movement if not launched
		if rl.IsKeyDown(rl.KeyA) {
			paddlePos.X -= 10
			if !ball.isLaunched {
				ball.Pos.X = paddlePos.X + paddleWidth/2
			}
		} else if rl.IsKeyDown(rl.KeyD) {
			paddlePos.X += 10
			if !ball.isLaunched {
				ball.Pos.X = paddlePos.X + paddleWidth/2
			}
		}

		// Load and save
		if rl.IsKeyPressed(rl.KeyO) {
			ball.Save("ball.json")
		}
		if rl.IsKeyPressed(rl.KeyP) {
			ball.Load("ball.json")
		}

		// Launch ball
		if rl.IsKeyPressed(rl.KeySpace) {
			ball.Vel = rl.NewVector2(0, -200)
			ball.isLaunched = true
		}
		ball.VelocityTick()

		ball.BounceOffPaddle(paddlePos, paddleWidth, paddleHeight)

		// Removes objects when collided
		for i := len(objects) - 1; i >= 0; i-- {
			if ball.BounceOffObjects((objects[i])) {
				objects = append(objects[:i], objects[i+1:]...)
			}
		}

		// Reset Game
		if ball.Pos.Y > float32(screenHeight) || len(objects) <= 0 {
			paddlePos.X = float32(screenWidth)/2 - paddleWidth/2

			ball.Pos = rl.NewVector2(paddlePos.X+paddleWidth/2, paddlePos.Y-ballSize)
			ball.Vel = rl.NewVector2(20, 0)
			ball.isLaunched = false

			for row := 0; row < rows; row++ {
				for col := 0; col < cols; col++ {
					posX := startingPosX + col*spacingX
					posY := startingPosY + row*spacingY
					newObject := NewObject(rl.NewVector2(float32(posX), float32(posY)), 55, 55)
					objects = append(objects, newObject)
				}
			}
		}

		rl.EndDrawing()
	}
}
