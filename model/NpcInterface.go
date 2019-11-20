package model

import "github.com/google/uuid"

type NpcInterface interface {
	GetId() uuid.UUID
	GetNpcId() int
	GetNpcType() int

	GetPrimaryDirection() Direction
	GetSecondaryDirection() Direction
	GetUpdateFlag() *UpdateFlag

	GetPosition() *Position
}
