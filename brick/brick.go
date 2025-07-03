package brick

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
)

const (
	BrickTypeBasic = rune('1')
	BrickTypeSteel = rune('2')
)

var (
	BrickColorBasic = color.RGBA{
		R: uint8(0),
		G: uint8(255),
		B: uint8(0),
		A: uint8(255),
	}
	BrickColorSteel = color.RGBA{
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
	Image        *ebiten.Image
	Rect         entities.Rect
	Health       int
	Destructable bool
	BrickType    rune
	//Entity
}

// Hit reduces Bricks health accordily and returns damage dealt
func (b *Brick) Hit() int {
	if b == nil {
		return 0
	}
	switch b.BrickType {
	case BrickTypeBasic:
		b.Health--
		return 1 // return damage dealt
	case BrickTypeSteel:
		return 0 // steel brick is indestructible
	default:
		return 0
	}
}

// // Destroy sets *Brick to nil
// func destroy(b *Brick) *Brick{
// 	b = nil
// }

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

// NewBrick returns a pointer to a new Brick of known brickType and it's health.
// Returns nil if brickType not known.
// '1': Basic Brick
// '2': Steel Brick
func NewBrick(x, y float64, brickType rune) (*Brick, int) {
	brick := Brick{}
	brick.Rect.X = x
	brick.Rect.Y = y
	brick.Rect.W = config.BrickWidth
	brick.Rect.H = config.BrickHeight
	brick.BrickType = brickType
	brick.Image = ebiten.NewImage(int(brick.Rect.W), int(brick.Rect.H))
	switch brick.BrickType {
	case BrickTypeBasic:
		brick.Destructable = true
		brick.Image.Fill(BrickColorBasic)
		brick.Health = 1
	case BrickTypeSteel:
		//fmt.Println("Steel brick add")
		brick.Destructable = false
		brick.Image.Fill(BrickColorSteel)
		brick.Health = 0
	default:
		return nil, 0
	}
	//brick.Image.Fill(BrickColorBasic)
	return &brick, brick.Health
}

// BRICK BASIC ^
