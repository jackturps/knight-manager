package game

import "fmt"

type WorldEventFunc = func()

func HouseAnnoysHouseEvent(flavourText string, tensionAmount int) {
	sourceHouse := RandomSelect(Game.Houses)
	possibleTargets := RemoveItem(Game.Houses, sourceHouse)
	targetHouse := RandomSelect(possibleTargets)
	targetHouse.DiplomaticRelations[sourceHouse].Tension += tensionAmount
	currentTension := targetHouse.DiplomaticRelations[sourceHouse].Tension
	fmt.Printf(flavourText + " Tensions increased to %d.\n", sourceHouse.GetTitle(), targetHouse.GetTitle(), currentTension)
}

// TODO: Maybe give houses a stat for how likely they are to antagonise others? Tyranny or something?
// TODO: Make more personal events, knights killing other knights etc.
var WorldEvents = []WorldEventFunc{
	func() { HouseAnnoysHouseEvent("%s imposed a trade embargo on %s.", 2) },
	func() { HouseAnnoysHouseEvent("%s raided a village in %s's lands.", 3) },
	func() { HouseAnnoysHouseEvent("A %s noble offended a %s noble during a feast.", 1) },
	func() { HouseAnnoysHouseEvent("A %s noble had a %s noble assassinated.", 3) },
	func() { HouseAnnoysHouseEvent("%s is blackmailing %s.", 2) },
	func() { HouseAnnoysHouseEvent("A %s noble killed a %s noble in a duel.", 2) },
	func() { HouseAnnoysHouseEvent("A %s noble started a brawl with a %s noble during a feast.", 1) },
	func() { HouseAnnoysHouseEvent("%s imposed tolls on all roads leading to %s's lands.", 1) },
	func() { HouseAnnoysHouseEvent("%s deployed a garrison on %s's border.", 2) },
}

func DoWorldEvent() {
	worldEvent := RandomSelect(WorldEvents)
	worldEvent()
}