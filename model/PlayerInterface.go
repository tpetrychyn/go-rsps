package model

type PlayerInterface interface {
	Character

	GetEquipmentItemContainer() *ItemContainer

	GetNearbyPlayers() []PlayerInterface
	GetLoadedPlayers() []PlayerInterface
	AddLoadedPlayer(PlayerInterface)
	RemoveLoadedPlayer(int)

	GetNearbyNpcs() map[int]NpcInterface
	GetLoadedNpcs() map[int]NpcInterface
	AddLoadedNpc(NpcInterface)
	RemoveLoadedNpc(int)

	GetName() string
}
