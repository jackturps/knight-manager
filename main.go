package main

import (
	"fmt"
	"knightmanager/game"
	"knightmanager/names"
	"math/rand"
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
		"a customary funeral tithe. Wealthier houses are more generous in their offering to the gods. " +
		"Spend it wisely.\n\nIn addition to your public duties, the church also has a hidden agenda. " +
		"The prophets speak of a new Mesiah amongst the ranks of the great houses. They have seen " +
		"many possible faces in the flames. Ensure at least one of the chosen knights make it to year 50.\n\n",
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


	var bestKnight, worstKnight, randomKnight *game.Knight
	for _, knight := range game.Game.Knights {
		if bestKnight == nil ||
			(knight.Prowess > bestKnight.Prowess) ||
			(knight.Prowess == bestKnight.Prowess && knight.House.Might > bestKnight.House.Might) {
			bestKnight = knight
		}

		if worstKnight == nil ||
			(knight.Prowess < worstKnight.Prowess) ||
			(knight.Prowess == worstKnight.Prowess && knight.House.Might < worstKnight.House.Might) {
			worstKnight = knight
		}
	}
	possibleKnights := game.RemoveItem(game.RemoveItem(game.Game.Knights, bestKnight), worstKnight)
	randomKnight = game.RandomSelect(possibleKnights)

	bestKnight.IsChosen = true
	worstKnight.IsChosen = true
	randomKnight.IsChosen = true

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

		// TODO: Only roll for start war after an insighting incident so every war has a cause?
		game.StartWars()

		// TODO: Roll house's wealth to see who gets knights?
		// Round robin which houses get new knights.
		for idx := 0; idx < numNewKnightsPerSeason; idx++ {
			house := game.Game.Houses[knightedHouseIdx]
			game.GenerateKnight(house)
			knightedHouseIdx = (knightedHouseIdx + 1) % len(game.Game.Houses)
		}

		remainingChosen := 0
		for _, knight := range game.Game.Knights {
			if knight.IsChosen {
				remainingChosen++
			}
		}
		if remainingChosen == 0 {
			fmt.Printf(
				"All possible Mesiahs have died or been stripped of their nobility. You are cast " +
				"out from the church.\n",
			)
			break
		}

		if game.Game.Cycle == 50 {
			fmt.Printf(
				"You protected the mesiah(s)!!! %d STAR VICTORY!!!\n", remainingChosen,
			)
		}
	}
}
