package model

type Character interface {
	GetId() int

	GetCurrentHitpoints() int
	GetMaxHitpoints() int
	TakeDamage(int)

	GetPosition() *Position
	SetPosition(*Position)
	GetPrimaryDirection() Direction
	SetPrimaryDirection(Direction)
	GetSecondaryDirection() Direction
	SetSecondaryDirection(Direction)
	GetLastDirection() Direction
	SetLastDirection(Direction)
	GetLastKnownRegion() *Position
	SetLastKnownRegion(*Position)

	GetUpdateFlag() *UpdateFlag
	GetMarkedForDeletion() bool
	GetInteractingWith() Character
}
