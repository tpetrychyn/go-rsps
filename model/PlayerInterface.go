package model

type PlayerInterface interface {
	Character

	GetEquipmentItemContainer() *ItemContainer

	GetNearbyPlayers() []PlayerInterface
	GetLoadedPlayers() []PlayerInterface
	AddLoadedPlayer(PlayerInterface)
	RemoveLoadedPlayer(int)

	GetNearbyNpcs() []NpcInterface
	GetLoadedNpcs() []NpcInterface
	AddLoadedNpc(NpcInterface)
	RemoveLoadedNpc(int)

	GetName() string
}
