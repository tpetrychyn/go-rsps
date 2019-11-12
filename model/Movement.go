package model

type Movement struct {
	Position           *Position
	LastKnownRegion    *Position
	PrimaryDirection   Direction
	SecondaryDirection Direction
	LastDirection      Direction
	IsRunning          bool
}
