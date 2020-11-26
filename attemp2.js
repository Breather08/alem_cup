const maxX = 12,
  maxY = 10,
  bombRadius = 2;

const map = `..........;;;
.!;!;!;!;!;!.
;..........;!
;!;!;!;!;!!!.
;..;.....;..;
.!;!;!;!;!;!;
;....;.;.;!.;
;!.!;!;!;!.!.
;.;.;....;!..
.!!!.!;!;!.!;
.;...;;;...;;`
  .split("\n")
  .map((item) => item.split(""));

function getNeighbours(current, target) {
  return [
    Tile({
      x: current.x,
      y: current.y - 1,
      cost: current.cost + 1,
      target: target,
      direction: "up"
    }),
    Tile({
      x: current.x,
      y: current.y + 1,
      cost: current.cost + 1,
      target: target,
      direction: "down"
    }),
    Tile({
      x: current.x - 1,
      y: current.y,
      cost: current.cost + 1,
      target: target,
      direction: "left"
    }),
    Tile({
      x: current.x + 1,
      y: current.y,
      cost: current.cost + 1,
      target: target,
      direction: "right"
    }),
  ];
}

function getDistance(current, target) {
  return Math.abs(target.x - current.x) + Math.abs(target.y - current.y);
}

function Tile({ x, y, cost = 0, target, parent = null, direction }) {
  let distance = 0;

  if (target) {
    distance = getDistance({ x: x, y: y }, target);
  }

  return {
    x: x,
    y: y,
    cost: cost,
    distance: distance,
    costDistance: 0,
    parent: parent,
    direction
  };
}

function buildPath(tile) {
    let path = []
    while (tile) {
        map[tile.y][tile.x] = '*'
        if (tile.direction) path.push(tile.direction)
        tile = tile.parent
    }
    console.log(path.reverse())
    map.forEach(col => console.log(col.join("")))
}

function AStar(start, finish) {
  let visitedTiles = [];
  let activeTiles = [start];

  while (activeTiles.length > 0) {
    activeTiles.sort((a, b) => a.costDistance - b.costDistance)

    let current = activeTiles.shift();

    if (current.x === finish.x && current.y === finish.y) {
      console.log("Path found");
      buildPath(current)
      return;
    }

    let neighbours = getNeighbours(current, finish).filter((nTile) => {
      if (
        nTile.x >= 0 &&
        nTile.x <= maxX &&
        nTile.y >= 0 &&
        nTile.y <= maxY &&
        map[nTile.y][nTile.x] !== "!"
      ) {
        return true;
      }
      return false;
    });

    visitedTiles.push(current);

    for (let i = 0; i < neighbours.length; i++) {
      let nTile = neighbours[i];
      if (
        visitedTiles.some((vTile) => vTile.x === nTile.x && vTile.y === nTile.y)
      ) {
        continue;
      }

      if (map[nTile.y][nTile.x] === ";") {
        nTile.cost += 10;
      }

      nTile.costDistance = nTile.cost + nTile.distance;
      nTile.parent = current;

      activeTiles.push(nTile);
    }
  }
  console.log("Failed");
  return;
}

function findSafeCoord(start) {
  function explosiveCells(bombCoords) {
    let explosives = [bombCoords];

    // Check Left
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x - bCounter,
        y: bombCoords.y,
      };

      if (0 <= explCell.x && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }

    // Check Right
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x + bCounter,
        y: bombCoords.y,
      };

      if (explCell.x <= maxX && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }

    // Check Up
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x,
        y: bombCoords.y - bCounter,
      };

      if (0 <= explCell.y && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }

    // Check Down
    for (let bCounter = 1; bCounter <= bombRadius; bCounter++) {
      let explCell = {
        x: bombCoords.x,
        y: bombCoords.y + bCounter,
      };

      if (explCell.y <= maxY && map[explCell.y][explCell.x] === ".") {
        explosives.push(explCell);
      } else {
        break;
      }
    }
    return explosives;
  }

  function checkFunc(arr, coord) {
    return !arr.some((item) => item.x === coord.x && item.y === coord.y);
  }

  function safeCoord(start) {
    const explosives = explosiveCells(start);

    let visited = [start];
    let queue = [start];
    let finish;

    while (queue.length > 0) {
      let current = queue.shift();

      const neighbours = getNeighbours(current);

      finish = current;

      if (checkFunc(explosives, finish)) {
        return finish;
      }

      neighbours.forEach((coord, i) => {
        if (
          coord.x >= 0 &&
          coord.x <= maxX &&
          coord.y >= 0 &&
          coord.y <= maxY
        ) {
          if (map[coord.y][coord.x] === "." && checkFunc(visited, coord)) {
            queue.push(coord);
            visited.push(coord);
          }
        }
      });
    }
  }

  return safeCoord(start);
}

function main() {
  const finish = Tile({ x: 7, y: 2 });
  const start = Tile({ x: 8, y: 4, target: finish });

  AStar(start, finish);
}

main();

// console.log(findSafeCoord({ x: 0, y: 1 }));
