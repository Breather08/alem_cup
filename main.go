package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/bradfitz/slice"
)

type Tile struct {
	X            int
	Y            int
	distance     int
	cost         int
	CostDistance int
	Parent       *Tile
}

var targetX = 12
var targetY = 10

func (t *Tile) setCostDistance() {
	t.CostDistance = t.cost + t.distance
}

func (t *Tile) setDistance(targetX, targetY int) {
	t.distance = int(math.Abs(float64(targetX-t.X)) + math.Abs(float64(targetY-t.Y)))
}

func getPossibleTiles(gameMap []string, currentTile, targetTile *Tile) []*Tile {
	// initialize 4 direction
	tilesDir := []*Tile{
		&Tile{X: currentTile.X, Y: currentTile.Y - 1, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X, Y: currentTile.Y + 1, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X - 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile},
		&Tile{X: currentTile.X + 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile},
	}

	possible := []*Tile{}

	// initialize possible directions
	for _, tile := range tilesDir {
		tile.setDistance(targetTile.X, targetTile.Y)
		if (tile.X >= 0 && tile.X <= targetX) && (tile.Y >= 0 && tile.Y <= targetY) {
			if gameMap[tile.Y][tile.X] != '!' {
				// meeting breakable box (spends 11-12 ticks, depending on position)
				if gameMap[tile.Y][tile.X] == ';' {
					tile.cost += 10
				}
				tile.setCostDistance()
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

func sortTiles(activeTiles []*Tile) {
	slice.Sort(activeTiles, func(i, j int) bool {
		return activeTiles[i].CostDistance < activeTiles[j].CostDistance
	})
}

func makeMapArray() []string {
	fileMap, err := ioutil.ReadFile("map1.txt")
	if err != nil {
		log.Println(err)
	}
	mapStr := ""
	for _, k := range fileMap {
		mapStr += string(k)
	}

	return strings.Split(mapStr, "\n")
}

func getResult(checkTile *Tile, mapArr []string) (path []string) {
	var tile = checkTile
	fmt.Println("Retracing steps backwards...")
	for tile != nil {
		box := ""
		if mapArr[tile.Y][tile.X] == ';' {
			box = ":box"
		}
		if mapArr[tile.Y][tile.X] == '.' || mapArr[tile.Y][tile.X] == ';' {
			var newMapRow = []rune(mapArr[tile.Y])
			newMapRow[tile.X] = '*'
			mapArr[tile.Y] = string(newMapRow)
		}
		fmt.Println(tile.X, ":", tile.Y)
		if tile.Parent != nil {
			if tile.X > tile.Parent.X {
				path = append([]string{"right" + box}, path...)
			} else if tile.X < tile.Parent.X {
				path = append([]string{"left" + box}, path...)
			} else {
				if tile.Y == tile.Parent.Y {
					path = append([]string{"up" + box}, path...)
				} else {
					path = append([]string{"down" + box}, path...)
				}
			}
		}
		tile = tile.Parent
		if tile == nil {
			fmt.Println("\nMap:")
			for _, k := range mapArr {
				fmt.Println(k)
			}
			fmt.Println()
		}
	}
	return
}

func reverse(arr []string) (res []string) {
	opposites := map[string]string{
		"left":  "right",
		"right": "left",
		"up":    "down",
		"down":  "up",
	}
	for i := range arr {
		res = append(res, opposites[arr[i]])
	}
	return res
}

func concat(a, b []string) []string {
	return append(a, b...)
}

func boxPath(path []string, k string) (res []string) {
	res = append(res, "bomb")
	res = append(res, reverse(path)...)
	res = append(res, []string{"stay", "stay"}...)
	res = append(res, path...)
	res = append(res, k[:len(k)-4])
	return
}

func boxHandling(path []string) (res []string) {
	for i, k := range path {
		if k[len(k)-1] == 'x' {
			res = append(res, boxPath(path[i-3:i], k)...)
		} else {
			res = append(res, k)
		}
	}
	return
}

func cellToBoxPriority(mapArr []string, x, y int) (res []string) {
	maxX := len(mapArr[0])
	maxY := len(mapArr)
	for i := y; i < maxY; i++ {
		for j := x; j < maxX; j++ {

			if mapArr[i][j] != '.' {
				continue
			}
			boxCount := 0

			// Check right
			if j+1 < maxX && mapArr[i][j+1] == ';' {
				boxCount++
			} else if j+2 < maxX && mapArr[i][j+1] != '!' && mapArr[i][j+2] == ';' {
				boxCount++
			}

			// Check left
			if j-1 >= 0 && mapArr[i][j-1] == ';' {
				boxCount++
			} else if j-2 >= 0 && mapArr[i][j-1] != '!' && mapArr[i][j-2] == ';' {
				boxCount++
			}

			// Check down
			if i+1 < maxY && mapArr[i+1][j] == ';' {
				boxCount++
			} else if i+2 < maxY && mapArr[i+1][j] != '!' && mapArr[i+2][j] == ';' {
				boxCount++
			}

			// Check up
			if i-1 >= 0 && mapArr[i-1][j] == ';' {
				boxCount++
			} else if i-2 >= 0 && mapArr[i-1][j] != '!' && mapArr[i-2][j] == ';' {
				boxCount++
			}
			
			// Assign value to map
			mapArr[i] = strings.Replace(mapArr[i], ".", strconv.Itoa(boxCount), 1)
		}
		// fmt.Println()
	}
	return mapArr
}

func AStar() {
	mapArr := makeMapArray()
	// fmt.Println("Map: \n", )
	for _, k := range cellToBoxPriority(mapArr, 0, 0) {
		fmt.Println(strings.Join(strings.Split(k, ""), " "))
	}

	start := &Tile{
		X: 0,
		Y: 0,
	}

	finish := &Tile{
		X: 12,
		Y: 10,
	}

	start.setDistance(finish.X, finish.Y)

	activeTiles := []*Tile{start}
	visitedTiles := []*Tile{}

	path := []string{}

	for len(activeTiles) > 0 {

		// Sorting tiles in stack to chose better option
		sortTiles(activeTiles)

		// Best option
		var checkTile = activeTiles[0]

		// Bim! Printing
		if checkTile.X == finish.X && checkTile.Y == finish.Y {
			path = getResult(checkTile, mapArr)
		}

		visitedTiles = append(visitedTiles, checkTile)
		activeTiles = remove(activeTiles, checkTile)

		possible := getPossibleTiles(mapArr, checkTile, finish)

		// Creating label to be able to skip unnecessary counting
	Loop:
		for _, walkableTile := range possible {
			// Prevent entering visited cell
			for _, k := range visitedTiles {
				if k.X == walkableTile.X && k.Y == walkableTile.Y {
					continue Loop
				}
			}

			// Check if new tile has better value
			for _, k := range activeTiles {
				if k.X == walkableTile.X && k.Y == walkableTile.Y {
					existingTile := k
					if existingTile.CostDistance > walkableTile.CostDistance {
						// if so, just replace it
						activeTiles = append(activeTiles, walkableTile)
						activeTiles = remove(activeTiles, existingTile)
						continue Loop
					}
				}
			}
			activeTiles = append(activeTiles, walkableTile)
		}
	}

	fmt.Println(path)
	fmt.Println(boxHandling(path))
	fmt.Println(path)
}

func main() {
	AStar()
}
