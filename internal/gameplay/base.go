package gameplay

// type Room struct {
// 	x1, y1, x2, y2 int
// }

// type ConnectionType int

// const (
// 	ConnectionTypeWall = iota
// 	ConnectionTypeDoor
// 	ConnectionTypeHole
// )

// type Connection struct {
// 	x1, y1, x2, y2 int
// 	connectionType ConnectionType
// }

// func buildBase(tileMap *TileMap, rooms []Room) {
// 	tileMap.ClearAll()

// 	connections := []Connection{}

// 	for _, room := range rooms {
// 		connections = append(connections,
// 			Connection{room.x1, room.y1, room.x2, room.y1, ConnectionTypeWall}, // n
// 			Connection{room.x1, room.y2, room.x2, room.y2, ConnectionTypeWall}, // s
// 			Connection{room.x2, room.y1, room.x2, room.y2, ConnectionTypeWall}, // e
// 			Connection{room.x1, room.y1, room.x1, room.y2, ConnectionTypeWall}, // w
// 		)
// 	}

// 	for _, room := range rooms {
// 		buildRoomStructureBits(tileMap, connections, room)
// 	}
// 	for _, room := range rooms {
// 		buildRoomAestheticBits(tileMap, room)
// 	}
// }

// func buildRoomStructureBits(tileMap *TileMap, connections []Connection, room Room) {
// 	setFloorBits(tileMap, connections, room)
// 	setWallsBits(tileMap, connections, room)
// 	setDoorsBits(tileMap, connections, room)
// }

// func setFloorBits(tileMap *TileMap, connections []Connection, room Room) {
// 	for i := 0; i < (room.x2-room.x1)/tileMap.gridSize; i++ {
// 		for j := 0; j < (room.y2-room.y1)/tileMap.gridSize; j++ {
// 			col := (room.x1 / tileMap.gridSize) + i
// 			row := (room.y1 / tileMap.gridSize) + j

// 			tileMap.SetBit(row, col, FLOOR_BIT)
// 		}
// 	}
// }

// func setWallsBits(tileMap *TileMap, connections []Connection, room Room) {
// 	for _, connection := range connections {
// 		if connection.connectionType == ConnectionTypeWall {
// 			if connection.x2-connection.x1 == 0 {
// 				// TODO - should be grid size?
// 				for i := (connection.y1 / 64); i < connection.y2/64; i++ {
// 					tileMap.SetBit(i, (connection.x1 / 64), VWALL_BIT)
// 				}
// 			} else {
// 				for i := (connection.x1 / 64); i < connection.x2/64; i++ {
// 					tileMap.SetBit((connection.y1/64)-1, i, HWALL_BIT)
// 				}
// 			}
// 		}
// 	}
// }

// func setDoorsBits(tileMap *TileMap, connections []Connection, room Room) {
// 	// TODO
// }

// // private void setDoorsBits(ConnectivityGraph graph, Room room) {
// // 	for (Connection connection : graph.getConnections(room)) {
// // 		if (connection.getType() == Connection.ConnectionType.DOOR) {
// // 			if (connection.x2 - connection.x1 == 0) {
// // 				for (int i = (int) (connection.y1 / 64); i < connection.y2 / 64; i++) {
// // 					setBit((int) (connection.x1 / 64), i, VDOOR_BIT);
// // 				}
// // 			} else {
// // 				for (int i = (int) (connection.x1 / 64); i < connection.x2 / 64; i++) {
// // 					setBit(i, (int) (connection.y1 / 64) - 1, HDOOR_BIT);
// // 				}
// // 			}
// // 		}
// // 	}
// // }

// func buildRoomAestheticBits(tileMap *TileMap, room Room) {
// 	for i := -1; i <= (room.x2-room.x1)/tileMap.gridSize; i++ {
// 		for j := -1; j <= (room.y2-room.y1)/tileMap.gridSize; j++ {
// 			col := (room.x1 / tileMap.gridSize) + i
// 			row := (room.y1 / tileMap.gridSize) + j

// 			buildWallTiles(tileMap, row, col)
// 			buildCornerTiles(tileMap, row, col)
// 			buildTerminusTiles(tileMap, row, col)
// 			buildDoorTiles(tileMap, row, col)
// 		}
// 	}
// }

// func buildWallTiles(tileMap *TileMap, row, col int) {
// 	// horizontal

// 	if tileMap.GetBit(row, col, HWALL_BIT) {
// 		tileMap.SetBit(row, col, INTERIOR_WALL_N_BIT)
// 		tileMap.SetBit(row+1, col, INTERIOR_WALL_S_BIT)

// 		if !tileMap.GetBit(row, col, FLOOR_BIT) {
// 			tileMap.SetBit(row, col, EXTERIOR_WALL_N_BIT)
// 		}

// 		if !tileMap.GetBit(row+1, col, FLOOR_BIT) {
// 			tileMap.SetBit(row+1, col, EXTERIOR_WALL_S_BIT)
// 		}
// 	}

