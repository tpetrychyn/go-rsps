package model

import "github.com/google/uuid"

type PlayerInterface interface {
	GetPlayerId() uuid.UUID
	GetLastDirection() Direction
	GetPrimaryDirection() Direction
	GetSecondaryDirection() Direction
	GetPosition() *Position
	GetLastKnownRegion() *Position
	GetEquipmentItemContainer() *ItemContainer
	GetUpdateFlag() *UpdateFlag

	GetNearbyPlayers() []PlayerInterface
	GetLoadedPlayers() map[uuid.UUID]PlayerInterface
	AddLoadedPlayer(PlayerInterface)
	RemoveLoadedPlayer(uuid.UUID)

	GetName() string
}
