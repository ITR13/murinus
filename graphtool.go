package main

import (
	"fmt"
)

const (
	DisplayStageGraph        = false
	ShowDistanceInsteadOfDir = false
)

type Graph struct {
	nodes [][]*Node
}

type Node struct {
	sides [4]*Side
}

type Side struct {
	dirToPush Direction
	distance  int
}

func (tileStage *TileStage) MakeGraph(snake bool) *Graph {
	fmt.Println("Making node array")
	nodes := make([][]*Node, stageWidth)
	for x := int32(0); x < stageWidth; x++ {
		nodes[x] = make([]*Node, screenHeight)
	}

	nodeCount := 0
	sideCount := 0

	for x := int32(0); x < stageWidth; x++ {
		for y := int32(0); y < stageHeight; y++ {
			if tileStage.tiles[x][y] != Wall &&
				(!snake || tileStage.tiles[x][y] != SnakeWall) {
				nodes[x][y] = &Node{}
				nodeCount++
			}
		}
	}

	for x := int32(0); x < stageWidth; x++ {
		for y := int32(0); y < stageHeight; y++ {
			if nodes[x][y] != nil {
				for i := Up; i <= Left; i++ {
					x2, y2 := NewPos(x, y, i)
					if nodes[x2][y2] != nil {
						sideCount++
						nodes[x][y].sides[i] = &Side{i, 0}
					}
				}
			}
		}
	}

	fmt.Printf("The stagegraph has %d nodes, and %d sides\n",
		nodeCount, sideCount)

	//TODO Try to find a way to do this more efficiently
	//TODO (also more beautifully)
	setSide := func(dir Direction, reverse bool) {
		max1, max2 := int(stageHeight), int(stageWidth)
		upDown := dir == Left || dir == Right
		dirToSet := Right
		if upDown {
			dirToSet++
			max1, max2 = max2, max1
		}
		if !reverse {
			dirToSet = (dirToSet + 2) % 4
		}

		for i := 0; i < max1; i++ {
			distance := -1
			for j := 0; j < max2; j++ {
				var node *Node
				if !upDown {
					if reverse {
						node = nodes[max2-j-1][i]
					} else {
						node = nodes[j][i]
					}
				} else {
					if reverse {
						node = nodes[i][max2-j-1]
					} else {
						node = nodes[i][j]
					}
				}
				if node != nil {
					side := node.sides[dir]
					if side == nil {
						if distance != -1 {
							distance++
							node.sides[dir] = &Side{dirToSet, distance}
						}
					} else {
						if side.distance == 0 {
							distance = 0
						} else if distance != -1 {
							distance++
							if side.distance > distance {
								side.dirToPush = dirToSet
								side.distance = distance
							} else {
								if side.distance == distance {
									side.dirToPush = 5
								}
								distance = -1
							}
						}
					}
				} else {
					distance = -1
				}
			}
		}
	}

	for i := Up; i <= Left; i++ {
		setSide(i, false)
		setSide(i, true)
	}

	for x := int32(0); x < stageWidth; x++ {
		for y := int32(0); y < stageHeight; y++ {
			if nodes[x][y] == nil {
				if tileStage.tiles[x][y] != Wall &&
					(!snake || tileStage.tiles[x][y] != SnakeWall) {
					fmt.Printf("Expected node on %d,%d, got nil\n", x, y)
					panic("Position lacks node (Not A Wall)")
				}
			} else {
				if tileStage.tiles[x][y] == Wall ||
					(snake && tileStage.tiles[x][y] == SnakeWall) {
					fmt.Printf("Illegal node on %d,%d, in wall\n", x, y)
					panic("Illegal node position (In A Wall)")
				} else {
					for i := Up; i <= Left; i++ {
						side := nodes[x][y].sides[i]
						if side == nil {
							//nodes[x][y].sides[i] = &Side{i, 2}
						}
					}
				}
			}
		}
	}

	for i := Up; i <= Left && DisplayStageGraph; i++ {
		fmt.Println("----------------------")
		for y := int32(0); y < stageHeight; y++ {
			for x := int32(0); x < stageWidth; x++ {
				node := nodes[x][y]
				if node == nil {
					fmt.Print("#")
				} else {
					side := node.sides[i]
					if side == nil {
						fmt.Print("*")
					} else {
						if ShowDistanceInsteadOfDir {
							fmt.Print(side.distance)
						} else {
							if side.dirToPush == Up {
								fmt.Print("↑")
							} else if side.dirToPush == Right {
								fmt.Print("→")
							} else if side.dirToPush == Down {
								fmt.Print("↓")
							} else if side.dirToPush == Left {
								fmt.Print("←")
							} else {
								fmt.Print(side.dirToPush)
							}
						}
					}
				}
			}
			fmt.Println()
		}
	}

	return &Graph{nodes}
}
