package maps

type Base struct {
	Rooms           []Room
	NavigationGraph *NavigationGraph
}

func ConstructBase(tileMap *TileMap) *Base {
	rooms, walls, doors := partitionRooms(tileMap)
	navigationGraph := constructNavigationGraph(rooms, walls, doors)

	return &Base{
		Rooms:           rooms,
		NavigationGraph: navigationGraph,
	}
}