// 	// vertical

// 	if tileMap.GetBit(row, col, VWALL_BIT) {
// 		tileMap.SetBit(row, col, INTERIOR_WALL_E_BIT)
// 		tileMap.SetBit(row, col-1, INTERIOR_WALL_W_BIT)

// 		if !tileMap.GetBit(row, col, FLOOR_BIT) {
// 			tileMap.SetBit(row, col, EXTERIOR_WALL_E_BIT)
// 		}

// 		if !tileMap.GetBit(row, col-1, FLOOR_BIT) {
// 			tileMap.SetBit(row, col-1, EXTERIOR_WALL_W_BIT)
// 		}
// 	}
// }

// func buildCornerTiles(tileMap *TileMap, row, col int) {
// 	// convex

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row+1, col+1, VWALL_BIT) && !tileMap.GetBit(row, col, FLOOR_BIT) && !tileMap.GetBit(row+1, col+1, FLOOR_BIT) {
// 		tileMap.SetBit(row, col+1, EXTERIOR_CORNER_CONVEX_NE_BIT)
// 	}

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row+1, col, VWALL_BIT) && !tileMap.GetBit(row+1, col-1, FLOOR_BIT) && !tileMap.GetBit(row, col, FLOOR_BIT) {
// 		tileMap.SetBit(row, col-1, EXTERIOR_CORNER_CONVEX_NW_BIT)
// 	}

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row, col+1, VWALL_BIT) && !tileMap.GetBit(row+1, col, FLOOR_BIT) && !tileMap.GetBit(row, col+1, FLOOR_BIT) {
// 		tileMap.SetBit(row+1, col+1, EXTERIOR_CORNER_CONVEX_SE_BIT)
// 	}

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row, col, VWALL_BIT) && !tileMap.GetBit(row, col-1, FLOOR_BIT) && !tileMap.GetBit(row+1, col, FLOOR_BIT) {
// 		tileMap.SetBit(row+1, col-1, EXTERIOR_CORNER_CONVEX_SW_BIT)
// 	}

// 	// concave

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row+1, col+1, VWALL_BIT) && !tileMap.GetBit(row+1, col, FLOOR_BIT) {
// 		tileMap.SetBit(row+1, col, EXTERIOR_CORNER_CONCAVE_NE_BIT)
// 	}

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row+1, col, VWALL_BIT) && !tileMap.GetBit(row+1, col, FLOOR_BIT) {
// 		tileMap.SetBit(row+1, col, EXTERIOR_CORNER_CONCAVE_NW_BIT)
// 	}

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row, col+1, VWALL_BIT) && !tileMap.GetBit(row, col, FLOOR_BIT) {
// 		tileMap.SetBit(row, col, EXTERIOR_CORNER_CONCAVE_SE_BIT)
// 	}

// 	if tileMap.GetBit(row, col, HWALL_BIT) && tileMap.GetBit(row, col, VWALL_BIT) && !tileMap.GetBit(row, col, FLOOR_BIT) {
// 		tileMap.SetBit(row, col, EXTERIOR_CORNER_CONCAVE_SW_BIT)
// 	}
// }

// func buildTerminusTiles(tileMap *TileMap, col, row int) {
// 	// left corner

// 	if tileMap.GetBit(row, col, HWALL_BIT) && (tileMap.GetBit(row, col, VWALL_BIT) || tileMap.GetBit(row+1, col, VWALL_BIT)) {
// 		setTerminus(tileMap, row, col-1)
// 	}

// 	// right corner

// 	if tileMap.GetBit(row, col, HWALL_BIT) && (tileMap.GetBit(row, col+1, VWALL_BIT) || tileMap.GetBit(row+1, col+1, VWALL_BIT)) {
// 		setTerminus(tileMap, row, col)
// 	}
// }

// func setTerminus(tileMap *TileMap, row, col int) {
// 	tileMap.SetBit(row, col, TERMINUS_NW_BIT)
// 	tileMap.SetBit(row+1, col, TERMINUS_SW_BIT)
// 	tileMap.SetBit(row, col+1, TERMINUS_NE_BIT)
// 	tileMap.SetBit(row+1, col+1, TERMINUS_SE_BIT)
// }

// func buildDoorTiles(tileMap *TileMap, row, col int) {
// 	// horizontal

// 	if tileMap.GetBit(row, col, HDOOR_BIT) {
// 		tileMap.SetBit(row, col, DOOR_N_BIT)
// 		tileMap.SetBit(row+1, col, DOOR_S_BIT)
// 	}

// 	// vertical

// 	if tileMap.GetBit(row, col, VDOOR_BIT) {
// 		tileMap.SetBit(row, col, DOOR_E_BIT)
// 		tileMap.SetBit(row, col-1, DOOR_W_BIT)
// 	}
// }
