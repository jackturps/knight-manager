package game

import (
	"fmt"
	"knightmanager/names"
	"strings"
)

var MaxMight = 5
var MaxWealth = 5

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
	Wealth int

	Knights []*Knight
	DiplomaticRelations map[*House]*DiplomaticRelation
}

func (house *House) GetTitle() string {
	return fmt.Sprintf("House %s", house.Name)
}

// NumWars returns the number of wars a house is currently participating in.
func (house *House) NumWars() int {
	numWars := 0
	for _, war := range Game.Wars {
		if war.Attackers.Leader == house {
			numWars++
		} else if war.Defenders.Leader == house {
			numWars++
		} else if Exists(war.Attackers.Allies, house) {
			numWars++
		} else if Exists(war.Defenders.Allies, house) {
			numWars++
		}
	}
	return numWars
}

func (house *House) GetAdjustedMight() int {
	// Reduce might for each extra war the house is in.
	return house.Might - Max[int](0, house.NumWars() - 1)
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

	Cycle int
}

func AssignKnightToHouse(knight *Knight, house *House) {
	house.Knights = append(house.Knights, knight)
	knight.House = house
}

func SponsorKnight(bishop *GloryBishop, knight *Knight) {
	bishop.SponsoredKnights = append(bishop.SponsoredKnights, knight)
	knight.Sponsor = bishop
}

func DestroyHouse(destroyedHouse *House) {
	for _, house := range Game.Houses {
		delete(house.DiplomaticRelations, destroyedHouse)
	}
	for _, war := range CopySlice(Game.Wars) {
		// End the war if the destroyedHouse is a primary fighter.
		if war.Attackers.Leader == destroyedHouse {
			fmt.Printf("%s could no longer fight in the war against %s. The war is over.\n", destroyedHouse.GetTitle(), war.Defenders.Leader.GetTitle())
			Game.Wars = RemoveItem(Game.Wars, war)
		}
		if war.Defenders.Leader == destroyedHouse {
			fmt.Printf("%s could no longer fight in the war against %s. The war is over.\n", destroyedHouse.GetTitle(), war.Attackers.Leader.GetTitle())
			Game.Wars = RemoveItem(Game.Wars, war)
		}

		// If the destroyedHouse was just an ally, remove them from the allies.
		war.Attackers.Allies = RemoveItem(war.Attackers.Allies, destroyedHouse)
		war.Defenders.Allies = RemoveItem(war.Defenders.Allies, destroyedHouse)
	}
	for _, knight := range destroyedHouse.Knights {
		Game.Knights = RemoveItem(Game.Knights, knight)
	}
	Game.Houses = RemoveItem(Game.Houses, destroyedHouse)
}

