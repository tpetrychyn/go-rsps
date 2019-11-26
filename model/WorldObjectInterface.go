package model

type WorldObjectInterface interface {
	GetObjectId() int
	GetPosition() *Position
	Tick()
	ShouldRefresh() bool
}

type DefaultWorldObject struct {
	ObjectId int
	Position *Position
}

func (d *DefaultWorldObject) GetObjectId() int {
	return d.ObjectId
}

func (d *DefaultWorldObject) GetPosition() *Position {
	return d.Position
}

func (d *DefaultWorldObject) ShouldRefresh() bool {
	return true
}

func (d *DefaultWorldObject) Tick() {
	return
}
