package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
)

type Updater interface {
	Update()
}

type Drawer interface {
	Draw()
}

// base struct for entities, with position and size
// TODO: rename to Rect?
type Entity struct {
	X, Y, W, H int
}

type Rect struct {
	X, Y, W, H float64
}

// PADDLE vvvvv
type Paddle struct {
	Image *ebiten.Image
	Rect  Rect
	//Entity
	Updater
	Drawer
}

func (p *Paddle) Update() {
	x, _ := ebiten.CursorPosition()
	// center paddle to cursor
	x = x - int(p.Rect.W/2)

	// constrain to walls
	if x < 0 {
		x = 0
	} else if x+int(p.Rect.W) > config.PlayAreaWidth {
		x = int(config.PlayAreaWidth - p.Rect.W)
	}
	p.Rect.X = float64(x)
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(p.Rect.X), float64(p.Rect.Y)) // ??
	screen.DrawImage(p.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// PADDLE ^^^^^

// BRICK vvvvvv

const (
	BrikcTypeNone           = 0
	BrickTypeBasic          = 1
	BrickTypeIndestructible = 2
)

type Brick struct {
	Image     *ebiten.Image
	Rect      Rect
	Health    int
	BrickType int
	Updater
	Drawer
	//Entity
}

func (b *Brick) Update() {
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

// // returns a pointer to a new brick.
// func NewBrick(brickType int, row int, column int) *Brick {
// 	brick := Brick{}
// 	brick.Rect = Rect{}
// 	return &brick
// }

// BRICK ^^^^^^

// BALL vvvvvv
type Ball struct {
	Image *ebiten.Image
	Rect  Rect
	//Speed          int
	SpeedMultiplier float64
	SpeedX, SpeedY  float64
	Grabbed         bool
	// TODO: direction
	Updater
	Drawer
	//Entity
}

func (b *Ball) Update() {
	//b.Rect.Y = b.Rect.Y + b.Speed
	if !b.Grabbed {
		b.Rect.Y = b.Rect.Y + b.SpeedY
		b.Rect.X = b.Rect.X + b.SpeedX
	}
}

func (b *Ball) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(b.Rect.X), float64(b.Rect.Y)) // ??
	screen.DrawImage(b.Image, &DIO)
}

// BALL ^^^^^^
