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
	direction    string
	isBox        bool
	retreatPath  []string
	Parent       *Tile
}

func (t *Tile) setCostDistance() {
	t.CostDistance = t.cost + t.distance
}

func (t *Tile) setDistance(targetX, targetY int) {
	t.distance = int(math.Abs(float64(targetX-t.X)) + math.Abs(float64(targetY-t.Y)))
}

func removeTile(s []*Tile, target *Tile) (res []*Tile) {
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

func getPossibleTiles(gameMap []string, currentTile, targetTile *Tile) []*Tile {
	// initialize 4 direction
	tilesDir := []*Tile{
		&Tile{X: currentTile.X, Y: currentTile.Y - 1, cost: currentTile.cost + 1, Parent: currentTile, direction: "up"},
		&Tile{X: currentTile.X, Y: currentTile.Y + 1, cost: currentTile.cost + 1, Parent: currentTile, direction: "down"},
		&Tile{X: currentTile.X - 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile, direction: "left"},
		&Tile{X: currentTile.X + 1, Y: currentTile.Y, cost: currentTile.cost + 1, Parent: currentTile, direction: "right"},
	}

	possible := []*Tile{}

	// initialize possible directions
	for _, tile := range tilesDir {
		tile.setDistance(targetTile.X, targetTile.Y)
		if (tile.X >= 0 && tile.X <= targetTile.X) && (tile.Y >= 0 && tile.Y <= targetTile.Y) {
			if gameMap[tile.Y][tile.X] != '!' {
				// meeting breakable box (spends 11-12 ticks, depending on position)
				if gameMap[tile.Y][tile.X] == ';' {
					tile.isBox = true
					tile.cost += 10
				}
				tile.setCostDistance()
				possible = append(possible, tile)
			}
		}
	}

	return possible
}

func getResult(checkTile *Tile, mapArr []string) (path []string) {
	var tile = checkTile
	for tile != nil {

		fmt.Println(tile)
		if len(tile.retreatPath) > 0 {
			path = append(tile.retreatPath, path...)
		}

		if tile.direction != "" {
			path = append([]string{tile.direction}, path...)
		}

		if mapArr[tile.Y][tile.X] != '!' {
			var newMapRow = []rune(mapArr[tile.Y])
			newMapRow[tile.X] = '*'
			mapArr[tile.Y] = string(newMapRow)
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
		if opposites[arr[i]] != "" {
			res = append([]string{opposites[arr[i]]}, res...)
		}
	}
	return res
}

func boxPath(path []string, retreat []string, k string) (res []string) {
	retreat = reverse(path)
	res = append(res, "bomb")
	res = append(res, retreat...)
	res = append(res, []string{"stay", "stay"}...)
	res = append(res, path...)
	res = append(res, k[:len(k)-4])
	return
}

func boxHandling(path []string) (res []string) {
	for i, k := range path {
		if k[len(k)-1] == 'x' {
			res = append(res, boxPath(path[i-3:i], []string{}, k)...)
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

func (t *Tile) setRetreatPath(mapArr []string) {
	retreatPath := AStar(mapArr, t.Parent.X, t.Parent.Y, true)
	retreat := append([]string{"bomb"}, retreatPath...)
	retreat = append(retreat, []string{"stay", "stay"}...)
	retreat = append(retreat, reverse(retreat)...)
	retreat = append(retreat, t.direction)
	t.retreatPath = retreat
	t.isBox = false
}

func AStar(mapArr []string, startX, startY int, isBox bool) (path []string) {
	// fmt.Println("Map: \n", )
	// for _, k := range cellToBoxPriority(mapArr, 0, 0) {
	// 	fmt.Println(strings.Join(strings.Split(k, ""), " "))
	// }

	blockedMap := []string{}

	for _, k := range mapArr {
		blockedMap = append(blockedMap, strings.ReplaceAll(k, ";", "!"))
	}

	start := &Tile{
		X: startX,
		Y: startY,
	}

	finish := &Tile{
		X: 12,
		Y: 10,
	}

	checkTile := &Tile{}

	start.setDistance(finish.X, finish.Y)

	activeTiles := []*Tile{start}
	visitedTiles := []*Tile{}
	counter := 0
	for len(activeTiles) > 0 {
		counter++

		// Sorting tiles in stack to chose better option
		sortTiles(activeTiles)

		// Best option
		checkTile = activeTiles[0]

		if checkTile.isBox {
			checkTile.setRetreatPath(blockedMap)
		}

		// Bim! Printing
		if checkTile.X == finish.X && checkTile.Y == finish.Y {
			path = getResult(checkTile, mapArr)
			return
		}

		visitedTiles = append(visitedTiles, checkTile)
		activeTiles = removeTile(activeTiles, checkTile)

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
						activeTiles = removeTile(activeTiles, existingTile)
						continue Loop
					}
				}
			}
			activeTiles = append(activeTiles, walkableTile)
		}
	}

	// return getResult(checkTile, mapArr)
	return
}

func main() {
	mapArr := makeMapArray()
	path := AStar(mapArr, 0, 0, false)
	fmt.Println(path)
}