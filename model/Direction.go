package model

type dir = uint

type direction struct {
	North dir
	NorthEast dir
	East dir
	SouthEast dir
	South dir
	SouthWest dir
	West dir
	NorthWest dir
}

var Direction = &direction{
	North: 1,
	NorthEast: 2,
	East: 4,
	SouthEast: 7,
	South: 6,
	SouthWest: 5,
	West: 3,
	NorthWest: 0,
}