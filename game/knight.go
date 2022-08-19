package game

import (
	"fmt"
	"math"
)

type BattleResult = int
const (
	Victory BattleResult = iota
	Defeat
)

type Gender = int
const (
	Male Gender = iota
	Female
)

type Knight struct {
	Name string
	Gender Gender

	Prowess int
	// TODO: Maybe rename to brashness?
	Bravery int
	Weapon *Weapon

	Spouse   *Knight

	Blessings int
	IsChosen bool

	BattleResults []BattleResult

	House   *House
	Sponsor *GloryBishop
}

func NewKnight(name string, gender Gender, prowess int, bravery int, weapon *Weapon, house *House, sponsor *GloryBishop) *Knight {
	knight := &Knight{
		Name:          name,
		Gender:        gender,
		Prowess:       prowess,
		Bravery:       bravery,
		Weapon:        weapon,
		Spouse:        nil,
		IsChosen:      false,
		BattleResults: make([]BattleResult, 0),
		House:         house,
		Sponsor:       sponsor,
	}
	if house != nil {
		AssignKnightToHouse(knight, house)
	}
	return knight
}

func KillKnight(knight *Knight) {
	if knight.Sponsor != nil {
		titheAmount := 5 * knight.House.Wealth
		fmt.Printf(
			"%s paid %d coin in customary funeral tithes for %s.\n",
			knight.House.GetTitle(), titheAmount, knight.GetTitle(),
		)
		knight.Sponsor.Coin += titheAmount
	}

	// Make their spouse a widow :(.
	if knight.Spouse != nil {
		fmt.Printf("%s was made a widow.\n", knight.Spouse.GetTitle())
		knight.Spouse.Spouse = nil
	}

	knight.House.Knights = RemoveItem(knight.House.Knights, knight)
	if knight.Sponsor != nil {
		knight.Sponsor.SponsoredKnights = RemoveItem(knight.Sponsor.SponsoredKnights, knight)
	}
	Game.Knights = RemoveItem(Game.Knights, knight)
}

// GetRecentReputation returns a knights reputation based on
// their recent accomplishments. This will be a float, 0-1 is a
// bad reputation, 1+ is a good reputation.
func (knight *Knight) GetRecentReputation() float64 {
	maxMemoryLength := 5

	memoryLength := int(math.Min(float64(maxMemoryLength), float64(len(knight.BattleResults))))
	if memoryLength == 0 {
		return 1
	}

	// Fill out any missing battles with average results.
	numRecentVictories := (maxMemoryLength -  memoryLength) / 2
	memoryStartIdx := len(knight.BattleResults) - memoryLength
	for _, battleResult := range knight.BattleResults[memoryStartIdx:] {
		if battleResult == Victory {
			numRecentVictories++
		}
	}

	/**
	 * If a knight won all of their recent battles recentOpinion will be 2. If they won
	 * 0 it will be 0.5. If they won half it will 1.
	 */
	recentReputation := 0.5 + ((float64(numRecentVictories) * 1.5) / float64(maxMemoryLength))
	return recentReputation
}

func (knight *Knight) GetCost() int {
	underlyingValue := knight.Prowess * knight.House.Might
	// TODO: Maybe adjust the math so we can't go below 1 without need a min?
	return 1 + int(float64(underlyingValue) * knight.GetRecentReputation())
}

func (knight *Knight) GetTitle() string {
	var genderedTitle string
	if knight.Gender == Male {
		genderedTitle = "Ser"
	} else if knight.Gender == Female {
		genderedTitle = "Lady"
	}

	blessedTitle := ""
	if knight.Blessings > 0 {
		blessedTitle = " the Blessed"
	}

	title := fmt.Sprintf("%s %s %s%s", genderedTitle, knight.Name, knight.House.Name, blessedTitle)
	if knight.Sponsor != nil {
		title = ColouredText(GreenTextCode, title)
	}
	if knight.IsChosen {
		title = ColouredText(BlueBackgroundCode, title)
	}
	return title
}
