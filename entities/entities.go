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
type Entity struct {
	X, Y, W, H int
}

// PADDLE vvvvv
type Paddle struct {
	Image *ebiten.Image
	Entity
	Updater
	Drawer
}

func (p *Paddle) Update() {
	p.X, _ = ebiten.CursorPosition()
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(p.X), float64(p.Y)) // ??
	screen.DrawImage(p.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// PADDLE ^^^^^

// BRICK vvvvvv
type Brick struct {
	Image *ebiten.Image
	Updater
	Drawer
	Entity
}

func (b *Brick) Update() {
	// TODO: collision with Ball
}

func (b *Brick) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(b.X), float64(b.Y)) // ??
	screen.DrawImage(b.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// BRICK ^^^^^^
