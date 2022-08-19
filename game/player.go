package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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

	fmt.Printf("%s fights with a %s\n", knight.GetTitle(), knight.Weapon.Type)

	if knight.Spouse == nil {
		fmt.Printf("%s is unmarried.\n", knight.GetTitle())
	} else {
		fmt.Printf("%s is married to %s.\n", knight.GetTitle(), knight.Spouse.GetTitle())
	}

	fmt.Printf(
		"%s has fought in %d battles, their results are: %s\n",
		knight.GetTitle(), len(knight.BattleResults), battleResultString,
	)

	slayedKnightsText := ""
	isFirstSlayedKnight := true
	for _, slayedKnight := range knight.SlayedKnights {
		if !isFirstSlayedKnight {
			slayedKnightsText += ", "
		}
		isFirstSlayedKnight = false
		slayedKnightsText += slayedKnight.GetTitle()
	}

	fmt.Printf("%s has killed %d knight(s) in battles: %s\n", knight.GetTitle(), len(knight.SlayedKnights), slayedKnightsText)
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
		fmt.Printf("Introducing the knights of %s[might: %d, wealth: %d]! Their banner is %s.\n", house.GetTitle(), house.Might, house.Wealth, house.Banner.GetDescription())
		for _, knight := range house.Knights {
			fmt.Printf(
				"%s! [prowess: %d, bravery: %d, cost: %d]\n",
				knight.GetTitle(), knight.Prowess, knight.Bravery, knight.GetCost(),
			)
		}
		fmt.Printf("\n")
	}
}

func DisplayWars() {
	for _, war := range Game.Wars {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(
			w, "Turn\tAttackers[morale: %d]\tDefenders[morale: %d]\n",
			war.Attackers.Morale, war.Defenders.Morale,
		)

		attacker := war.Attackers.Leader
		defender := war.Defenders.Leader

		turnIcon := ""
		if war.attackingHouseIdx == 0 {
			turnIcon = "*"
		}

		fmt.Fprintf(
			w,"%s\t%s[might: %d]\t%s[might: %d]\n",
			turnIcon,
			attacker.GetTitle(), attacker.Might,
			defender.GetTitle(), defender.Might,
		)

		for allyIdx := 0; allyIdx < Max[int](len(war.Attackers.Allies), len(war.Defenders.Allies)); allyIdx++ {
			if war.attackingHouseIdx - 1 == allyIdx {
				fmt.Fprintf(w, "*\t")
			} else {
				fmt.Fprintf(w, "\t")
			}

			if allyIdx < len(war.Attackers.Allies) {
				ally := war.Attackers.Allies[allyIdx]
				fmt.Fprintf(
					w, "%s[might: %d]\t",
					ally.GetTitle(), ally.Might,
				)
			} else {
				fmt.Fprintf(w, "\t")
			}

			if allyIdx < len(war.Defenders.Allies) {
				ally := war.Defenders.Allies[allyIdx]
				fmt.Fprintf(
					w, "%s[might: %d]\n",
					ally.GetTitle(), ally.Might,
				)
			} else {
				fmt.Fprintf(w, "\n")
			}
		}
		w.Flush()
		fmt.Printf("\n")
	}
}

func DisplayDiplomacy() {
	tensionSeverity := []string{
		GreenTextCode,
		YellowTextCode,
		RedTextCode,
	}

	// NOTE: All fields need to be wrapped in some colour codes so they all get fucked
	// up in the same way and the table works.
	// TODO: Use a table that can handle coloured text.
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	fmt.Fprintf(w, "%s\t", ColouredText(DefaultColourCode, ""))
	for _, house := range Game.Houses {
		fmt.Fprintf(w, "%s\t", ColouredText(DefaultColourCode, house.Name))
	}
	fmt.Fprint(w, "\n")

	for _, sourceHouse := range Game.Houses {
		fmt.Fprintf(w, "%s\t", ColouredText(DefaultColourCode, sourceHouse.Name))
		for _, targetHouse := range Game.Houses {
			if sourceHouse == targetHouse {
				fmt.Fprintf(w, "%s\t", ColouredText(DefaultColourCode ,"X"))
			} else {
				// TODO: Colour numbers on severity?
				tension := sourceHouse.DiplomaticRelations[targetHouse].Tension
				tensionColour := tensionSeverity[Min[int](2, tension / 3)]
				tensionText := ColouredText(tensionColour, strconv.Itoa(tension))
				fmt.Fprintf(w, "%s\t", tensionText)
			}
		}
		fmt.Fprint(w, "\n")
	}

	w.Flush()
}

func DoPlayerTurn() {
	fmt.Printf("Year %d - You have %d coin and %d glory.\n", Game.Cycle, Game.Player.Coin, Game.Player.Glory)

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
				"sponsor <knight-name>: pay a knight's cost in coin to sponsor them, gaining glory from their victories and coin when they die.\n" +
					"marry <knight-name> <knight-name>: marry two knights, moving a knight from the weaker house into the stronger house. This reduces tension between the houses.\n" +
					"research <knight-name|house-name>: discover information about a knight or house.\n" +
					"houses: display information about all houses.\n" +
					"wars: display information about all in progress wars.\n" +
					"tensions: show the tensions between each of the houses.\n" +
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
			DisplayWars()
		} else if command[0] == "marry" {
			if len(command) < 3 {
				fmt.Printf("Specify knights to marry(marry <first-name> <first-name>)\n")
				continue
			}

			knight1 := FindKnightByName(command[1])
			if knight1 == nil {
				fmt.Printf("Could not find knight '%s'\n", command[1])
				continue
			}
			knight2 := FindKnightByName(command[2])
			if knight2 == nil {
				fmt.Printf("Could not find knight '%s'\n", command[2])
				continue
			}

			MarryKnights(knight1, knight2)
		} else if command[0] == "tensions" {
			DisplayDiplomacy()
		} else if command[0] == "bless" {
			if len(command) < 2 {
				fmt.Printf("Specify knight to bless(bless <first-name>)\n")
				continue
			}
			knight := FindKnightByName(command[1])
			if knight == nil {
				fmt.Printf("Could not find knight '%s'\n", command[1])
				continue
			}

			gloryCost := (knight.Blessings + 1) * 10
			if Game.Player.Glory < gloryCost {
				fmt.Printf(
					"It costs %d glory to bless %s, you have %d.\n",
					gloryCost, knight.GetTitle(), Game.Player.Glory,
				)
				continue
			}

			Game.Player.Glory -= gloryCost
			knight.Blessings++
			fmt.Printf("%s will now have +%dd in duels.\n", knight.GetTitle(), knight.Blessings)
		}
	}
}

