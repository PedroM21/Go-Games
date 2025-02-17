package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1600, 900, "raylib [core] example - basic window")
	rl.InitAudioDevice()
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	var gameInProgress bool = true

	// Sound effects and music
	bgMusic := rl.LoadMusicStream("audio/Nokia Espionage.ogg")
	laserShot := rl.LoadSound("audio/Laser.ogg")
	explosion := rl.LoadSound("audio/retro-explode.ogg")
	heal := rl.LoadSound("audio/Heal.ogg")
	cargoPickUp := rl.LoadSound("audio/Cargo.ogg")
	takeDamage := rl.LoadSound("audio/damage.ogg")
	rl.PlayMusicStream(bgMusic)
	rl.SetSoundVolume(laserShot, 0.25)
	rl.SetSoundVolume(explosion, 0.25)
	rl.SetSoundVolume(heal, 0.5)

	// Camera
	cam := rl.NewCamera2D(
		rl.NewVector2(0, 0),
		rl.NewVector2(0, 0),
		0,
		1,
	)
	cam.Offset = rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2))
	var rotationSpeed float32 = 250

	// Planet variables
	planetPos := rl.NewVector2(0, 0)
	var maxHealth int32 = 10
	var currentHealth int32 = maxHealth
	var cargoCount int32 = 0
	var currentCargo int32 = 0
	var planetRadius float32 = 125
	planetColor := rl.NewColor(107, 147, 214, 255)
	planet := NewPlanet(planetPos, float32(maxHealth), planetRadius, planetColor, false)

	// Ship variables
	shipColor := rl.NewColor(255, 255, 255, 255)
	shipSprite := rl.LoadTexture("textures/ship.png")
	shipSpeed := 500
	shipScale := 2
	shipRotation := 0
	ship := NewShip(shipSprite, planetPos, float32(shipRotation), float32(shipScale), float32(shipSpeed), shipColor)

	// Asteroid variables
	var asteroidRadius float32 = 100
	asteroidSpawnRange := planetRadius + float32(rl.GetRandomValue(300, 700))
	randomAsteroidColors := []rl.Color{rl.NewColor(27, 37, 43, 255), rl.NewColor(179, 139, 109, 255), rl.NewColor(133, 127, 114, 255), rl.NewColor(54, 69, 79, 255)}
	var asteroidSpeed float32 = 50
	asteroids := make([]Asteroid, 0, 10)

	// Projectile variables
	var projectiles []Projectile
	var projectileRadius float32 = 10

	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(bgMusic)
		rl.BeginDrawing()

		// Handle Ship Movement
		if rl.IsKeyDown(rl.KeyA) {
			ship.Pos.X -= float32(shipSpeed) * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyD) {
			ship.Pos.X += float32(shipSpeed) * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyW) {
			ship.Pos.Y -= float32(shipSpeed) * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyS) {
			ship.Pos.Y += float32(shipSpeed) * rl.GetFrameTime()
		}

		// Handle Rotation
		if rl.IsKeyDown(rl.KeyQ) {
			ship.Rotation -= rotationSpeed * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyR) {
			ship.Rotation += rotationSpeed * rl.GetFrameTime()
		}

		// Handle projectiles
		if rl.IsKeyPressed(rl.KeySpace) {
			projectileVelocity := rl.NewVector2(float32(math.Sin(float64(rl.Deg2rad*ship.Rotation)))*500, -float32(math.Cos(float64(rl.Deg2rad*ship.Rotation)))*500)
			newProjectile := NewProjectile(ship.Pos, projectileVelocity, projectileRadius)
			projectiles = append(projectiles, newProjectile)
			rl.PlaySound(laserShot)
		}

		// Handle movement of projectiles and asteroids
		for i := range projectiles {
			projectiles[i].UpdateProjectile()

		}
		DestroyProjectile(&projectiles, 2000)

		for i := range asteroids {
			asteroids[i].MoveAsteroid(planetPos)
		}

		// Add asteroids to the slice
		if len(asteroids) < 10 {
			for i := 0; i < 10-len(asteroids); i++ {
				var asteroidPos rl.Vector2
				for {
					asteroidPos = rl.NewVector2(planetPos.X+float32(rl.GetRandomValue(-1000, 1000)), planetPos.Y+float32(rl.GetRandomValue(-1000, 1000)))
					distance := rl.Vector2Distance(planetPos, asteroidPos)

					if distance > asteroidSpawnRange {
						break
					}
				}
				randomIndex := rl.GetRandomValue(0, int32(len(randomAsteroidColors)-1))
				randomColor := randomAsteroidColors[randomIndex]
				asteroidVel := rl.NewVector2(10, 10)
				asteroids = append(asteroids, NewAsteroid(asteroidPos, asteroidVel, asteroidRadius, asteroidSpeed, randomColor))
			}
		}

		rl.ClearBackground(rl.Black)
		cam.Target = ship.Pos

		if gameInProgress {
			rl.BeginMode2D(cam)
			planet.DrawPlanet()

			for _, asteroid := range asteroids {
				asteroid.DrawAsteroid()
			}

			for _, projectile := range projectiles {
				projectile.DrawProjectile()
			}

			ship.DrawShip()
			rl.EndMode2D()
			planetHealth := fmt.Sprintf("Planet Health: %d", currentHealth)
			planetCargo := fmt.Sprintf("Cargo Count: %d", cargoCount)
			rl.DrawText(planetHealth, 5, 5, 20, rl.LightGray)
			rl.DrawText(planetCargo, 5, 25, 20, rl.LightGray)

			// Destroy asteroid and projectiles upon collision and handle heal/cargo scores
			for i := len(asteroids) - 1; i >= 0; i-- {
				if ship.CollectCargo(&asteroids[i]) {
					currentCargo += 1
					rl.PlaySound(cargoPickUp)
					asteroids = append(asteroids[:i], asteroids[i+1:]...)
					continue
				}
				if ship.DepositeCargo(&planet) && currentCargo != 0 {
					cargoCount += currentCargo
					rl.PlaySound(heal)
					currentHealth += currentCargo
					if currentHealth > maxHealth {
						currentHealth = maxHealth
					}
					currentCargo = 0
				}
				if planet.CollisionWithAsteroid(&asteroids[i]) {
					currentHealth -= 1
					asteroids = append(asteroids[:i], asteroids[i+1:]...)
					rl.PlaySound(takeDamage)

					if currentHealth <= 0 {
						gameInProgress = false
					}
					continue
				}
				for j := len(projectiles) - 1; j >= 0; j-- {
					if projectiles[j].DestroyAsteroid(&asteroids[i]) {
						if asteroids[i].Radius > 20 {
							asteroid1, asteroid2 := SplitAsteroids(asteroids[i])
							asteroids = append(asteroids, asteroid1, asteroid2)
						}
						projectiles = append(projectiles[:j], projectiles[j+1:]...)
						rl.PlaySound(explosion)

						asteroids = append(asteroids[:i], asteroids[i+1:]...)
						break
					}
				}
			}
		} else {
			rl.DrawRectangle(0, 0, 1920, 1080, rl.Red)
			rl.DrawText("Game Over! Press R to restart.", 600, 400, 20, rl.White)
			if rl.IsKeyPressed(rl.KeyR) {
				gameInProgress = true

				// Reset health and cargo
				currentHealth = maxHealth
				cargoCount = 0
				currentCargo = 0

				// Reset ship
				ship.Pos = planetPos
				ship.Rotation = 0

				// Clear exisiting asteroids and projectiles
				asteroids = make([]Asteroid, 0, 10)
				projectiles = make([]Projectile, 0)

				// Set the camera
				cam.Target = ship.Pos

				// Reset music
				rl.StopMusicStream(bgMusic)
				rl.PlayMusicStream(bgMusic)
			}
		}
		rl.EndDrawing()
	}
}
