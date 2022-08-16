package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Research(entityName string) {
	if knight := FindKnightByName(entityName); knight != nil {
		ResearchKnight(knight)
	} else if house := FindHouseByName(entityName); house != nil {
		ResearchHouse(house)
	} else {
		fmt.Printf("Could not find knight or house with name '%s'\n", entityName)
	}
}

func ResearchKnight(knight *Knight) {
	battleResultString := ""
	for _, battleResult := range knight.BattleResults {
		if battleResult == Victory {
			battleResultString += "V "
		} else {
			battleResultString += "D "
		}
	}
	fmt.Printf(
		"%s has fought in %d battles, their results are: %s\n",
		knight.GetTitle(), len(knight.BattleResults), battleResultString,
	)
}

func ResearchHouse(house *House) {
	for targetHouse, relation := range house.DiplomaticRelations {
		fmt.Printf(
			"%s's tensions with %s are at %d\n",
			house.GetTitle(), targetHouse.GetTitle(), relation.Tension,
		)
	}
}

func DisplayHouses() {
	for _, house := range Game.Houses {
		fmt.Printf("Introducing the knights of %s[might: %d, wealth: %d]! Their banner is %s.\n", house.GetTitle(), house.Might, house.Wealth, house.Banner)
		for _, knight := range house.Knights {
			fmt.Printf(
				"%s! [prowess: %d, bravery: %d, cost: %d]\n",
				knight.GetTitle(), knight.Prowess, knight.Bravery, knight.GetCost(),
			)
		}
		fmt.Printf("\n")
	}
}

func DoPlayerTurn() {
	fmt.Printf("You have %d coin\n", Game.Player.Coin)

	// Player interaction loop.
	// TODO: For the love of god clean this up.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		command := strings.Split(input, " ")
		if command[0] == "done" {
			break
		} else if command[0] == "help" {
			fmt.Printf(
				"sponsor <knight-name>: pay a knight's cost in coin to sponsor them, gaining glory from their victories\n" +
					"research <knight-name|house-name>: discover information about a knight or house\n" +
					"houses: display information about all houses\n" +
					"wars: display information about all in progress wars\n" +
					"done: finalise your sponsorships for this season\n",
			)
		} else if command[0] == "research" {
			if len(command) < 2 {
				fmt.Printf("Specify a knight or house(research <first-name>)\n")
				continue
			}

			entityName := command[1]
			Research(entityName)
		} else if command[0] == "sponsor" {
			if len(command) < 2 {
				fmt.Printf("Specify a knight(sponsor <first-name>)\n")
				continue
			}

			knightName := command[1]
			foundKnight := FindKnightByName(knightName)
			if foundKnight == nil {
				fmt.Printf("Could not find knight '%s'\n", knightName)
				continue
			}
			if foundKnight.Sponsor != nil {
				fmt.Printf("%s is already sponsored\n", foundKnight.GetTitle())
				continue
			}

			cost := foundKnight.GetCost()
			if cost > Game.Player.Coin {
				fmt.Printf("The church coffers run low, %s costs %d coin but you only have %d.\n", foundKnight.GetTitle(), cost, Game.Player.Coin)
				continue
			}

			Game.Player.Coin -= foundKnight.GetCost()
			SponsorKnight(Game.Player, foundKnight)
			fmt.Printf(
				"You have sponsored %s, %d coin remaining\n",
				foundKnight.GetTitle(), Game.Player.Coin,
			)
		} else if command[0] == "houses" {
			DisplayHouses()
		} else if command[0] == "wars" {
			for _, war := range Game.Wars {
				attacker := war.Attackers.Leader
				defender := war.Defenders.Leader

				fmt.Printf(
					"%s[might: %d, morale: %d] is at war with %s[might: %d, morale: %d]\n",
					attacker.GetTitle(), attacker.Might, war.Attackers.Morale,
					defender.GetTitle(), defender.Might, war.Defenders.Morale,
				)
				for _, ally := range war.Attackers.Allies {
					fmt.Printf("%s[might: %d] is allied with %s\n", ally.GetTitle(), ally.Might, attacker.GetTitle())
				}
				for _, ally := range war.Defenders.Allies {
					fmt.Printf("%s[might: %d] is allied with %s\n", ally.GetTitle(), ally.Might, defender.GetTitle())
				}

				fmt.Printf("\n")
			}
		}
	}
}

