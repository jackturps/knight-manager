package main

import (
	"bufio"
	"fmt"
	"knightmanager/names"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

var MaxMight = 5

type House struct {
	Name string
	Banner string
	/**
	 * Make these values getters so they can be derived from other systems in the future
	 * if necessary.
	 */
	Might int
	Knights []*Knight
}

type BattleResult = int

const (
	Victory BattleResult = iota
	Defeat
)

type Knight struct {
	Name string

	Prowess int
	// TODO: Maybe rename to brashness?
	Bravery int

	BattleResults []BattleResult

	House   *House
	Sponsor *GloryBishop
}

func NewKnight(name string, prowess int, bravery int, house *House, sponsor *GloryBishop) *Knight {
	knight := &Knight{
		Name:          name,
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
	memoryLength := int(math.Min(5, float64(len(knight.BattleResults))))
	if memoryLength == 0 {
		return 1
	}

	// TODO: Make this less swingy from turn 0 to turn 1. Maybe assume an average history if less than max?
	numRecentVictories := 0
	for idx := 0; idx < memoryLength; idx++ {
		battleIdx := len(knight.BattleResults) - (idx + 1)
		battleResult := knight.BattleResults[battleIdx]
		if battleResult == Victory {
			numRecentVictories++
		}
	}

	/**
	 * If a knight won all of their recent battles recentOpinion will be 2. If they won
	 * 0 it will be 0.5. If they won half it will 1.
	 */
	recentReputation := float64(numRecentVictories * 2) / float64(memoryLength)
	return recentReputation
}

func (knight *Knight) GetCost() int {
	underlyingValue := knight.Prowess * knight.House.Might
	return int(float64(underlyingValue) * knight.GetRecentReputation())
}

func (knight *Knight) GetTitle() string {
	return fmt.Sprintf("Lady %s %s", knight.Name, knight.House.Name)
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

	NameGenerator names.NameGenerator
}

func RemoveItem[V comparable](list []V, item V) []V {
	// Copy to prevent in place modification of input slice. Turns out append modifies!
	listCopy := make([]V, len(list))
	copy(listCopy, list)

	for idx, other := range list {
		if other == item {
			return append(listCopy[:idx], listCopy[idx+1:]...)
		}
	}
	return listCopy
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

func GenerateWorld() {
	numHouses := 5
	numKnights := 10

	gameState.Houses = make([]*House, 0, 5)
	for idx := 0; idx < numHouses; idx++ {
		house := &House{
			Name: gameState.NameGenerator.GenerateName(),
			Banner: GenerateBanner(),
			Might: RandomRange(1, MaxMight + 1),
		}
		gameState.Houses = append(gameState.Houses, house)
	}

	gameState.Knights = make([]*Knight, 0, numKnights)
	for idx := 0; idx < numKnights; idx++ {
		GenerateKnight()
	}
}

func GenerateKnight() {
	knight := NewKnight(
		gameState.NameGenerator.GenerateName(),
		RandomRange(1, 6),
		RandomRange(1, 6),
		RandomSelect(gameState.Houses),
		nil,
	)
	// TODO: Should go in knight constructor?
	gameState.Knights = append(gameState.Knights, knight)
}

func GenerateBanner() string {
	// TODO: Move these to input files or something.
	symbols := []string{
		"stag", "wolf", "crab", "crow", "lion", "elephant", "snake", "cross", "heart", "arrow", "ship", "rose", "sword",
		"hanged man", "wheel", "octopus", "horse", "star", "fist", "sunrise", "star", "crescent moon", "beaver",
		"sparrow", "eagle", "chain", "spear", "shield", "apple", "raindrop", "cloud", "lightning bolt", "crystal",
	}
	colours := []string {
		"crimson", "aqua", "light grey", "dark grey", "black", "white", "pink", "golden", "yellow", "blue", "red",
		"purple", "turquoise", "amber", "violet", "orange", "navy", "magenta", "silver", "copper", "teal", "green",
	}
	adjectives := []string {
		"flaming", "submerged", "bloody", "crowned", "upside down", "striped", "spotted", "mirrored",
	}

	// TODO: This could be made simpler with a tracery grammar.
	// Combine parts of banner.
	symbol := RandomSelect(symbols)
	colour := RandomSelect(colours)
	shouldUseAdjective := RandomRange(0, 5) == 0
	var banner string
	if shouldUseAdjective {
		adjective := RandomSelect(adjectives)
		banner = fmt.Sprintf("%s %s %s", adjective, colour, symbol)
	} else {
		banner = fmt.Sprintf("%s %s", colour, symbol)
	}

	startsWithVowel := strings.Contains("aeiou", banner[0:1])
	if startsWithVowel {
		banner = fmt.Sprintf("an %s", banner)
	} else {
		banner = fmt.Sprintf("a %s", banner)
	}
	return banner
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

func ChooseHouseChampion(house *House) *Knight {
	/**
	 * Choose a champion for the house by rolling the bravery of all
	 * knights and choosing the bravest. Prowess is used as a tie
	 * breaker.
	 */
	maxBraveryHits := -1
	var bravestKnight *Knight = nil
	for _, knight := range house.Knights {
		braveryHits := RollHits(knight.Bravery)
		if braveryHits > maxBraveryHits {
			bravestKnight = knight
		} else if braveryHits == maxBraveryHits {
			if knight.Prowess > bravestKnight.Prowess {
				bravestKnight = knight
			}
		}
	}
	return bravestKnight
}

func RunBattle(houses []*House) {
	attackingHouse := RandomSelect(houses)
	possibleTargets := RemoveItem(houses, attackingHouse)
	defendingHouse := RandomSelect(possibleTargets)
	fmt.Printf("House %s attacks House %s!\n", attackingHouse.Name, defendingHouse.Name)


	attackerAdvantage := 0
	defenderAdvantage := 0

	// NOTE: Maybe battles could have multiple "fronts" and we'd have a champion per front.
	// number of fronts could depend on the terrain or some other factor?
	attackingChampion := ChooseHouseChampion(attackingHouse)
	defendingChampion := ChooseHouseChampion(defendingHouse)

	if attackingChampion == nil && defendingChampion == nil {
		fmt.Printf("Neither house could field a champion!\n")
	} else if attackingChampion == nil {
		defenderAdvantage = 1
		fmt.Printf("House %s could not field a champion, giving House %s a tactical edge!\n", attackingHouse.Name, defendingHouse.Name)
	} else if defendingChampion == nil {
		attackerAdvantage = 1
		fmt.Printf("House %s could not field a champion, giving House %s a tactical edge!\n", defendingHouse.Name, attackingHouse.Name)
	} else {
		attackerHits := RollHits(attackingChampion.Prowess)
		defenderHits := RollHits(defendingChampion.Prowess)
		// TODO: Maybe only kill if the margin is big enough.
		// TODO: Award glory to victorious champions.
		if attackerHits == defenderHits {
			fmt.Printf(
				"%s met %s on the battlefield, their duel raged until it met a stalemate[%d/%d vs %d/%d]!\n",
				attackingChampion.GetTitle(), defendingChampion.GetTitle(),
				attackerHits, attackingChampion.Prowess,
				defenderHits, defendingChampion.Prowess,
			)
		} else if attackerHits > defenderHits {
			// TODO: Reduce duplication between here and below.
			fmt.Printf(
				"%s slaid %s on the battlefield after an intense duel[%d/%d vs %d/%d], giving house %s a tactical edge!\n",
				attackingChampion.GetTitle(), defendingChampion.GetTitle(),
				attackerHits, attackingChampion.Prowess,
				defenderHits, defendingChampion.Prowess,
				attackingHouse.Name,
			)

			if attackingChampion.Sponsor != nil {
				glory := int(5 * float64(defendingChampion.Prowess) * defendingChampion.GetRecentReputation())
				gameState.Player.Glory += glory
				fmt.Printf("The Church earned %d glory for sponsoring %s. You're sponsorships have earned the church %d glory in total.\n", glory, attackingChampion.GetTitle(), gameState.Player.Glory)
			}

			attackerAdvantage = 1
			KillKnight(defendingChampion)
		} else {
			fmt.Printf(
				"%s slaid %s on the battlefield after an intense duel[%d/%d vs %d/%d], giving house %s a tactical edge!\n",
				defendingChampion.GetTitle(), attackingChampion.GetTitle(),
				defenderHits, defendingChampion.Prowess,
				attackerHits, attackingChampion.Prowess,
				defendingHouse.Name,
			)

			if defendingChampion.Sponsor != nil {
				glory := int(5 * float64(attackingChampion.Prowess) * attackingChampion.GetRecentReputation())
				gameState.Player.Glory += glory
				fmt.Printf("The Church earned %d glory for sponsoring %s. You're sponsorships have earned the church %d glory in total.\n", glory, defendingChampion.GetTitle(), gameState.Player.Glory)
			}

			defenderAdvantage = 1
			KillKnight(attackingChampion)
		}
	}




	attackerHits := RollHits(attackingHouse.Might + attackerAdvantage)
	defenderHits := RollHits(defendingHouse.Might + defenderAdvantage)

	var winner, loser *House
	var winnerHits, loserHits int

	if attackerHits > defenderHits {
		winner, winnerHits, loser, loserHits = attackingHouse, attackerHits, defendingHouse, defenderHits
	} else {
		winner, winnerHits, loser, loserHits = defendingHouse, defenderHits, attackingHouse, attackerHits
	}
	// TODO: Print advantages?
	fmt.Printf(
		"House %s[%d/%d hits] defeated House %s[%d/%d hits]!\n",
		winner.Name, winnerHits, winner.Might, loser.Name, loserHits, loser.Might,
	)

	// Award more glory to underdogs and less to bullies.
	glory := (MaxMight + 1) + (loser.Might - winner.Might)
	for _, knight := range winner.Knights {
		knight.BattleResults = append(knight.BattleResults, Victory)
		if knight.Sponsor != nil {
			gameState.Player.Glory += glory
			fmt.Printf("The Church earned %d glory for sponsoring %s. You're sponsorships have earned the church %d glory in total.\n", glory, knight.GetTitle(), gameState.Player.Glory)
		}
	}
	for _, knight := range loser.Knights {
		knight.BattleResults = append(knight.BattleResults, Defeat)
		// TODO: Maybe a chance of death here?
	}

	fmt.Printf("\n")
}

func DisplayState() {
	for _, house := range gameState.Houses {
		fmt.Printf("Introducing the knights of House %s[might: %d]! Their banner is %s.\n", house.Name, house.Might, house.Banner)
		for _, knight := range house.Knights {

			fmt.Printf("%s! [prowess: %d, bravery: %d, cost: %d]\n", knight.GetTitle(), knight.Prowess, knight.Bravery, knight.GetCost())
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n===================\n\n")
}

func DoPlayerTurn() {
	// Player interaction loop.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		command := strings.Split(input, " ")
		if command[0] == "done" {
			fmt.Printf("\n====================\n\n")
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
				fmt.Printf("%s is already sponsored\n", foundKnight.GetTitle())
				continue
			}

			cost := foundKnight.GetCost()
			if cost > gameState.Player.Coin {
				fmt.Printf("The church coffers run low, %s costs %d coin but you only have %d.\n", foundKnight.GetTitle(), cost, gameState.Player.Coin)
				continue
			}

			gameState.Player.Coin -= foundKnight.GetCost()
			SponsorKnight(gameState.Player, foundKnight)
			fmt.Printf(
				"You have sponsored %s, %d coin remaining\n",
				foundKnight.GetTitle(), gameState.Player.Coin,
			)
		}
	}
}

var gameState *GameState

func main() {
	rand.Seed(time.Now().UnixNano())

	gameState = &GameState{}
	gameState.NameGenerator = names.NewSelectorNameGenerator("input_names.txt")
	GenerateWorld()

	gameState.Player = &GloryBishop{
		Coin: 15,
		Glory: 0,
	}

	for {
		DisplayState()
		DoPlayerTurn()
		numBattles := 3
		for idx := 0; idx < numBattles; idx++ {
			RunBattle(gameState.Houses)
		}
		fmt.Printf("\n=================\n\n")

		numNewKnights := 1
		for idx := 0; idx < numNewKnights; idx++ {
			GenerateKnight()
		}

		gameState.Player.Coin += 5
	}
}
