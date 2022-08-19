package game

import "fmt"

type Weapon struct {
	Type           string
	ActionVerb     string
	GetKillMessage func(aliveKnight *Knight, deadKnight *Knight) string
}

var Spear = &Weapon {
	Type: "spear",
	ActionVerb: "piercer",
	GetKillMessage: func(aliveKnight *Knight, deadKnight *Knight) string {
		return fmt.Sprintf(
			"%s impaled %s",
			aliveKnight.GetTitle(), deadKnight.GetTitle(),
		)
	},
}

var Sword = &Weapon {
	Type: "sword",
	ActionVerb: "slayer",
	GetKillMessage: func(aliveKnight *Knight, deadKnight *Knight) string {
		return fmt.Sprintf(
			"%s pierced %s's heart",
			aliveKnight.GetTitle(), deadKnight.GetTitle(),
		)
	},
}

var Hammer = &Weapon {
	Type: "war hammer",
	ActionVerb: "crusher",
	GetKillMessage: func(aliveKnight *Knight, deadKnight *Knight) string {
		return fmt.Sprintf(
			"%s caved in %s's chest",
			aliveKnight.GetTitle(), deadKnight.GetTitle(),
		)
	},
}

var AllWeapons = []*Weapon{
	Spear, Sword, Hammer,
}