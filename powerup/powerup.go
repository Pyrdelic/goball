package powerup

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/node"
)

const (
	MultiBall int = iota
	Death
	ExtraLife
	ExpandPad
)

type PowerUp struct {
	Image          *ebiten.Image
	Rect           entities.Rect
	SpeedX, SpeedY float64
	PowerUpType    int
}

func (pu *PowerUp) Update() node.Message {
	pu.SpeedY += config.PowerUpGravity
	pu.Rect.X += pu.SpeedX
	pu.Rect.Y += pu.SpeedY
	return node.Message{TypeStr: "PowerUp"}
}

func (pu *PowerUp) Draw(screen *ebiten.Image) {
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(pu.Rect.X), float64(pu.Rect.Y))
	screen.DrawImage(pu.Image, &DIO)
}

func NewPowerUp(x, y, powerUpType int, speedX, speedY float64) *PowerUp {
	powerUp := PowerUp{}
	powerUp.Rect.X = float64(x)
	powerUp.Rect.Y = float64(y)
	powerUp.SpeedX, powerUp.SpeedY = speedX, speedY
	powerUp.PowerUpType = powerUpType
	powerUp.Image = ebiten.NewImage(config.PowerUpWidth, config.PowerUpHeight)
	var fillColor color.Color
	switch powerUp.PowerUpType {
	case MultiBall:
		fillColor = color.RGBA{
			R: uint8(127),
			G: uint8(127),
			B: uint8(127),
			A: uint8(255),
		}
	// case Death:
	// case ExtraLife:
	// case ExpandPad:
	default:
		fillColor = color.White
	}
	powerUp.Image.Fill(fillColor)
	return &powerUp
}
