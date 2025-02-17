package main

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1600, 900, "Final Project")
	rl.InitAudioDevice()
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	screenWidth := 1600
	screenHeight := 900

	mainMenu := true
	gameMenu := false
	fightMenu := false
	gameOver := false
	pauseGame := false

	var damageTexts []DamageText
	var defeatedCount int = 0

	// audio
	bgMusic := rl.LoadMusicStream("audio/rpgcity.ogg")
	buttonSound := rl.LoadSound("audio/toy-button.mp3")
	rl.PlayMusicStream(bgMusic)
	rl.SetSoundVolume(buttonSound, 0.25)
	rl.SetMusicVolume(bgMusic, 0.5)

	// set color theme
	colorTheme := NewColorTheme(
		rl.NewColor(255, 255, 255, 255),
		rl.NewColor(128, 255, 255, 255),
		rl.NewColor(0, 0, 0, 255),
	)

	// Load textures
	playerWalkRight := rl.LoadTexture("sprites/PlayerWR.png")
	playerWalkLeft := rl.LoadTexture("sprites/PlayerWL.png")
	playerIdle := rl.LoadTexture("sprites/Character.png")

	enemySprites := []rl.Texture2D{
		rl.LoadTexture("sprites/squirrel.png"),
		rl.LoadTexture("sprites/crow.png"),
	}

	// Create animations
	idleAnim := NewAnimation(rl.NewVector2(50, 300), playerIdle, 5, .15)
	walkRightAnim := NewAnimation(rl.NewVector2(20, 20), playerWalkRight, 5, .15)
	walkLeftAnim := NewAnimation(rl.NewVector2(20, 20), playerWalkLeft, 5, .15)

	// Randomize enemy spawn
	selectedSprite := enemySprites[rand.IntN(len(enemySprites))]
	enemyIdleAnim := NewAnimation(rl.NewVector2(900, 300), selectedSprite, 5, .15)

	// Create player and enemies
	player := NewPlayer(300, idleAnim, walkRightAnim, walkLeftAnim)
	enemy := NewEnemy(enemyIdleAnim)

	// Create hud elements
	playerHealthBar := NewProgressBar(50, 50, 300, 50, 50, 50, rl.Red)
	playerManaBar := NewProgressBar(50, 110, 200, 50, 50, 50, rl.DarkBlue)
	playerExpBar := NewProgressBar(50, 170, 300, 10, 50, 50, rl.NewColor(8, 189, 8, 255))
	enemyHealthBar := NewProgressBar(1250, 50, 300, 50, 50, 50, rl.Red)

	// Create buttons for main menu
	startButton := NewButton(0, 300, 300, 100, colorTheme)
	quitButton := NewButton(0, 500, 300, 100, colorTheme)
	startButton.SetText("Start Game", 20)
	quitButton.SetText("Quit Game", 20)
	startButton.CenterButtonX()
	quitButton.CenterButtonX()
	mainMenuButton := NewButton(0, 300, 300, 100, colorTheme)
	mainMenuButton.SetText("Main Menu", 20)
	mainMenuButton.CenterButtonX()

	// Buttons for Fight Menu
	attackButton := NewButton(100, 650, 150, 75, colorTheme)
	spellButton := NewButton(100, 750, 150, 75, colorTheme)
	healButton := NewButton(350, 650, 150, 75, colorTheme)
	attackButton.SetText("Attack", 30)
	spellButton.SetText("Magic", 30)
	healButton.SetText("Heal", 30)

	// Functionality
	startButton.AddOnClickFunc(func() {
		mainMenu = false
		gameMenu = true

		player, enemy, damageTexts, playerHealthBar, playerManaBar, playerExpBar, enemyHealthBar = RestartGame()
		defeatedCount = 0
		rl.PlaySound(buttonSound)
		rl.StopMusicStream(bgMusic)
		rl.PlayMusicStream(bgMusic)
	})

	quitButton.AddOnClickFunc(func() {
		rl.PlaySound(buttonSound)
		rl.CloseWindow()
	})

	mainMenuButton.AddOnClickFunc(func() {
		rl.PlaySound(buttonSound)
		mainMenu = true
		pauseGame = false
	})

	attackButton.AddOnClickFunc(func() {
		rl.PlaySound(buttonSound)
		playerDamage := player.Damage(&enemy)
		damageText := NewDamageText(550, 700, "The Enemy", int32(playerDamage))
		damageTexts = append(damageTexts, damageText)

		if enemy.Stats.Health > 0 {
			enemyDamage := enemy.Damage(&player)
			enemyDamageText := NewDamageText(550, 750, "The Player", int32(enemyDamage))
			damageTexts = append(damageTexts, enemyDamageText)
		} else {
			defeatedCount++
			selectedSprite := enemySprites[rand.IntN(len(enemySprites))]
			enemyIdleAnim := NewAnimation(rl.NewVector2(900, 300), selectedSprite, 5, .15)
			enemy = NewEnemy(enemyIdleAnim)

			enemy.ScaleUp(defeatedCount)

			fightMenu = false
			gameMenu = true
			player.Pos = rl.NewVector2(50, 300)
		}

		if player.Stats.Health <= 0 {
			player.Stats.Health = 0
			fightMenu = false
			gameOver = true
		}
	})

	spellButton.AddOnClickFunc(func() {
		rl.PlaySound(buttonSound)
		playerMagicDamage := player.MagicDamage(&enemy)
		mDamageText := NewDamageText(550, 700, "The Enemy", int32(playerMagicDamage))
		damageTexts = append(damageTexts, mDamageText)

		if enemy.Stats.Health > 0 {
			enemyMagicDamage := enemy.MagicDamage(&player)
			enemyMDamageText := NewDamageText(550, 750, "The Player", int32(enemyMagicDamage))
			damageTexts = append(damageTexts, enemyMDamageText)
		} else {
			defeatedCount++
			enemy.ScaleUp(defeatedCount)
			fightMenu = false
			gameMenu = true
			player.Pos = rl.NewVector2(50, 300)
		}

		if player.Stats.Health <= 0 {
			player.Stats.Health = 0
			fightMenu = false
			gameOver = true
		}
	})

	healButton.AddOnClickFunc(func() {
		rl.PlaySound(buttonSound)
		player.HealHealth()
	})

	for !rl.WindowShouldClose() {
		player.Update()

		for i := len(damageTexts) - 1; i >= 0; i-- {
			d := &damageTexts[i]
			d.Update()

			// Remove the text after the duration is over
			if d.Duration <= 0 {
				damageTexts = append(damageTexts[:i], damageTexts[i+1:]...)
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Update hud bars
		healthProgress := float32(player.Stats.Health) / float32(player.Stats.MaxHealth)
		playerHealthBar.SetProgress(float32(healthProgress))
		enemyProgress := float32(enemy.Stats.Health) / float32(enemy.Stats.MaxHealth)
		enemyHealthBar.SetProgress(float32(enemyProgress))

		manaProgress := float32(player.Stats.Mana) / float32(player.Stats.MaxMana)
		playerManaBar.SetProgress(manaProgress)

		expProgress := float32(player.Stats.Experience) / 100.0
		playerExpBar.SetProgress(expProgress)

		// Change menu based on game conditions
		if mainMenu {
			rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.Orange)
			rl.DrawText("MAIN MENU", int32(screenWidth/2)-120, int32(screenHeight/2)-300, 40, rl.White)
			startButton.UpdateButton()
			quitButton.UpdateButton()

		} else if gameMenu {
			rl.UpdateMusicStream(bgMusic)
			rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.Purple)
			score := fmt.Sprintf("Defeated Enemies: %d", defeatedCount)
			rl.DrawText(score, 20, 20, 40, rl.White)
			player.Draw()
			enemy.Draw()

			if player.CheckCollision(&enemy) {
				gameMenu = false
				fightMenu = true
			}

			if rl.IsKeyPressed(rl.KeyP) {
				gameMenu = false
				pauseGame = true
			}

		} else if fightMenu {
			rl.UpdateMusicStream(bgMusic)
			rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.Blue)
			levelText := fmt.Sprintf("Level: %d", player.Stats.Level)
			rl.DrawText(levelText, 50, 10, 40, rl.White)

			playerHealthBar.DrawBar()
			playerManaBar.DrawBar()
			playerExpBar.DrawBar()
			enemyHealthBar.DrawBar()
			attackButton.UpdateButton()
			spellButton.UpdateButton()
			healButton.UpdateButton()

			// Set player and enemy in battle area
			player.Pos = rl.NewVector2(100, 300)
			player.Draw()
			enemy.Pos = rl.NewVector2(1300, 300)
			enemy.Draw()
			fmt.Println("Current Health:", player.Stats.Health)

			for _, d := range damageTexts {
				d.Draw()
			}

			if rl.IsKeyPressed(rl.KeyP) {
				fightMenu = false
				pauseGame = true
			}

		} else if gameOver {
			rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.Red)
			rl.DrawText("GAME OVER", int32(screenWidth/2)-120, int32(screenHeight/2)-300, 40, rl.White)
			finalScore := fmt.Sprintf("Final Score: %d", defeatedCount)
			rl.DrawText(finalScore, int32(screenWidth/2)-120, int32(screenHeight/2)-200, 40, rl.White)
			mainMenuButton.UpdateButton()
			quitButton.UpdateButton()
		} else if pauseGame {
			rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.Orange)
			rl.DrawText("PAUSE MENU", 20, 20, 40, rl.White)
			mainMenuButton.UpdateButton()
			quitButton.UpdateButton()
		}

		rl.EndDrawing()

	}
}

