package main

import "fmt"

type Graph struct {
	nodes []*Node
	edge  [][]*Edge
}

type Edge struct {
	me         *Node
	neighbours []*Edge
	neighDir   []Direction
	distance   int32
	dir        Direction
}

type Node struct {
	neighbours []*Node
}

func (tileStage *TileStage) MakeGraph(snake bool) *Graph {
	edges := make([][]*Edge, stageWidth)
	for x := int32(0); x < stageWidth; x++ {
		edges[x] = make([]*Edge, screenHeight)
	}
	var getEdge func(int32, int32) *Edge
	nodecount := 0
	edgeCount := 0
	getEdge = func(x, y int32) *Edge {
		if edges[x][y] == nil {
			edgeCount++
			sides := 0
			for i := Up; i <= Left; i++ {
				x2, y2 := NewPos(x, y, i)
				if tileStage.tiles[x2][y2] != Wall &&
					(!snake || tileStage.tiles[x2][y2] != SnakeWall) {
					sides++
				}
			}
			if sides > 2 {
				edges[x][y] = &Edge{&Node{nil}, make([]*Edge, sides),
					make([]Direction, sides), 0, 0}
				nodecount++
			} else {
				edges[x][y] = &Edge{nil, make([]*Edge, sides),
					make([]Direction, sides), -1, 0}
			}
			c := 0
			for i := Up; i <= Left; i++ {
				x2, y2 := NewPos(x, y, i)
				if tileStage.tiles[x2][y2] != Wall &&
					(!snake || tileStage.tiles[x2][y2] != SnakeWall) {
					edges[x][y].neighbours[c] = getEdge(x2, y2)
					edges[x][y].neighDir[c] = i
					c++
				}
			}
		}
		return edges[x][y]
	}
	edges[1][1] = getEdge(1, 1)
	fmt.Println(edgeCount, nodecount)

	var setMe func(*Edge, int32, *Node)
	setMe = func(edge *Edge, distance int32, node *Node) {
		if edge.distance == -1 || edge.distance > distance {
			edge.distance = distance
			edge.me = node
			for i := 0; i < len(edge.neighbours); i++ {
				setMe(edge.neighbours[i], distance+1, node)
			}
		} else if edge.distance == distance {
			edge.me = nil
			for i := 0; i < len(edge.neighbours); i++ {
				setMe(edge.neighbours[i], distance+1, node)
			}
		}
	}

	nodes := make([]*Node, nodecount)
	c := 0
	for x := int32(0); x < stageWidth; x++ {
		for y := int32(0); y < stageHeight; y++ {
			if edges[x][y] != nil && edges[x][y].distance == 0 {
				edges[x][y].distance = -1
				nodes[c] = edges[x][y].me
				setMe(edges[x][y], 0, nodes[c])
				c++
			}
		}
	}

	for x := int32(1); x < stageWidth-1; x++ {
		for y := int32(1); y < stageHeight-1; y++ {
			if edges[x][y] != nil && edges[x][y].me != nil &&
				edges[x][y].distance > 0 {
				for i := 0; i <= len(edges[x][y].neighbours); i++ {
					if edges[x][y].distance > edges[x][y].neighbours[i].distance {
						edges[x][y].dir = edges[x][y].neighDir[i]
						break
					}
				}
			}
		}
	}

	for x := int32(0); x < stageWidth; x++ {
		for y := int32(0); y < stageHeight; y++ {
			if edges[x][y] == nil {
				if tileStage.tiles[x][y] != Wall &&
					(!snake || tileStage.tiles[x][y] != SnakeWall) {
					fmt.Printf("Expected edge on %d,%d, got nil\n", x, y)
					panic("Position lacks edge (Not A Wall)")
				}
			} else {
				if tileStage.tiles[x][y] == Wall ||
					(snake && tileStage.tiles[x][y] == SnakeWall) {
					fmt.Printf("Illegal edge on %d,%d, in wall\n", x, y)
					panic("Illegal edge position (In A Wall)")
				} else {
					if edges[x][y].me != nil && edges[x][y].distance > 0 {
						x2, y2 := NewPos(x, y, edges[x][y].dir)
						if tileStage.tiles[x2][y2] == Wall &&
							(!snake || tileStage.tiles[x2][y2] == SnakeWall) {
							fmt.Printf("Edge pointing towards at wall from"+
								" %d,%d with %d, at %d,%d\n", x, y,
								edges[x][y].dir, x2, y2)
							panic("Illegal edge direction (Pointing At Wall)")
						}
					}
					for i := 0; i < len(edges[x][y].neighDir); i++ {
						x2, y2 := NewPos(x, y, edges[x][y].neighDir[i])
						if tileStage.tiles[x2][y2] == Wall &&
							(!snake || tileStage.tiles[x2][y2] == SnakeWall) {
							fmt.Printf("Edge with neighdir towards at wall "+
								"from %d,%d with %d (%d), at %d,%d\n", x, y,
								edges[x][y].neighDir[i], i, x2, y2)
							panic("Illegal edge neighdir (Pointing At Wall)")
						}
					}
				}
			}
		}
	}
	return &Graph{nodes, edges}
}
