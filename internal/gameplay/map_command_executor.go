package gameplay

type Palette int

const (
	FLOOR_TOOL Palette = iota
	HWALL_TOOL
	VWALL_TOOL
	HDOOR_TOOL
	VDOOR_TOOL
)

const MaxUndoStack = 100

var factories = map[Palette][2]MapCommandFactory{
	FLOOR_TOOL: {NewAddFloorMapCommandFactory(), NewRemoveFloorMapCommandFactory()},
	HWALL_TOOL: {NewAddHorizontalWallMapCommandFactory(), NewRemoveHorizontalWallMapCommandFactory()},
	VWALL_TOOL: {NewAddVerticalWallMapCommandFactory(), NewRemoveVerticalWallMapCommandFactory()},
	HDOOR_TOOL: {NewAddHorizontalDoorMapCommandFactory(), NewRemoveHorizontalDoorMapCommandFactory()},
	VDOOR_TOOL: {NewAddVerticalDoorMapCommandFactory(), NewRemoveVerticalDoorMapCommandFactory()},
}

type MapCommandExecutor struct {
	tileMap *TileMap
	undoLog []MapCommand
	redoLog []MapCommand
	factory MapCommandFactory
}

func NewMapCommandExecutor(tileMap *TileMap) *MapCommandExecutor {
	return &MapCommandExecutor{
		tileMap: tileMap,
	}
}

func (e *MapCommandExecutor) HasAction(tile Palette, row, col int) bool {
	return e.factoryFor(tile, row, col).Valid(e.tileMap, row, col)
}

func (e *MapCommandExecutor) PrepareAction(tile Palette, row, col int) {
	e.factory = e.factoryFor(tile, row, col)
}

func (e *MapCommandExecutor) ExecuteAction(tile Palette, row, col int) bool {
	if e.factory != nil && e.factory.Valid(e.tileMap, row, col) {
		command := e.factory.Create(e.tileMap, row, col)
		command.Execute()

		e.redoLog = nil
		e.undoLog = append(e.undoLog, command)
		if len(e.undoLog) > MaxUndoStack {
			e.undoLog = e.undoLog[1:]
		}

		return true
	}

	return false
}

func (e *MapCommandExecutor) Undo() bool {
	if len(e.undoLog) == 0 {
		return false
	}

	e.redoLog = append(e.redoLog, e.undoLog[len(e.undoLog)-1])
	if len(e.redoLog) > MaxUndoStack {
		e.redoLog = e.redoLog[1:]
	}

	e.undoLog[len(e.undoLog)-1].Unexecute()
	e.undoLog = e.undoLog[:len(e.undoLog)-1]
	return true
}

func (e *MapCommandExecutor) Redo() bool {
	if len(e.redoLog) == 0 {
		return false
	}

	e.undoLog = append(e.undoLog, e.redoLog[len(e.redoLog)-1])
	if len(e.undoLog) > MaxUndoStack {
		e.undoLog = e.undoLog[1:]
	}

	e.redoLog[len(e.redoLog)-1].Execute()
	e.redoLog = e.redoLog[:len(e.redoLog)-1]
	return true
}

func (e *MapCommandExecutor) ClearUndoState() {
	e.undoLog = nil
	e.redoLog = nil
}

func (e *MapCommandExecutor) factoryFor(tile Palette, row, col int) MapCommandFactory {
	factory1 := factories[tile][0]
	factory2 := factories[tile][1]

	if factory2.Valid(e.tileMap, row, col) {
		return factory2
	}

	return factory1
}
