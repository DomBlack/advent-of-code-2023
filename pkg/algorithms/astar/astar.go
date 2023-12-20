package astar

import (
	goErrs "errors"
	"fmt"
	"math"

	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/heaps"
	"github.com/DomBlack/advent-of-code-2023/pkg/datastructures/vectors/vec2"
	. "github.com/DomBlack/advent-of-code-2023/pkg/maps"
	"github.com/cockroachdb/errors"
)

var ErrNoPath = goErrs.New("no path found")

// Search performs an A* search on the given map from the start position to the goal position.
//
// h is the heuristic function that returns the estimated cost from the given position to the goal.
// func Search[TileType TileWithCost](m *Map[TileType], start Pos, goal Pos, h func(from Pos) int, recordFrames bool, pathHeadFlag TileType, pathTailFlag TileType) (cost int, path []Pos, err error) {
func Search[S State, TileType TileWithCost](
	m *Map[TileType],
	start S, isGoal func(S) bool,
	neighbours func(S) []S,
	h func(from S) int,
	pathHeadFlag TileType, pathTailFlag TileType,
) (cost int, path []Pos, err error) {
	// Remove the path flags when we're done
	defer func() {
		for i := range m.Tiles {
			m.Tiles[i] = m.Tiles[i] &^ pathHeadFlag &^ pathTailFlag
		}
	}()

	nodes := make(map[S]*Node[S])

	// nodeFor returns the Node for the given position, creating it if it doesn't exist.
	nodeFor := func(state S) *Node[S] {
		if n, ok := nodes[state]; ok {
			return n
		}

		n := &Node[S]{State: state, gScore: math.MaxInt, fScore: math.MaxInt}
		nodes[state] = n
		return n
	}

	// The set of discovered nodes that may need to be (re-)expanded.
	openSet := heaps.NewMinHeap[*Node[S]](len(m.Tiles))

	// Create the starting Node
	startNode := nodeFor(start)
	startNode.gScore = 0
	startNode.fScore = h(startNode.State)
	openSet.Insert(startNode)

	for openSet.Len() > 0 {
		current := openSet.Remove()
		m.RemoveFlagAt(current.State.Pos(), pathHeadFlag)
		m.AddFlagAt(current.State.Pos(), pathTailFlag)

		// Did we get to the end?
		if isGoal(current.State) {
			return current.gScore, reconstructPath(current), nil
		}

		for _, neighbourState := range neighbours(current.State) {
			n := nodeFor(neighbourState)
			pos := neighbourState.Pos()
			tile, valid := m.Get(pos)
			if !valid {
				return 0, nil, errors.Newf("invalid neighbour position: %v", neighbourState.Pos())
			}

			tentativeGScore := current.gScore + tile.Cost()
			if tentativeGScore < n.gScore {
				n.Parent = current
				n.gScore = tentativeGScore
				n.fScore = n.gScore + h(neighbourState)
				m.AddFlagAt(pos, pathHeadFlag)

				if !openSet.Contains(n) {
					openSet.Insert(n)
				} else {
					openSet.Update(n)
				}
			}
		}

		m.CaptureFrame(fmt.Sprintf("Open States: %d", openSet.Len()), 1)
	}

	return 0, nil, errors.WithStack(ErrNoPath)
}

type DirectionalPos struct {
	Pos
	Direction vec2.Vec2
}

type TileWithCost interface {
	Tile

	// Cost returns the cost to move onto this tile on the map
	Cost() int
}

// Node represents a Node in the A* search
type Node[S State] struct {
	State  S
	gScore int      // The cost of the cheapest path from start to this Node currently known
	fScore int      // Our best guess as to the total cost from start to goal through this Node
	Parent *Node[S] // The cheapest Parent Node (nil if this is the start Node)
}

func (n *Node[S]) Less(a *Node[S]) bool {
	return n.fScore < a.fScore
}

// reconstructPath reconstructs the path from the start Node to the given Node
func reconstructPath[S State](n *Node[S]) []Pos {
	path := make([]Pos, 0)
	for n != nil {
		path = append(path, n.State.Pos())
		n = n.Parent
	}

	// Reverse the path
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}
	return path
}
