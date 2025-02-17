package main

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Stats
	Animation
	Movespeed    float32
	CurrentState string

	IdleAnimation      Animation
	WalkRightAnimation Animation
	WalkLeftAnimation  Animation
}

func NewPlayer(mspeed float32, idleAnim, walkRightAnim, walkLeftAnim Animation) Player {
	health := rand.IntN(11) + 10
	mana := rand.IntN(5) + 6

	newPlayer := Player{
		Stats: Stats{
			Health:     health,
			MaxHealth:  health,
			Mana:       mana,
			MaxMana:    mana,
			Strength:   rand.IntN(4) + 5,
			Magic:      rand.IntN(7) + 6,
			Defense:    rand.IntN(3) + 3,
			Resistance: rand.IntN(4) + 4,
			Level:      1,
			Experience: 0.0,
		},
		Animation:    idleAnim,
		Movespeed:    mspeed,
		CurrentState: "idle",
	}

	newPlayer.IdleAnimation = idleAnim
	newPlayer.WalkRightAnimation = walkRightAnim
	newPlayer.WalkLeftAnimation = walkLeftAnim

	return newPlayer
}

func (p *Player) Update() {
	switch {
	case rl.IsKeyDown(rl.KeyD):
		p.CurrentState = "walkRight"
		p.Pos.X += p.Movespeed * rl.GetFrameTime()
		p.SpriteSheet = p.WalkRightAnimation.SpriteSheet
		p.UpdateTime()

	case rl.IsKeyDown(rl.KeyA):
		p.CurrentState = "walkLeft"
		p.Pos.X -= p.Movespeed * rl.GetFrameTime()
		p.SpriteSheet = p.WalkLeftAnimation.SpriteSheet
		p.UpdateTime()

	case rl.IsKeyDown(rl.KeyW):
		p.CurrentState = "idle"
		p.Pos.Y -= p.Movespeed * rl.GetFrameTime()
		p.SpriteSheet = p.IdleAnimation.SpriteSheet
		p.UpdateTime()

	case rl.IsKeyDown(rl.KeyS):
		p.CurrentState = "idle"
		p.Pos.Y += p.Movespeed * rl.GetFrameTime()
		p.SpriteSheet = p.IdleAnimation.SpriteSheet
		p.UpdateTime()

	default:
		p.CurrentState = "idle"
		p.SpriteSheet = p.IdleAnimation.SpriteSheet
		p.UpdateTime()
	}
}

func (p *Player) CheckCollision(e *Enemy) bool {
	playerRect := rl.Rectangle{
		X:      p.Pos.X,
		Y:      p.Pos.Y,
		Width:  float32(p.SpriteSheet.Width) / float32(p.MaxIndex),
		Height: float32(p.SpriteSheet.Height),
	}

	enemyRect := rl.Rectangle{
		X:      e.Pos.X,
		Y:      e.Pos.Y,
		Width:  float32(e.SpriteSheet.Width) / float32(e.MaxIndex),
		Height: float32(e.SpriteSheet.Height),
	}

	return rl.CheckCollisionRecs(playerRect, enemyRect)
}

func (p *Player) Draw() {
	p.DrawAnimation()
}

func (p *Player) Damage(target *Enemy) int {
	damage := p.Stats.Strength - target.Stats.Defense
	if damage <= 0 {
		damage = 1
	}
	target.Stats.Health -= damage
	if target.Stats.Health <= 0 {
		target.Stats.Health = 0
		p.Experience += float32(rand.IntN(46) + 30)
		p.LevelUp()

	}

	return damage
}

func (p *Player) MagicDamage(target *Enemy) int {
	spellCost := 4
	if p.Stats.Mana < spellCost {
		return 0
	}

	damage := p.Stats.Magic - target.Stats.Resistance
	if damage <= 0 {
		damage = 1
	}
	target.Stats.Health -= damage
	if target.Stats.Health <= 0 {
		target.Stats.Health = 0
		p.Experience += float32(rand.IntN(46) + 30)
		p.LevelUp()
	}
	p.Stats.Mana -= spellCost

	return damage
}

func (p *Player) HealHealth() int {
	healCost := 3
	healAmount := p.Stats.MaxHealth / 4

	if p.Stats.Mana < healCost {
		return 0
	}

	p.Stats.Health += healAmount
	if p.Stats.Health > p.Stats.MaxHealth {
		p.Stats.Health = p.Stats.MaxHealth
	}
	p.Stats.Mana -= healCost

	return healAmount
}
