package game

import "fmt"

func MarryKnights(knight1 *Knight, knight2 *Knight) {
	// TODO: Spend glory to marry knights.
	if HousesAreAtWar(knight1.House, knight2.House) {
		fmt.Printf(
			"%s and %s are at war, they refuse to marry %s and %s.\n",
			knight1.House.GetTitle(), knight2.House.GetTitle(),
			knight1.GetTitle(), knight2.GetTitle(),
		)
		return
	}
	if knight1.House == knight2.House {
		fmt.Printf(
			"%s and %s are from the same house, they cannot be wed.\n",
			knight1.GetTitle(), knight2.GetTitle(),
		)
		return
	}
	if knight1.Spouse != nil {
		fmt.Printf(
			"%s is already married to %s, they cannot be wed again.\n",
			knight1.GetTitle(), knight1.Spouse.GetTitle(),
		)
		return
	}
	if knight2.Spouse != nil {
		fmt.Printf(
			"%s is already married to %s, they cannot be wed again.\n",
			knight2.GetTitle(), knight2.Spouse.GetTitle(),
		)
		return
	}

	requiredGlory := 50
	if Game.Player.Glory < requiredGlory {
		fmt.Printf("Arranging a marriage costs %d glory, you only have %d.\n", requiredGlory, Game.Player.Glory)
		return
	}
	Game.Player.Glory -= requiredGlory

	var movingKnight, stayingKnight *Knight
	if knight1.House.Might > knight2.House.Might {
		stayingKnight, movingKnight = knight1, knight2
	} else if knight2.House.Might > knight1.House.Might {
		stayingKnight, movingKnight = knight2, knight1
	} else {
		// If might matches the user can decide based on the order they give.
		movingKnight, stayingKnight = knight1, knight2
	}

	fmt.Printf(
		"Marrying %s to %s. %s will become a member of %s.\n",
		movingKnight.GetTitle(), stayingKnight.GetTitle(), movingKnight.GetTitle(), stayingKnight.House.GetTitle(),
	)

	tensionReducedAmount := 5
	fmt.Printf(
		"Tensions between %s and %s are reduced by %d.\n",
		movingKnight.House.GetTitle(), stayingKnight.House.GetTitle(), tensionReducedAmount,
	)
	// NOTE: Modify tensions before moving knights so we don't lose the reference to the house of
	// the moving knight.
	relation1 := movingKnight.House.DiplomaticRelations[stayingKnight.House]
	relation2 := stayingKnight.House.DiplomaticRelations[movingKnight.House]
	relation1.Tension = Max(relation1.Tension - tensionReducedAmount, 0)
	relation2.Tension = Max(relation2.Tension - tensionReducedAmount, 0)

	movingKnight.House.Knights = RemoveItem(movingKnight.House.Knights, movingKnight)
	stayingKnight.House.Knights = append(stayingKnight.House.Knights, movingKnight)
	movingKnight.House = stayingKnight.House

	movingKnight.Spouse = stayingKnight
	stayingKnight.Spouse = movingKnight
}