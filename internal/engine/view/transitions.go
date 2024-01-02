package view

import "github.com/efritz/lunar-fever/internal/common/math"

type TransitionableView interface {
	View
	BeginExiting()
}

type TransitionView struct {
	view               View
	viewManager        *Manager
	transitionOnTime   int64
	transitionOffTime  int64
	state              ViewState
	isExiting          bool
	transitionPosition float64
}

type ViewState int

const (
	ViewStateActive ViewState = iota
	ViewStateHidden
	ViewStateTransitionOn
	ViewStateTransitionOff
)

func NewTransitionView(view View, viewManager *Manager) TransitionableView {
	return &TransitionView{
		view:              view,
		viewManager:       viewManager,
		state:             ViewStateTransitionOn,
		transitionOnTime:  500,
		transitionOffTime: 500,
	}
}

func (v *TransitionView) InnerView() View {
	return v.view
}

func (v *TransitionView) BeginExiting() {
	v.isExiting = true
}

func (v *TransitionView) Init() {
	v.view.Init()
}

func (v *TransitionView) Exit() {
	v.view.Exit()
}

func (v *TransitionView) Update(elapsedMs int64, hasFocus bool) {
	if v.isExiting {
		v.state = ViewStateTransitionOff

		if v.transitionOff(elapsedMs) {
			v.viewManager.Remove(v)
			v.Exit()
		}
	} else {
		if !hasFocus {
			if v.transitionOff(elapsedMs) {
				v.state = ViewStateHidden
			} else {
				v.state = ViewStateTransitionOff
			}
		} else {
			if v.transitionOn(elapsedMs) {
				v.state = ViewStateActive
			} else {
				v.state = ViewStateTransitionOn
			}

			v.view.Update(elapsedMs, hasFocus)
		}
	}
}

func (v *TransitionView) transitionOn(elapsedMs int64) (finished bool) {
	delta := 1.0
	if v.transitionOnTime != 0 {
		delta = float64(elapsedMs) / float64(v.transitionOnTime)
	}

	return v.transition(delta)
}

func (v *TransitionView) transitionOff(elapsedMs int64) (finished bool) {
	delta := 1.0
	if v.transitionOffTime != 0 {
		delta = float64(elapsedMs) / float64(v.transitionOffTime)
	}

	return v.transition(-delta)
}

func (v *TransitionView) transition(delta float64) (finished bool) {
	v.transitionPosition, finished = math.Clamp(v.transitionPosition+delta, 0, 1)
	return finished
}

func (v *TransitionView) Render(elapsedMs int64) {
	v.view.Render(elapsedMs)
}

func (v *TransitionView) IsOverlay() bool {
	return v.view.IsOverlay()
}
