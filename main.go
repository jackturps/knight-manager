package main

import (
	"fmt"
	"knightmanager/game"
	"knightmanager/names"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Printf(
		"The Church brings glory to the many gods by using it's resources " +
		"to make the gods' values more prevalent in the world. Glory is brought " +
		"to each god differently. You are a bishop of the God of War, You bring " +
		"glory to Them by sponsoring knights of the great houses, and having those knights " +
		"see success in battle. Of course, good knights do not come cheap - but if their " +
		"house falls on hard times bargains may present themselves. Be warned, a dead " +
		"knight knows no glory.\n\nThe church provides you 5 coin per season to sponsor " +
		"knights. Spend it wisely.\n\n",
	)

	game.Game = &game.GameState{}
	game.Game.FemaleNameGenerator = names.NewSelectorNameGenerator("female_input_names.txt")
	game.Game.MaleNameGenerator = names.NewSelectorNameGenerator("male_input_names.txt")
	game.GenerateWorld()

	game.Game.Player = &game.GloryBishop{
		Coin: 15,
		Glory: 0,
	}

	knightedHouseIdx := game.RandomRange(0, len(game.Game.Houses))
	numNewKnightsPerSeason := 3
	numBattlesPerSeason := 3

	for {
		game.DisplayState()
		game.DoPlayerTurn()
		for idx := 0; idx < numBattlesPerSeason; idx++ {
			game.RunBattle(game.Game.Houses)
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Printf("\n=================\n\n")


		// Round robin which houses get new knights.
		for idx := 0; idx < numNewKnightsPerSeason; idx++ {
			house := game.Game.Houses[knightedHouseIdx]
			game.GenerateKnight(house)
			knightedHouseIdx = (knightedHouseIdx + 1) % len(game.Game.Houses)
		}

		game.Game.Player.Coin += 5
	}
}
