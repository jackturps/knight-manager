package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
			fmt.Printf("\n====================\n\n")
			break
		} else if command[0] == "help" {
			fmt.Printf(
				"sponsor <first-name>: pay a knight's cost in coin to sponsor them, gaining glory from their victories\n" +
					"research <first-name>: discover information about a knight\n" +
					"done: finalise your sponsorships for this season\n",
			)
		} else if command[0] == "research" {
			if len(command) < 2 {
				fmt.Printf("Specify a knight(research <first-name>)\n")
				continue
			}

			knightName := command[1]
			knight := FindKnightByName(knightName)
			if knight == nil {
				fmt.Printf("Could not find knight '%s'\n", knightName)
				continue
			}

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
		}
	}
}

