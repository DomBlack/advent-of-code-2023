package astar

import (
	"github.com/DomBlack/advent-of-code-2023/pkg/maps"
)

type State interface {
	comparable
	Pos() maps.Pos
}
