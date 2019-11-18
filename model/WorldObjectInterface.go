package model

type WorldObjectInterface interface {
	GetObjectId() int
	GetPosition() *Position
	Tick()
	ShouldRefresh() bool
}
