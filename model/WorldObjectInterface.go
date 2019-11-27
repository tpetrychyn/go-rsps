package model

type WorldObjectInterface interface {
	GetObjectId() int
	GetPosition() *Position
	GetFace() int
	GetType() int
	Tick()
	ShouldRefresh() bool
}

type DefaultWorldObject struct {
	ObjectId int
	Position *Position
	Face     int
	Type     int
}

func (d *DefaultWorldObject) GetObjectId() int {
	return d.ObjectId
}

func (d *DefaultWorldObject) GetPosition() *Position {
	return d.Position
}

func (d *DefaultWorldObject) ShouldRefresh() bool {
	return false
}

func (d *DefaultWorldObject) Tick() {
	return
}

func (d *DefaultWorldObject) GetFace() int {
	return d.Face
}

func (d *DefaultWorldObject) GetType() int {
	return d.Type
}