func RestartGame() (Player, Enemy, []DamageText, ProgressBar, ProgressBar, ProgressBar, ProgressBar) {
	damageTexts := []DamageText{}

	playerIdle := rl.LoadTexture("sprites/Character.png")
	playerWalkRight := rl.LoadTexture("sprites/PlayerWR.png")
	playerWalkLeft := rl.LoadTexture("sprites/PlayerWL.png")

	idleAnim := NewAnimation(rl.NewVector2(50, 300), playerIdle, 5, .15)
	walkRightAnim := NewAnimation(rl.NewVector2(20, 20), playerWalkRight, 5, .15)
	walkLeftAnim := NewAnimation(rl.NewVector2(20, 20), playerWalkLeft, 5, .15)
	player := NewPlayer(300, idleAnim, walkRightAnim, walkLeftAnim)

	enemySprites := []rl.Texture2D{
		rl.LoadTexture("sprites/squirrel.png"),
		rl.LoadTexture("sprites/crow.png"),
	}
	selectedSprite := enemySprites[rand.IntN(len(enemySprites))]
	enemyIdleAnim := NewAnimation(rl.NewVector2(900, 300), selectedSprite, 5, .15)
	enemy := NewEnemy(enemyIdleAnim)

	playerHealthBar := NewProgressBar(50, 50, 300, 50, 50, 50, rl.Red)
	playerManaBar := NewProgressBar(50, 110, 200, 50, 50, 50, rl.DarkBlue)
	playerExpBar := NewProgressBar(50, 170, 300, 10, 50, 50, rl.NewColor(8, 189, 8, 255))
	enemyHealthBar := NewProgressBar(1250, 50, 300, 50, 50, 50, rl.Red)

	return player, enemy, damageTexts, playerHealthBar, playerManaBar, playerExpBar, enemyHealthBar
}
