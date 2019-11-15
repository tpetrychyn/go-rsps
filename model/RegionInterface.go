package model

import "github.com/google/uuid"

type RegionInterface interface {
	GetPlayers() map[uuid.UUID]PlayerInterface
}
