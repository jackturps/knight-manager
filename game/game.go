package game

import (
	"fmt"
	"knightmanager/names"
	"strings"
)

var MaxMight = 5

type DiplomaticRelation struct {
	Tension int
}

type House struct {
	Name string
	Banner string

	/**
	 * Make these values getters so they can be derived from other systems in the future
	 * if necessary.
	 */
	Might int

	Knights []*Knight
	DiplomaticRelations map[*House]*DiplomaticRelation
}

func (house *House) GetTitle() string {
	return fmt.Sprintf("House %s", house.Name)
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
	Wars []*War

	FemaleNameGenerator names.NameGenerator
	MaleNameGenerator names.NameGenerator
}

func AssignKnightToHouse(knight *Knight, house *House) {
	house.Knights = append(house.Knights, knight)
	knight.House = house
}

func SponsorKnight(bishop *GloryBishop, knight *Knight) {
	bishop.SponsoredKnights = append(bishop.SponsoredKnights, knight)
	knight.Sponsor = bishop
}

func RemoveHouse(house *House) {
	for _, knight := range house.Knights {
		Game.Knights = RemoveItem(Game.Knights, knight)
	}
	Game.Houses = RemoveItem(Game.Houses, house)
}

func GenerateHouse() *House {
	house := &House{
		Name:   Game.FemaleNameGenerator.GenerateName(),
		Banner: GenerateBanner(),
		Might:  RandomRange(1, MaxMight + 1),
		DiplomaticRelations: make(map[*House]*DiplomaticRelation, 0),
	}
	Game.Houses = append(Game.Houses, house)
	InitNewDiplomaticRelations()
	return house
}

func InitNewDiplomaticRelations() {
	for _, srcHouse := range Game.Houses {
		for _, dstHouse := range Game.Houses {
			// Can't have relation with your own house.
			if srcHouse == dstHouse {
				continue
			}
			// Don't overwrite existing relationship.
			if _, houseFound := srcHouse.DiplomaticRelations[dstHouse]; houseFound {
				continue
			}

			srcHouse.DiplomaticRelations[dstHouse] = &DiplomaticRelation{
				Tension: 0,
			}
		}
	}
}

func GenerateWorld() {
	numHouses := 6
	numKnights := 10

	// Generate houses.
	Game.Houses = make([]*House, 0, 5)
	for idx := 0; idx < numHouses; idx++ {
		GenerateHouse()
	}

	// Generate knights.
	Game.Knights = make([]*Knight, 0, numKnights)
	for idx := 0; idx < numKnights; idx++ {
		GenerateKnight(RandomSelect(Game.Houses))
	}
}

func GenerateKnight(house *House) {
	gender := RandomSelect([]Gender{Female, Male})
	var name string
	if gender == Female {
		name = Game.FemaleNameGenerator.GenerateName()
	} else {
		name = Game.MaleNameGenerator.GenerateName()
	}

	knight := NewKnight(
		name, gender,
		RandomRange(1, 6), RandomRange(1, 6),
		house,
		nil,
	)
	// TODO: Should go in knight constructor?
	Game.Knights = append(Game.Knights, knight)
}

func GenerateBanner() string {
	// TODO: Move these to input files or something.
	symbols := []string{
		"stag", "wolf", "crab", "crow", "lion", "elephant", "snake", "cross", "heart", "arrow", "ship", "rose", "sword",
		"hanged man", "wheel", "octopus", "horse", "star", "fist", "sunrise", "star", "crescent moon", "beaver",
		"sparrow", "eagle", "chain", "spear", "shield", "apple", "raindrop", "cloud", "lightning bolt", "crystal",
		"demon", "angel", "dragon", "griffin", "unicorn", "hydra", "bull", "goat", "sheep",
	}
	colours := []string {
		"crimson", "aqua", "light grey", "dark grey", "black", "white", "pink", "golden", "yellow", "blue", "red",
		"purple", "turquoise", "amber", "violet", "orange", "navy", "magenta", "silver", "copper", "teal", "green",
	}
	adjectives := []string {
		"flaming", "submerged", "bloody", "crowned", "upside down", "striped", "spotted", "mirrored", "frozen",
		"shattered",
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
	Game.Knights = RemoveItem(Game.Knights, knight)
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
			maxBraveryHits = braveryHits
			bravestKnight = knight
		} else if braveryHits == maxBraveryHits {
			if knight.Prowess > bravestKnight.Prowess {
				bravestKnight = knight
			}
		}
	}
	return bravestKnight
}

