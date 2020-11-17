package main

import (
	"fmt"
	"github.com/bradfitz/slice"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

type Tile struct {
	X            int
	Y            int
	distance     int
	cost         int
	CostDistance int
	Parent       *Tile
}

var maxX = 12
var maxY = 10

func (t *Tile) setCostDistance() {
	t.CostDistance = t.cost + t.distance
}
// Manhattan
func (t *Tile) setDistance(targetX, targetY int) {
	t.distance = int(math.Abs(float64(targetX-t.X)) + math.Abs(float64(targetY-t.Y)))
}

func getPossibleTiles(gameMap []string, currentTile, targetTile *Tile) []*Tile {
	// currentTile.setCostDistance()
	tilesDir := []*Tile{
		&Tile{X: currentTile.X, Y: currentTile.Y - 1, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X, Y: currentTile.Y + 1, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X - 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X + 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile},
	}

	possible := []*Tile{}

	for _, tile := range tilesDir {
		tile.setDistance(targetTile.X, targetTile.Y)
		if (tile.X >= 0 && tile.X <= maxX) && (tile.Y >= 0 && tile.Y <= maxY) {
			if gameMap[tile.Y][tile.X] == '.' || gameMap[tile.Y][tile.X] == ';' {
				possible = append(possible, tile)
			}
		}
	}

	return possible
}

func remove(s []*Tile, target *Tile) (res []*Tile) {
	for _, k := range s {
		if k != target {
			res = append(res, k)
		}
	}
	return
}

func AStar() {
	fileMap, err := ioutil.ReadFile("map1.txt")
	if err != nil {
		log.Println(err)
	}
	mapStr := ""
	for _, k := range fileMap {
		mapStr += string(k)
	}

	mapArr := strings.Split(mapStr, "\n")

	start := &Tile{
		X: 0,
		Y: 0,
	}

	finish := &Tile{
		X: 12,
		Y: 10,
	}

	// start.setDistance(finish.X, finish.Y)

	activeTiles := []*Tile{start}
	visitedTiles := []*Tile{}

	itest := 0

	for len(activeTiles) > 0 {
		itest++
		slice.Sort(activeTiles, func(i, j int) bool {
			return activeTiles[i].CostDistance < activeTiles[j].CostDistance
		})

		var checkTile = activeTiles[0]
		// fmt.Println(checkTile.X, ":", checkTile.Y)
		
		if checkTile.X == finish.X && checkTile.Y == finish.Y {
			var tile = checkTile
			fmt.Println("Retracing steps backwards...")
			for true {
				fmt.Println(tile.X, ":", tile.Y)
				// if mapArr[tile.Y][tile.X] == '.' || mapArr[tile.Y][tile.X] == ';' {
					var newMapRow = []rune(mapArr[tile.Y])
					newMapRow[tile.X] = '*'
					mapArr[tile.Y] = string(newMapRow)
				// }
				tile = tile.Parent
				if tile == nil {
					// fmt.Println("Map:")
					for _,k := range mapArr {
						fmt.Println(k)
					}
				}
			}
			return
		}

		visitedTiles = append(visitedTiles, checkTile)
		activeTiles = remove(activeTiles, checkTile)

		possible := getPossibleTiles(mapArr, checkTile, finish)
	Loop:
		for _, walkableTile := range possible {
			for _, k := range visitedTiles {
				if k.X == walkableTile.X && k.Y == walkableTile.Y {
					continue Loop
				}
			}

			for _, k := range activeTiles {
				if k.X == walkableTile.X && k.Y == walkableTile.Y {
					existingTile := k
					if existingTile.CostDistance < checkTile.CostDistance {
						fmt.Println("hi")
						activeTiles = remove(activeTiles, existingTile)
						activeTiles = append(activeTiles, walkableTile)
						continue Loop
					}
				}
			}
			activeTiles = append(activeTiles, walkableTile)
		}

	}
}

func main() {
	AStar()
}
