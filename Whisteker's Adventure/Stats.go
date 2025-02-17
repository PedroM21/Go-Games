package main

import (
	"fmt"
	"math/rand/v2"
)

type Stats struct {
	Health        int
	Mana          int
	Strength      int
	Magic         int
	Defense       int
	Resistance    int
	MaxHealth     int
	MaxMana       int
	MaxStrength   int
	MaxMagic      int
	MaxDefense    int
	MaxResistance int
	Level         int
	Experience    float32
}

func (p *Player) LevelUp() bool {
	if p.Experience >= 100 {
		p.Level++
		p.Experience -= 100

		// Random stat growth
		p.Stats.Health += rand.IntN(5) + 5
		p.Stats.Mana += rand.IntN(2) + 2
		p.Stats.Strength += rand.IntN(5) + 5
		p.Stats.Magic += rand.IntN(6) + 4
		p.Stats.Defense += rand.IntN(2) + 3
		p.Stats.Resistance += rand.IntN(3) + 4
		p.Stats.Mana = p.Stats.MaxMana
		fmt.Println("You have leveled up")

		return true
	}

	return false
}

func (e *Enemy) ScaleUp(deafeatedCount int) bool {
	if deafeatedCount%3 == 0 {
		e.Stats.Health += rand.IntN(3) + 4
		e.Stats.Strength += rand.IntN(3) + 4
		e.Stats.Magic += rand.IntN(5) + 3
		e.Stats.Defense += rand.IntN(4) + 4
		e.Stats.Resistance += rand.IntN(4) + 5
		fmt.Println("Enemies have scaled up")

		return true
	}

	return false
}
