package main

import (
	"bufio"
	"fmt"
	"knightmanager/names"
	"math/rand"
	"os"
	"strings"
	"time"
)

var MaxMight = 5

type House struct {
	Name string
	/**
	 * Make these values getters so they can be derived from other systems in the future
	 * if necessary.
	 */
	Might int
	Knights []*Knight
}

type Knight struct {
	Name string

	Prowess int

	// This would be in quotes if it could be, mostly represents brashness.
	Bravery int

	House   *House
	Sponsor *GloryBishop
}

// GloryBishop is a member of the church who sponsors knights for glory.
// the player will represent a glory bishop.
type GloryBishop struct {
	Coin int
	Glory int

	SponsoredKnights []*Knight
}

type GameState struct {
	Player *GloryBishop
	Knights []*Knight
	Houses []*House
}

func RemoveItem[V comparable](list []V, item V) []V {
	for idx, other := range list {
		if other == item {
			return append(list[:idx], list[idx+1:]...)
		}
	}
	return list
}

func RandomSelect[V any] (values []V) V {
	return values[rand.Intn(len(values))]
}

// RandomRange generates a random number between min and max. This includes min but excludes max.
func RandomRange(min int, max int) int {
	if min > max {
		panic(fmt.Sprintf("min(%d) must be less than or equal to max(%d)", min, max))
	}
	return rand.Intn(max - min) + min
}

func AssignKnightToHouse(knight *Knight, house *House) {
	house.Knights = append(house.Knights, knight)
	knight.House = house
}

func SponsorKnight(bishop *GloryBishop, knight *Knight) {
	bishop.SponsoredKnights = append(bishop.SponsoredKnights, knight)
	knight.Sponsor = bishop
}

func GenerateWorld() ([]*House, []*Knight) {
	numHouses := 5
	numKnights := 10

	var nameGenerator names.NameGenerator = names.NewSelectorNameGenerator("input_names.txt")


	houses := make([]*House, 0, 5)
	for idx := 0; idx < numHouses; idx++ {
		house := &House{
			Name: nameGenerator.GenerateName(),
			Might: RandomRange(1, MaxMight + 1),
		}
		houses = append(houses, house)
	}

	knights := make([]*Knight, 0, numKnights)
	for idx := 0; idx < numKnights; idx++ {
		knight := &Knight{
			Name: nameGenerator.GenerateName(),
			Prowess: RandomRange(1, 5),
			Bravery: RandomRange(1, 5),
		}

		assignedHouse := RandomSelect(houses)
		AssignKnightToHouse(knight, assignedHouse)
		knights = append(knights, knight)
	}

	return houses, knights
}

// Given a certain rating randomly determine the number of success. Effectively
// a dice pool system.
func RollHits(rating int) int {
	successes := 0
	// TODO: There's probably a way to do this in a single call to some probability curve.
	/**
	 * Roll a d6. If its greater than 3 it counts as a success. If its a 6 it counts
	 * as a success and can be rerolled.
	 */
	for idx := 0; idx < rating; idx++ {
		value := RandomRange(1, 7)
		if value >= 4 {
			successes++
		}
		if value == 6 {
			rating++
		}
	}
	return successes
}

func KillKnight(knight *Knight) {
	knight.House.Knights = RemoveItem(knight.House.Knights, knight)
	if knight.Sponsor != nil {
		knight.Sponsor.SponsoredKnights = RemoveItem(knight.Sponsor.SponsoredKnights, knight)
	}
	gameState.Knights = RemoveItem(gameState.Knights, knight)
}

func RunBattle(houses []*House) {
	attackingHouse := RandomSelect(houses)
	possibleTargets := RemoveItem(houses, attackingHouse)
	defendingHouse := RandomSelect(possibleTargets)
	fmt.Printf("House %s attacks House %s!\n", attackingHouse.Name, defendingHouse.Name)

	attackerHits := RollHits(attackingHouse.Might)
	defenderHits := RollHits(defendingHouse.Might)

	var winner, loser *House
	var winnerHits, loserHits int

	if attackerHits > defenderHits {
		winner, winnerHits, loser, loserHits = attackingHouse, attackerHits, defendingHouse, defenderHits
	} else {
		winner, winnerHits, loser, loserHits = defendingHouse, defenderHits, attackingHouse, attackerHits
	}
	fmt.Printf(
		"House %s[%d/%d hits] defeated House %s[%d/%d hits]!\n",
		winner.Name, winnerHits, winner.Might, loser.Name, loserHits, loser.Might,
	)

	// Award more glory to underdogs and less to bullies.
	glory := (MaxMight + 1) + (loser.Might - winner.Might)
	for _, knight := range winner.Knights {
		fmt.Printf(
			"Ser %s was awarded %d glory for their deeds in battle!\n", knight.Name, glory)
		if knight.Sponsor != nil {
			gameState.Player.Glory += glory
			fmt.Printf("The Church earned %d glory for sponsoring Ser %s. You're sponsorships have earned the church %d glory in total.\n", glory, knight.Name, gameState.Player.Glory)
		}
	}

	lossSeverity := winnerHits - loserHits
	// Copy the slice as we may remove items from the primary slice so we can't iterate it.
	loserKnights := loser.Knights
	for _, knight := range loserKnights {
		survivalHits := RollHits(knight.Prowess)
		if survivalHits < lossSeverity {
			fmt.Printf(
				"Ser %s fought valiantly[%d/%d hits] but was slain by the enemy forces[%d].\n",
				knight.Name, survivalHits, knight.Prowess, lossSeverity,
			)
			KillKnight(knight)
		}
	}

	fmt.Printf("\n")
}

var gameState *GameState

func main() {
	rand.Seed(time.Now().UnixNano())

	gameState = &GameState{}
	gameState.Houses, gameState.Knights = GenerateWorld()
	for _, house := range gameState.Houses {
		fmt.Printf("Introducing the knights of House %s! [might: %d]\n", house.Name, house.Might)
		for _, knight := range house.Knights {

			fmt.Printf("Ser %s of House %s! [prowess: %d]\n", knight.Name, knight.House.Name, knight.Prowess)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n===================\n\n")

	gameState.Player = &GloryBishop{
		Coin: 5,
		Glory: 0,
	}

	// Player interaction loop.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		command := strings.Split(input, " ")
		if command[0] == "done" {
			break
		} else if command[0] == "sponsor" {
			knightName := command[1]
			var foundKnight *Knight = nil
			for _, knight := range gameState.Knights {
				if knight.Name == knightName {
					foundKnight = knight
					break
				}
			}
			if foundKnight == nil {
				fmt.Printf("Could not find knight '%s'\n", knightName)
				continue
			}
			if foundKnight.Sponsor != nil {
				fmt.Printf("Ser %s is already sponsored\n", foundKnight.Name)
			}

			gameState.Player.Coin--
			SponsorKnight(gameState.Player, foundKnight)
			fmt.Printf(
				"You have sponsored Ser %s of House %s, %d coin remaining\n",
				foundKnight.Name, foundKnight.House.Name, gameState.Player.Coin,
			)
		}
	}

	numBattles := 10
	for idx := 0; idx < numBattles; idx++ {
		RunBattle(gameState.Houses)
	}
}
