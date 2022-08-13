package game

import (
	"fmt"
	"math"
)

type War struct {
	Attackers *Alliance
	Defenders *Alliance

	attackingHouseIdx int
}

type Alliance struct {
	Leader *House
	Allies []*House
	Morale int
}

func (alliance *Alliance) GetTotalMight() int {
	totalMight := alliance.Leader.Might
	for _, ally := range alliance.Allies {
		totalMight += ally.Might
	}
	return totalMight
}

func HouseWillJoinAlliance(allyHouse *House, alliance *Alliance, enemy *Alliance, otherEnemies []*House) bool {
	isEnemy := allyHouse == enemy.Leader
	isLeaderHouse := allyHouse == alliance.Leader
	isAlliedWithEnemy := Exists(enemy.Allies, allyHouse) || Exists(otherEnemies, allyHouse)
	// TODO: Check across all wars not just this one.
	if isEnemy || isLeaderHouse || isAlliedWithEnemy {
		return false
	}

	tensionWithTarget := allyHouse.DiplomaticRelations[enemy.Leader].Tension
	tensionWithLeader := allyHouse.DiplomaticRelations[alliance.Leader].Tension
	relativeTension := tensionWithTarget - tensionWithLeader

	joinAlliancePool := int(math.Max(0, float64(relativeTension + allyHouse.Might)))
	joinAllianceHits := RollHits(joinAlliancePool)
	willJoin := joinAllianceHits >= enemy.GetTotalMight()

	if willJoin {
		fmt.Printf(
			"%s allied with %s in the war against %s! [%d/%d vs %d]\n",
			allyHouse.GetTitle(), alliance.Leader.GetTitle(), enemy.Leader.GetTitle(),
			joinAllianceHits, joinAlliancePool, enemy.GetTotalMight(),
		)
	}
	return willJoin
}

func CreateWar(attackerHouse *House, defenderHouse *House) *War {
	// TODO: Set morale based on some stat. Maybe median bravery?
	war := &War{
		Attackers: &Alliance{
			Leader: attackerHouse,
			Allies: make([]*House, 0),
			Morale: 6,
		},
		Defenders: &Alliance{
			Leader: defenderHouse,
			Allies: make([]*House, 0),
			Morale: 6,
		},
		attackingHouseIdx: 0,
	}

	fmt.Printf("%s declared war against %s!\n", attackerHouse.GetTitle(), defenderHouse.GetTitle())

	randomizedHouses := RandomizeOrder(Game.Houses)

	// TODO: Get an ally for each house before increasing the house might to prevent favoring the attacker.
	attackerAllyIdx := 0
	defenderAllyIdx := 0
	for attackerAllyIdx < len(Game.Houses) || defenderAllyIdx < len(Game.Houses) {
		var attackerAlly *House = nil
		var defenderAlly *House = nil

		for ; attackerAllyIdx < len(Game.Houses); attackerAllyIdx++ {
			allyHouse := randomizedHouses[attackerAllyIdx]
			if HouseWillJoinAlliance(allyHouse, war.Attackers, war.Defenders, []*House{}) {
				attackerAlly = allyHouse
				attackerAllyIdx++
				break
			}
		}

		for ; defenderAllyIdx < len(Game.Houses); defenderAllyIdx++ {
			allyHouse := randomizedHouses[defenderAllyIdx]
			if HouseWillJoinAlliance(allyHouse, war.Defenders, war.Attackers, []*House{attackerAlly}) {
				defenderAlly = allyHouse
				defenderAllyIdx++
				break
			}
		}

		// Only append allies after both have selected to prevent biases to the first house to choose.
		if attackerAlly != nil {
			war.Attackers.Allies = append(war.Attackers.Allies, attackerAlly)
		}
		if defenderAlly != nil {
			war.Defenders.Allies = append(war.Defenders.Allies, defenderAlly)
		}
	}
	fmt.Printf("\n")

	return war
}

