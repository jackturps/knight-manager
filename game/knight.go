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

	BattleResults []BattleResult

	House   *House
	Sponsor *GloryBishop
}

func NewKnight(name string, gender Gender, prowess int, bravery int, house *House, sponsor *GloryBishop) *Knight {
	knight := &Knight{
		Name:          name,
		// TODO: Get a list of male names and make it possible to be a man.
		Gender:		   gender,
		Prowess:       prowess,
		Bravery:	   bravery,
		BattleResults: make([]BattleResult, 0),
		House:         house,
		Sponsor:       sponsor,
	}
	if house != nil {
		AssignKnightToHouse(knight, house)
	}
	return knight
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
	return fmt.Sprintf("%s %s %s", genderedTitle, knight.Name, knight.House.Name)
}
