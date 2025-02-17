package main

import (
	"math/rand/v2"
)

type Enemy struct {
	Stats
	Animation
	CurrentState  string
	IdleAnimation Animation
}

func NewEnemy(idleAnim Animation) Enemy {
	health := rand.IntN(9) + 8

	newEnemy := Enemy{
		Stats: Stats{
			Health:     health,
			MaxHealth:  health,
			Strength:   rand.IntN(3) + 4,
			Magic:      rand.IntN(6) + 5,
			Defense:    rand.IntN(3) + 3,
			Resistance: rand.IntN(4) + 4,
		},
		Animation:    idleAnim,
		CurrentState: "idle",
	}

	newEnemy.IdleAnimation = idleAnim

	return newEnemy
}

func (e *Enemy) Draw() {
	e.DrawAnimation()
}

func (e *Enemy) Damage(target *Player) int {
	damage := e.Stats.Strength - target.Stats.Defense
	if damage <= 0 {
		damage = 1
	}
	target.Stats.Health -= damage
	if target.Stats.Health <= 0 {
		target.Stats.Health = 0
	}

	return damage
}

func (e *Enemy) MagicDamage(target *Player) int {
	damage := e.Stats.Magic - target.Stats.Resistance
	if damage <= 0 {
		damage = 1
	}
	target.Stats.Health -= damage
	if target.Stats.Health <= 0 {
		target.Stats.Health = 0
	}

	return damage
}
