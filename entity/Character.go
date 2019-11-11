package entity

import "rsps/model"

type Character interface{
	GetPosition() *model.Position
	SetPosition(*model.Position)
	GetPrimaryDirection() model.Direction
	SetPrimaryDirection(model.Direction)
	GetSecondaryDirection() model.Direction
	SetSecondaryDirection(model.Direction)
	GetLastDirection() model.Direction
	SetLastDirection(model.Direction)
	GetLastKnownRegion() *model.Position
	SetLastKnownRegion(*model.Position)
}