package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strings"
)

type Coords struct {
	x int
	y int
}

type Entity struct {
	name string
	id   int
	Coords
	param1 int
	param2 int
}

type Tile struct {
	Coords
	direction    string
	distance     int
	cost         int
	CostDistance int
	retreatPath  []string
	isBox        bool
	Parent       *Tile
}

const (
	maxX       = 12
	maxY       = 10
	bombRadius = 2
)

func (t *Tile) setCostDistance() {
	t.CostDistance = t.cost + t.distance
}

func (t *Tile) setDistance(targetX, targetY int) {
	t.distance = int(math.Abs(float64(targetX-t.x)) + math.Abs(float64(targetY-t.y)))
}

func getResult(checkTile *Tile, mapArr []string) (path []string) {
	var tile = checkTile
	counter := 0
	for tile != nil {
		counter++
		if len(tile.retreatPath) > 0 {
			if len(path) > 1 {
				path = path[:len(path)-1]
			}
			path = append(tile.retreatPath, path...)
		}

		if tile.direction != "" {
			path = append([]string{tile.direction}, path...)
		}

		if mapArr[tile.y][tile.x] != '!' {
			var newMapRow = []rune(mapArr[tile.y])
			newMapRow[tile.x] = '*'
			mapArr[tile.y] = string(newMapRow)
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

func remove(s []*Tile, target *Tile) (res []*Tile) {
	for _, k := range s {
		if k != target {
			res = append(res, k)
		}
	}
	return
}

func sortTiles(activeTiles []*Tile) {
	sort.Slice(activeTiles, func(i, j int) bool {
		return activeTiles[i].CostDistance > activeTiles[j].CostDistance
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

func getPossibleTiles(gameMap []string, currentTile, targetTile *Tile, visitedTiles []*Tile) []*Tile {
	// initialize 4 direction
	tilesDir := []*Tile{
		&Tile{Coords: Coords{x: currentTile.x, y: currentTile.y - 1}, cost: currentTile.cost + 1, Parent: currentTile, direction: "up"},
		&Tile{Coords: Coords{x: currentTile.x, y: currentTile.y + 1}, cost: currentTile.cost + 1, Parent: currentTile, direction: "down"},
		&Tile{Coords: Coords{x: currentTile.x - 1, y: currentTile.y}, cost: currentTile.cost + 1, Parent: currentTile, direction: "left"},
		&Tile{Coords: Coords{x: currentTile.x + 1, y: currentTile.y}, cost: currentTile.cost + 1, Parent: currentTile, direction: "right"},
	}

	possible := []*Tile{}

	// initialize possible directions
Loop:
	for _, tile := range tilesDir {
		tile.setDistance(targetTile.x, targetTile.y)
		if (tile.x >= 0 && tile.x <= maxX) && (tile.y >= 0 && tile.y <= maxY) {
			if gameMap[tile.y][tile.x] == '.' || gameMap[tile.y][tile.x] == ';' {

				// Prevent entering visited cell
				for _, k := range visitedTiles {
					if k.x == tile.x && k.y == tile.y {
						continue Loop
					}
				}

				// meeting breakable box (spends 11-12 ticks, depending on position)
				if gameMap[tile.y][tile.x] == ';' {
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

func (t *Tile) setRetreatPath(mapArr []string) {

	retreatCoord := bfs(mapArr, Coords{x: t.Parent.x, y: t.Parent.y})
	fmt.Println(retreatCoord)
	return
	retreatPath := AStar(mapArr, retreatCoord)
	retreat := append([]string{"bomb"}, retreatPath...)
	retreat = append(retreat, []string{"stay", "stay"}...)
	retreat = append(retreat, reverse(retreat)...)
	retreat = append(retreat, t.direction)
	t.retreatPath = retreat

}

func replaceAtIndex(in string, r rune, i int) string {
    out := []rune(in)
    out[i] = r
    return string(out)
}

func updateMap(gameMap []string, bombCell Coords) {
	gameMap[bombCell.y] = replaceAtIndex(gameMap[bombCell.y], '.', bombCell.x)
}

func AStar(mapArr []string, startCoords Coords) (path []string) {

	checkTile := &Tile{}

	start := &Tile{
		Coords: startCoords,
	}

	finish := &Tile{
		Coords: Coords{
			x: maxX,
			y: maxY,
		},
	}

	start.setDistance(finish.x, finish.y)

	activeTiles := []*Tile{start}
	visitedTiles := []*Tile{}

	for len(activeTiles) > 0 {

		// Sorting tiles in stack to chose better option
		sortTiles(activeTiles)

		// Best option
		checkTile = activeTiles[0]
		
		// Set Bomb and update map
		if checkTile.isBox {
			checkTile.setRetreatPath(mapArr)
			updateMap(mapArr, checkTile.Coords)
		}

		if checkTile.x == finish.x && checkTile.y == finish.y {
			path = getResult(checkTile, mapArr)
			return
		}

		visitedTiles = append(visitedTiles, checkTile)
		activeTiles = remove(activeTiles, checkTile)

		possible := getPossibleTiles(mapArr, checkTile, finish, visitedTiles)

		for _, walkableTile := range possible {
			activeTiles = append(activeTiles, walkableTile)
		}
	}

	path = getResult(checkTile, mapArr)

	fmt.Println(path)
	return
}

func main() {
	mapArr := makeMapArray()
	AStar(mapArr, Coords{x: 0, y: 0})
}

func explosionArea(gameMap []string, bombC Coords) map[Coords]bool {
	explosionArea := map[Coords]bool{
		bombC: true,
	}
	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x - bCounter, bombC.y}
		if 0 <= explCell.x && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x + bCounter, bombC.y}
		if explCell.x <= maxX && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x, bombC.y - bCounter}
		if 0 <= explCell.y && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	for bCounter := 1; bCounter <= bombRadius; bCounter++ {
		explCell := Coords{bombC.x, bombC.y + bCounter}
		if explCell.y <= maxY && gameMap[explCell.y][explCell.x] == '.' {
			explosionArea[explCell] = true
		} else {
			break
		}
	}

	return explosionArea
}

func bfs(gameMap []string, start Coords) (finish Coords) {
	fmt.Println("Current start point: ", start)
	visited := map[Coords]bool{
		start: true,
	}

	explosive := explosionArea(gameMap, start)
	fmt.Print("Bomb Hitting Cells: ", explosive)

	queue := []Coords{start}

	counter := 0
	for len(queue) > 0 {

		fmt.Println("Queue: ", queue)
		fmt.Println("Visited Coordinates: ", visited)
		counter++

		coordsDir := []Coords{
			Coords{x: queue[0].x, y: queue[0].y - 1},
			Coords{x: queue[0].x, y: queue[0].y + 1},
			Coords{x: queue[0].x - 1, y: queue[0].y},
			Coords{x: queue[0].x + 1, y: queue[0].y},
		}

		finish = queue[0]

		if !explosive[finish] {
			return finish
		}

		fmt.Println("Current optimal coord: ", finish)

		queue = queue[1:]

		for _, coord := range coordsDir {
			if (coord.x >= 0 && coord.x <= maxX) && (coord.y >= 0 && coord.y <= maxY) {
				if gameMap[coord.y][coord.x] == '.' && !visited[coord] {
					queue = append(queue, coord)
					visited[coord] = true
				}
			}
		}

	}

	fmt.Println("Finishing at: ", finish)

	return
}

// func bfs(start Coords, nodes map[int][]int, fn func (int)) {
//     frontier := []Coords{start}
//     visited := map[Coords]bool{}
//     next := []Coords{}

//     for 0 < len(frontier) {
//         next = []Coords{}
//         for _, node := range frontier {
//             visited[node] = true
//             for _, n := range bfs_possibles(node, nodes, visited) {
//                 next = append(next, n)
//             }
//         }
//         frontier = next
//     }
// }

// func bfs_possibles(active []Coords, node Coords, visited map[Coords]bool) []Coords {
//     next := []Coords{}
//     iter := func (n Coords) bool { _, ok := visited[n]; return !ok }
//     for _, n := range gameMap {
//         if iter(n) {
//             next = append(next, n)
//         }
//     }
//     return next
// }

// func makeMapArray() (res [][]byte) {
// 	fileMap, err := ioutil.ReadFile("map1.txt")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	temp := []byte{}
// 	for _, k := range fileMap {
// 		if k != 10 {
// 			temp = append(temp, k)
// 		} else {
// 			res = append(res, temp)
// 			temp = []byte{}
// 		}
// 	}

// 	return
// }
