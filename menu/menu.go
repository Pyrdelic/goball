package menu

import "github.com/pyrdelic/goball/entities"

type Button struct {
	entities.Rect
}

// Clicked checks wheter x, y within Button Rect
func (b *Button) Clicked(x, y int) bool {
	if b == nil {
		return false
	}
	if x >= int(b.X) && x <= int(b.W) &&
		y >= int(b.Y) && y <= int(b.H) {
		return true
	}
	return false
}