func GenerateHouse() *House {
	house := &House{
		Name:   Game.FemaleNameGenerator.GenerateName(),
		Banner: GenerateBanner(),
		Might:  RandomRange(1, MaxMight + 1),
		// TODO: Make wealth related to might of house in some way?
		Wealth: RandomRange(1, MaxWealth + 1),
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
	attackingKnight := ChooseHouseChampion(attackingHouse)
	defendingKnight := ChooseHouseChampion(defendingHouse)

	if attackingKnight == nil && defendingKnight == nil {
		fmt.Printf("Neither house could field a champion!\n")
	} else if attackingKnight == nil {
		defenderAdvantage = 1
		fmt.Printf("%s could not field a champion, giving %s a tactical edge!\n", attackingHouse.GetTitle(), defendingHouse.GetTitle())
	} else if defendingKnight == nil {
		attackerAdvantage = 1
		fmt.Printf("%s could not field a champion, giving %s a tactical edge!\n", defendingHouse.GetTitle(), attackingHouse.GetTitle())
	} else {
		attackerHits := RollHits(attackingKnight.Prowess + attackingKnight.Blessings)
		defenderHits := RollHits(defendingKnight.Prowess + defendingKnight.Blessings)
		// TODO: Maybe only kill if the margin is big enough.
		// TODO: Award glory to victorious champions.
		if attackerHits == defenderHits {
			fmt.Printf(
				"%s met %s on the battlefield, their duel raged until it met a stalemate[%d/%dd+%dd vs %d/%dd+%dd]!\n",
				attackingKnight.GetTitle(), defendingKnight.GetTitle(),
				attackerHits, attackingKnight.Prowess, attackingKnight.Blessings,
				defenderHits, defendingKnight.Prowess, defendingKnight.Blessings,
			)
		} else if attackerHits > defenderHits {
			// TODO: Reduce duplication between here and below.
			fmt.Printf(
				"%s slayed %s on the battlefield after an intense duel[%d/%dd+%dd vs %d/%dd+%dd], giving %s a tactical edge!\n",
				attackingKnight.GetTitle(), defendingKnight.GetTitle(),
				attackerHits, attackingKnight.Prowess, attackingKnight.Blessings,
				defenderHits, defendingKnight.Prowess, defendingKnight.Blessings,
				attackingHouse.GetTitle(),
			)

			if attackingKnight.Sponsor != nil {
				glory := int(5 * float64(defendingKnight.Prowess) * defendingKnight.GetRecentReputation())
				Game.Player.Glory += glory
				fmt.Printf("The Church earned %d glory for sponsoring %s.\n", glory, attackingKnight.GetTitle())
			}

			attackerAdvantage = 1
			KillKnight(defendingKnight)
		} else {
			fmt.Printf(
				"%s slayed %s on the battlefield after an intense duel[%d/%dd+%dd vs %d/%dd+%dd], giving %s a tactical edge!\n",
				defendingKnight.GetTitle(), attackingKnight.GetTitle(),
				defenderHits, defendingKnight.Prowess, defendingKnight.Blessings,
				attackerHits, attackingKnight.Prowess, attackingKnight.Blessings,
				defendingHouse.GetTitle(),
			)

			if defendingKnight.Sponsor != nil {
				glory := int(5 * float64(attackingKnight.Prowess) * attackingKnight.GetRecentReputation())
				Game.Player.Glory += glory
				fmt.Printf("The Church earned %d glory for sponsoring %s.\n", glory, defendingKnight.GetTitle())
			}

			defenderAdvantage = 1
			KillKnight(attackingKnight)
		}
	}

	// Blessings only last one fight.
	if attackingKnight != nil {
		attackingKnight.Blessings = 0
	}
	if defendingKnight != nil {
		defendingKnight.Blessings = 0
	}

	attackerHits := RollHits(attackingHouse.GetAdjustedMight() + attackerAdvantage)
	defenderHits := RollHits(defendingHouse.GetAdjustedMight() + defenderAdvantage)

	var winner, loser *House
	var winnerHits, loserHits int

	if attackerHits > defenderHits {
		winner, winnerHits, loser, loserHits = attackingHouse, attackerHits, defendingHouse, defenderHits
	} else {
		winner, winnerHits, loser, loserHits = defendingHouse, defenderHits, attackingHouse, attackerHits
	}
	// TODO: Print advantages?
	fmt.Printf(
		"%s[%d/%dd hits] defeated %s[%d/%dd hits]!\n",
		winner.GetTitle(), winnerHits, winner.GetAdjustedMight(), loser.GetTitle(), loserHits, loser.GetAdjustedMight(),
	)

	// TODO: Remove glory for winning battle? Too easy?
	// Award more glory to underdogs and less to bullies.
	glory := (MaxMight + 1) + (loser.Might - winner.Might)
	for _, knight := range winner.Knights {
		knight.BattleResults = append(knight.BattleResults, Victory)
		if knight.Sponsor != nil {
			Game.Player.Glory += glory
			fmt.Printf("The Church earned %d glory for sponsoring %s.\n", glory, knight.GetTitle())
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
				"%s was overwhelmed by the enemy forces and killed[%d/%dd vs %d]\n",
				knight.GetTitle(), survivalHits, knight.Prowess, defeatSeverity,
			)
			KillKnight(knight)
		}
	}

	return attackerHits - defenderHits
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

