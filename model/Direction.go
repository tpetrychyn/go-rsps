package model

//type dir = uint

type Direction = int

const (
	North = Direction(1)
	NorthEast = Direction(2)
	East = Direction(4)
	SouthEast = Direction(7)
	South = Direction(6)
	SouthWest = Direction(5)
	West = Direction(3)
	NorthWest = Direction(0)
	None = Direction(-1)
)

func DirectionFromDeltas(deltaX int, deltaY int) Direction {
	if deltaY == 1 {
		if deltaX == 1 {
			return NorthEast
		}
		if deltaX == 0 {
			return North
		}
		return NorthWest
	}
	if deltaY == -1 {
		if deltaX == 1 {
			return SouthEast
		}
		if deltaX == 0 {
			return South
		}
		return SouthWest
	}

	if deltaX == 1 {
		return East
	}
	if deltaX == -1 {
		return West
	}

	return None
}

//type Direction struct {
//	North dir
//	NorthEast dir
//	East dir
//	SouthEast dir
//	South dir
//	SouthWest dir
//	West dir
//	NorthWest dir
//}

//var Direction = &direction{
//	North: 1,
//	NorthEast: 2,
//	East: 4,
//	SouthEast: 7,
//	South: 6,
//	SouthWest: 5,
//	West: 3,
//	NorthWest: 0,
//}