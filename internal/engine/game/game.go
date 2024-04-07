package game

import "time"

type Game struct {
	slowUpdateThreshold int
	isStarted           bool
	isRunning           bool
	isRunningSlowly     bool
	isFixedTimeStep     bool
	targetElapsedTimeMs int64
	alwaysUpdate        bool
	alwaysRender        bool
	isRenderSuppressed  bool
	hasCalledUpdate     bool
	hasCalledRender     bool

	delegate     Delegate
	startTimeMs  int64
	lastUpdateNs int64
	lastRenderNs int64
	slowUpdates  int
}

func New(delegate Delegate) *Game {
	return &Game{
		slowUpdateThreshold: 20,
		isFixedTimeStep:     true,
		targetElapsedTimeMs: 1000 / 120,
		alwaysUpdate:        true,
		alwaysRender:        true,
		delegate:            delegate,
	}
}

func (g *Game) Start() {
	if g.isStarted {
		panic("illegal state")
	}

	g.run()
}

func (g *Game) Stop() {
	if !g.isStarted {
		panic("illegal state")
	}

	g.isRunning = false
}

func (g *Game) run() {
	g.isStarted = true
	g.isRunning = true
	g.startTimeMs = time.Now().UnixMilli()

	g.delegate.Init()

	for g.isRunning {
		g.tick()
	}

	g.delegate.Exit()
}

func (g *Game) tick() {
	// If the loop is running in fixed timestep and at least one update has already been
	// called, wait until targetElapsedTime has passed before calling the next update.

	if g.isFixedTimeStep && g.hasCalledUpdate {
		elapsed := g.getMsSinceLastUpdate()

		if elapsed < g.targetElapsedTimeMs {
			time.Sleep(time.Duration(g.targetElapsedTimeMs-elapsed) * time.Millisecond)
		}
	}

	// Call update.

	if g.delegate.Active() || g.alwaysUpdate {
		//
		// TODO - if inactive, should the elapsed time reset?

		elapsedMs := g.getMsSinceLastUpdate()
		g.lastUpdateNs = time.Now().UnixNano()

		g.delegate.Update(elapsedMs)
		g.hasCalledUpdate = true

		if !g.isRunning {
			return
		}
	}

	// Determine if the game loop is running slowly. Maintain a counter for the number of
	// consecutive updates that took longer than targetElapsedTime. If this value becomes
	// greater than slowUpdateThreshold, the loop is running consistently slower than the
	// target elapsed time.

	g.isRunningSlowly = false

	if g.isFixedTimeStep && g.getMsSinceLastUpdate() > g.targetElapsedTimeMs {
		g.slowUpdates++
		if g.slowUpdates >= g.slowUpdateThreshold {
			g.isRunningSlowly = true
		}
	} else {
		g.slowUpdates = 0
	}

	// Call render if the game loop is not running slowly and the last call to update did
	// not suppress rendering for this tick.

	if (g.delegate.Active() || g.alwaysRender) && !g.isRenderSuppressed && !g.isRunningSlowly {
		//
		// TODO - if inactive, suppressed, or skipped, should the elapsed time reset?

		elapsedMs := g.getMsSinceLastRender()
		g.lastRenderNs = time.Now().UnixNano()

		g.delegate.Render(elapsedMs)
		g.hasCalledRender = true
	}

	g.isRenderSuppressed = false
}

func (g *Game) getMsSinceLastUpdate() int64 {
	if !g.hasCalledUpdate {
		return 0
	}

	return (time.Now().UnixNano() - g.lastUpdateNs) / 1e6
}

func (g *Game) getMsSinceLastRender() int64 {
	if !g.hasCalledRender {
		return 0
	}

	return (time.Now().UnixNano() - g.lastRenderNs) / 1e6
}
