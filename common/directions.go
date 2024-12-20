package common

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type CardinalDirection int

const (
	North CardinalDirection = iota
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)

var (
	CardinalDirections [8]CardinalDirection = [8]CardinalDirection{North, South, East, West, NorthEast, NorthWest, SouthEast, SouthWest}
)

func (cd CardinalDirection) GetCoord() Coord {
	switch cd {
	case North:
		return Coord{X: 0, Y: -1}
	case South:
		return Coord{X: 0, Y: 1}
	case East:
		return Coord{X: 1, Y: 0}
	case West:
		return Coord{X: -1, Y: 0}
	case NorthEast:
		return Coord{X: 1, Y: -1}
	case NorthWest:
		return Coord{X: -1, Y: -1}
	case SouthEast:
		return Coord{X: 1, Y: 1}
	case SouthWest:
		return Coord{X: -1, Y: 1}
	}
	return Coord{}
}

func (cd CardinalDirection) String() string {
	switch cd {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	case NorthEast:
		return "NorthEast"
	case NorthWest:
		return "NorthWest"
	case SouthEast:
		return "SouthEast"
	case SouthWest:
		return "SouthWest"
	}
	return ""
}
