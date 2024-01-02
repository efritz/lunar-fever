package view

type View interface {
	Init()
	Exit()
	Update(elapsedMs int64, hasFocus bool)
	Render(elapsedMs int64)
	IsOverlay() bool
}

type ViewWrapper interface {
	InnerView() View
}

func areViewsEqual(view, target View) bool {
	if w, ok := view.(ViewWrapper); ok {
		view = w.InnerView()
	}

	if w, ok := target.(ViewWrapper); ok {
		target = w.InnerView()
	}

	return view == target
}
