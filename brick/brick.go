package brick

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
)

const (
	BrickTypeBasic = 1
)

var (
	BrickColorBasic = color.RGBA{
		R: uint8(0),
		G: uint8(255),
		B: uint8(0),
		A: uint8(255),
	}
	ColorBrickSteel = color.RGBA{
		R: uint8(127),
		G: uint8(127),
		B: uint8(127),
		A: uint8(255),
	}
)

// type Brick interface {
// 	Hit() int
// 	Destroy()
// }

// base struct for Bricks
type Brick struct {
	Image     *ebiten.Image
	Rect      entities.Rect
	Health    int
	BrickType int
	//Entity
}

func (b *Brick) Hit() int {
	if b == nil {
		return 0
	}
	switch b.BrickType {
	case BrickTypeBasic:
		b.Health--
		if b.Health <= 0 {
			b.Destroy()
		}
		return 1 // return damage dealt
	default:
		return 0
	}
}

func (b *Brick) Destroy() {
	b = nil
}

func (b *Brick) Update() {
	if b == nil {
		return
	}
	// TODO: collision with Ball (destruction)
}

func (b *Brick) Draw(screen *ebiten.Image) {
	if b == nil {
		return
	}
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(b.Rect.X), float64(b.Rect.Y)) // ??
	screen.DrawImage(b.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// NewBrickBasic returns a pointer to a new BrickBasic.
func NewBrick(x, y float64, brickType int) *Brick {
	brick := Brick{}
	brick.Rect.X = x
	brick.Rect.Y = y
	brick.Rect.W = config.BrickWidth
	brick.Rect.H = config.BrickHeight
	brick.BrickType = brickType
	brick.Image = ebiten.NewImage(int(brick.Rect.W), int(brick.Rect.H))
	switch brick.BrickType {
	case BrickTypeBasic:
		brick.Image.Fill(BrickColorBasic)
		brick.Health = 1
	default:
		return nil
	}
	brick.Image.Fill(BrickColorBasic)
	return &brick
}

// BRICK BASIC ^
