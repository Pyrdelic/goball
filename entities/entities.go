package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	X, Y, W, H int
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
	p.Rect.X, _ = ebiten.CursorPosition()
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(p.Rect.X), float64(p.Rect.Y)) // ??
	screen.DrawImage(p.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// PADDLE ^^^^^

// BRICK vvvvvv
type Brick struct {
	Image  *ebiten.Image
	Rect   Rect
	Health int
	Updater
	Drawer
	//Entity
}

func (b *Brick) Update() {
	// TODO: collision with Ball (destruction)
}

func (b *Brick) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(b.Rect.X), float64(b.Rect.Y)) // ??
	screen.DrawImage(b.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// BRICK ^^^^^^

// BALL vvvvvv
type Ball struct {
	Image *ebiten.Image
	Rect  Rect
	Speed int
	// TODO: direction
	Updater
	Drawer
	//Entity
}

func (b *Ball) Update() {
	b.Rect.Y = b.Rect.Y + b.Speed
}

func (b *Ball) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(b.Rect.X), float64(b.Rect.Y)) // ??
	screen.DrawImage(b.Image, &DIO)
}

// BALL ^^^^^^
