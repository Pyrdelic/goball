package button

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/node"
)

const (
	Clicked int = iota + 1
)

type Button struct {
	entities.Rect
	Image          *ebiten.Image
	ClickedMessage int
}

// NewButton returns a pointer to a new Button.
func NewButton(x, y, w, h int) *Button {
	b := Button{}
	b.Rect.X = float64(x)
	b.Rect.Y = float64(y)
	b.Rect.W = float64(w)
	b.Rect.H = float64(h)
	b.ClickedMessage = Clicked
	b.Image = ebiten.NewImage(w, h)
	b.Image.Fill(color.White)
	return &b
}

// Clicked checks wheter x, y within Button Rect
func (b *Button) IsJustClicked() bool {
	if b == nil {
		return false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		x, y := ebiten.CursorPosition()
		if x >= int(b.Rect.X) && x <= int(b.Rect.X+b.Rect.W) &&
			y >= int(b.Rect.Y) && y <= int(b.Rect.Y+b.Rect.H) {
			return true
		}
	}
	return false
}

func (b *Button) Update() node.Message {
	if b.IsJustClicked() {
		fmt.Println("Button just clicked")
	}
	return node.Message{
		TypeStr: "Button",
		Msg:     Clicked,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	dio := ebiten.DrawImageOptions{}
	dio.GeoM.Translate(b.Rect.X, b.Rect.Y)
	screen.DrawImage(b.Image, &dio)
}
