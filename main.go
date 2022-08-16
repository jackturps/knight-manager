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
		"Spend it wisely.\n\n",
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

	for {
		game.DoPlayerTurn()

		for idx := 0; idx < 3; idx++ {
			game.DoWorldEvent()
		}
		fmt.Printf("\n")

		// TODO: Only roll for start war after an insighting incident so every war has a cause?
		game.StartWars()

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

		// TODO: Roll house's wealth to see who gets knights?
		// Round robin which houses get new knights.
		for idx := 0; idx < numNewKnightsPerSeason; idx++ {
			house := game.Game.Houses[knightedHouseIdx]
			game.GenerateKnight(house)
			knightedHouseIdx = (knightedHouseIdx + 1) % len(game.Game.Houses)
		}
	}
}
