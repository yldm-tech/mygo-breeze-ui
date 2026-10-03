package breeze

import "github.com/egoist/mygo/ui"

// StyleState applies a style function during this frame. It is kept in this
// package so components can express Tailwind-like hover, pressed, and
// disabled variants without a CSS selector engine.
func StyleState(e *ui.Element, fn func(*ui.Element)) *ui.Element {
	if fn != nil {
		fn(e)
	}
	return e
}
