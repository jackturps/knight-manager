package main

import (
	"fmt"
	"knightmanager/game"
	"knightmanager/names"
	"math/rand"
	"sort"
	"time"
)

/**
- Starting a war against a tiny house increases everyones tension? No one likes a tyrant.
 */

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Printf(
		"The Church brings glory to the many gods by using it's resources " +
		"to make the gods' values more prevalent in the world. Glory is brought " +
		"to each god differently. You are a bishop of the God of War, You bring " +
		"glory to Them by sponsoring knights of the great houses, and having those knights " +
		"see success in battle. Of course, good knights do not come cheap - but if their " +
		"house falls on hard times bargains may present themselves. Be warned, a dead " +
		"knight knows no glory.\n\nWhen sponsored knights die their house will pay the church " +
		"a customary funeral tithe. Wealthier houses are more generous in their offerings to the gods. " +
		"Spend it wisely.\n\nIn addition to your public duties, the prophets of the church also whisper " +
		"a secret agenda in your ear. In the flames they have seen the faces of several knights rising " +
		"to become mesiahs. The prophecy states that these would-be mesiahs are challenged only by a " +
		"group of heretical knights. Ensure that at least one of the %s outlives all of the %s\n\n",
		game.ColouredText(game.BlueBackgroundCode, "mesiahs"),
		game.ColouredText(game.RedBackgroundCode, "heretics"),
	)

	game.Game = &game.GameState{}
	game.Game.Wars = make([]*game.War, 0)
	game.Game.FemaleNameGenerator = names.NewSelectorNameGenerator("female_input_names.txt")
	game.Game.MaleNameGenerator = names.NewSelectorNameGenerator("male_input_names.txt")
	game.GenerateWorld()

	game.Game.Player = &game.GloryBishop{
		Coin: 30,
		Glory: 0,
	}

	knightedHouseIdx := game.RandomRange(0, len(game.Game.Houses))
	numNewKnightsPerSeason := 2


	sortedKnights := game.CopySlice(game.Game.Knights)
	sort.Slice(sortedKnights, func(x, y int) bool {
		knightScore1 := sortedKnights[x].House.Might + sortedKnights[x].Prowess
		knightScore2 := sortedKnights[y].House.Might + sortedKnights[y].Prowess
		return knightScore1 < knightScore2
	})

	sortedKnights[0].ChurchObjective = game.Protect
	sortedKnights[1].ChurchObjective = game.Protect
	sortedKnights[2].ChurchObjective = game.Protect

	sortedKnights[len(sortedKnights) - 1].ChurchObjective = game.Kill
	sortedKnights[len(sortedKnights) - 2].ChurchObjective = game.Kill
	sortedKnights[len(sortedKnights) - 3].ChurchObjective = game.Kill

	for {
		game.Game.Cycle++
		game.DoPlayerTurn()

		for idx := 0; idx < 3; idx++ {
			game.DoWorldEvent()
		}
		fmt.Printf("\n")

		for _, war := range game.CopySlice(game.Game.Wars) {
			// If a house is destroyed in another war this turn any of their other wars.
			// will end. We should only run battles for wars that are still going.
			if game.Exists(game.Game.Wars, war) {
				war.DoNextBattles()
			}
			if war.IsOver() {
				war.EndWar()
			}
		}
		if len(game.Game.Wars) > 0 {
			fmt.Printf("\n")
		}

		game.CheckForNicknames()

		// TODO: Only roll for start war after an insighting incident so every war has a cause?
		game.StartWars()

		// TODO: Roll house's wealth to see who gets knights?
		// Round robin which houses get new knights.
		for idx := 0; idx < numNewKnightsPerSeason; idx++ {
			house := game.Game.Houses[knightedHouseIdx]
			game.GenerateKnight(house)
			knightedHouseIdx = (knightedHouseIdx + 1) % len(game.Game.Houses)
		}

		numProtectedKnights := 0
		numKillKnights := 0
		for _, knight := range game.Game.Knights {
			if knight.ChurchObjective == game.Protect {
				numProtectedKnights++
			} else if knight.ChurchObjective == game.Kill {
				numKillKnights++
			}
		}
		if numProtectedKnights == 0 {
			fmt.Printf(
				"All possible Mesiahs have died or been stripped of their nobility. You are cast " +
				"out from the church.\n",
			)
			break
		}
		if numKillKnights == 0 {
			fmt.Printf(
				"You protected the mesiah(s)!!! %d STAR VICTORY!!!\n", numProtectedKnights,
			)
			break
		}
	}
}