// Returns the margin of the attacker. This will be <=0 if they lost
// and >0 if they won.
func RunBattle(attackingHouse *House, defendingHouse *House) int {
	// TODO: Reduce morale for every knight killed?
	fmt.Printf("%s attacks %s!\n", attackingHouse.GetTitle(), defendingHouse.GetTitle())

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
		fmt.Printf("%s could not field a champion, giving %s a tactical edge!\n", attackingHouse.GetTitle(), defendingHouse.GetTitle())
	} else if defendingChampion == nil {
		attackerAdvantage = 1
		fmt.Printf("%s could not field a champion, giving %s a tactical edge!\n", defendingHouse.GetTitle(), attackingHouse.GetTitle())
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
				"%s slayed %s on the battlefield after an intense duel[%d/%d vs %d/%d], giving %s a tactical edge!\n",
				attackingChampion.GetTitle(), defendingChampion.GetTitle(),
				attackerHits, attackingChampion.Prowess,
				defenderHits, defendingChampion.Prowess,
				attackingHouse.GetTitle(),
			)

			if attackingChampion.Sponsor != nil {
				glory := int(5 * float64(defendingChampion.Prowess) * defendingChampion.GetRecentReputation())
				Game.Player.Glory += glory
				fmt.Printf("The Church earned %d glory for sponsoring %s. Your sponsorships have earned the church %d glory in total.\n", glory, attackingChampion.GetTitle(), Game.Player.Glory)
			}

			attackerAdvantage = 1
			KillKnight(defendingChampion)
		} else {
			fmt.Printf(
				"%s slayed %s on the battlefield after an intense duel[%d/%d vs %d/%d], giving %s a tactical edge!\n",
				defendingChampion.GetTitle(), attackingChampion.GetTitle(),
				defenderHits, defendingChampion.Prowess,
				attackerHits, attackingChampion.Prowess,
				defendingHouse.GetTitle(),
			)

			if defendingChampion.Sponsor != nil {
				glory := int(5 * float64(attackingChampion.Prowess) * attackingChampion.GetRecentReputation())
				Game.Player.Glory += glory
				fmt.Printf("The Church earned %d glory for sponsoring %s. You're sponsorships have earned the church %d glory in total.\n", glory, defendingChampion.GetTitle(), Game.Player.Glory)
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
		"%s[%d/%d hits] defeated %s[%d/%d hits]!\n",
		winner.GetTitle(), winnerHits, winner.Might, loser.GetTitle(), loserHits, loser.Might,
	)

	// Award more glory to underdogs and less to bullies.
	glory := (MaxMight + 1) + (loser.Might - winner.Might)
	for _, knight := range winner.Knights {
		knight.BattleResults = append(knight.BattleResults, Victory)
		if knight.Sponsor != nil {
			Game.Player.Glory += glory
			fmt.Printf("The Church earned %d glory for sponsoring %s. You're sponsorships have earned the church %d glory in total.\n", glory, knight.GetTitle(), Game.Player.Glory)
		}
	}

	loserKnightsCopy := make([]*Knight, len(loser.Knights))
	copy(loserKnightsCopy, loser.Knights)
	for _, knight := range loserKnightsCopy {
		knight.BattleResults = append(knight.BattleResults, Defeat)

		defeatSeverity := (winnerHits - loserHits) / 2
		survivalHits := RollHits(knight.Prowess)
		if survivalHits < defeatSeverity {
			fmt.Printf(
				"%s was overwhelmed by the enemy forces and killed[%d/%d vs %d]\n",
				knight.GetTitle(), survivalHits, knight.Prowess, defeatSeverity,
			)
			KillKnight(knight)
		}
	}

	return attackerHits - defenderHits
}

func DisplayState() {
	for _, house := range Game.Houses {
		fmt.Printf("Introducing the knights of %s[might: %d]! Their banner is %s.\n", house.GetTitle(), house.Might, house.Banner)
		for _, knight := range house.Knights {
			fmt.Printf(
				"%s! [prowess: %d, bravery: %d, cost: %d]\n",
				knight.GetTitle(), knight.Prowess, knight.Bravery, knight.GetCost(),
			)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n===================\n\n")
}

func FindKnightByName(knightName string) *Knight {
	for _, knight := range Game.Knights {
		if knight.Name == knightName {
			return knight
		}
	}
	return nil
}

// TODO: Can this be done with generics or interfaces?
func FindHouseByName(houseName string) *House {
	for _, house := range Game.Houses {
		if house.Name == houseName {
			return house
		}
	}
	return nil
}

var Game *GameState