func StartWars() {
	for _, house := range RandomizeOrder(Game.Houses) {
		for targetHouse, relationship := range house.DiplomaticRelations {
			tensionHits := RollHits(relationship.Tension)

			// TODO: The ob should probably have another factor/be higher here, otherwise weak houses get trampled.
			// TODO: Opponent might should be in relation to your might. Subtract or divide?
			// TODO: Start multiple wars if we need to. Make wars a bit less common.
			if tensionHits >= targetHouse.Might {
				war := CreateWar(house, targetHouse)
				Game.Wars = append(Game.Wars, war)
				return
			}
		}
	}
	fmt.Printf("\n")
}

func (war *War) DoNextBattles() {
	allAttackers := make([]*House, 0, len(war.Attackers.Allies) + 1)
	allAttackers = append(allAttackers, war.Attackers.Leader)
	allAttackers = append(allAttackers, war.Attackers.Allies...)

	allDefenders := make([]*House, 0, len(war.Defenders.Allies) + 1)
	allDefenders = append(allDefenders, war.Defenders.Leader)
	allDefenders = append(allDefenders, war.Defenders.Allies...)

	maxHouseIdx := Max[int](len(allAttackers), len(allDefenders))

	// Every house on each side attacks a randome opponent. More allies means more attacks.
	if war.attackingHouseIdx < len(allAttackers) {
		attacker := allAttackers[war.attackingHouseIdx]
		defender := RandomSelect(allDefenders)
		attackerMargin := RunBattle(attacker, defender)
		if attackerMargin > 0 {
			war.Defenders.Morale -= attackerMargin
			fmt.Printf(
				"The morale of %s's alliance is now at %d\n",
				war.Defenders.Leader.GetTitle(), war.Defenders.Morale,
			)
		}
		fmt.Printf("\n")
	}

	if war.attackingHouseIdx < len(allDefenders) {
		attacker := allDefenders[war.attackingHouseIdx]
		defender := RandomSelect(allAttackers)
		attackerMargin := RunBattle(attacker, defender)
		if attackerMargin > 0 {
			war.Attackers.Morale -= attackerMargin
			fmt.Printf(
				"The morale of %s's alliance is now at %d\n",
				war.Attackers.Leader.GetTitle(), war.Attackers.Morale,
			)
		}
		fmt.Printf("\n")
	}

	war.attackingHouseIdx = (war.attackingHouseIdx + 1) % maxHouseIdx
}

func (war *War) IsOver() bool {
	return war.Attackers.Morale <= 0 || war.Defenders.Morale <= 0
}

func (war *War) EndWar() {
	if !war.IsOver() {
		panic("tried to end war when it wasn't over.")
	}

	// TODO: Destroy houses that lose with might of 1.
	attackLeader := war.Attackers.Leader
	defenseLeader := war.Defenders.Leader

	if war.Attackers.Morale <= 0 && war.Defenders.Morale <= 0 {
		fmt.Printf(
			"The war between %s and %s ended in a truce after significant losses on both sides. " +
				"The might of both houses is reduced by 1.\n\n",
			attackLeader.GetTitle(), defenseLeader.GetTitle(),
		)
		attackLeader.Might = Max[int](attackLeader.Might - 1, 1)
		defenseLeader.Might = Max[int](defenseLeader.Might - 1, 1)
	} else if war.Attackers.Morale <= 0 {
		GiveWarRewards(war.Defenders, war.Attackers)
	} else if war.Defenders.Morale <= 0 {
		GiveWarRewards(war.Attackers, war.Defenders)
	}

	attackLeader.DiplomaticRelations[defenseLeader].Tension = 0
	defenseLeader.DiplomaticRelations[attackLeader].Tension = 0
}

func GiveWarRewards(winner *Alliance, loser *Alliance) {
	fmt.Printf(
		"%s surrenders the war to %s. %s's might increases by 1, %s's might decreases by 1.\n\n",
		loser.Leader.GetTitle(), winner.Leader.GetTitle(),
		winner.Leader.GetTitle(), loser.Leader.GetTitle(),
	)
	loser.Leader.Might = Max[int](loser.Leader.Might - 1, 1)
	winner.Leader.Might = Min[int](winner.Leader.Might + 1, MaxMight)
}
