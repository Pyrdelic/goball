package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/node"
)

// type Updater interface {
// 	Update()
// }

// type Drawer interface {
// 	Draw()
// }

// type Node interface {
// 	Update()
// 	Draw(*ebiten.Image)
// }

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
}

// NewPaddle returns a pointer to a new Paddle
func NewPaddle() *Paddle {
	paddle := Paddle{}
	cursorX, _ := ebiten.CursorPosition()
	paddle.Rect.X = float64(cursorX)
	paddle.Rect.Y = 200
	paddle.Rect.W = config.PaddleStartingWidth
	paddle.Rect.H = 5
	paddle.Image = ebiten.NewImage(int(paddle.Rect.W), int(paddle.Rect.H))
	paddle.Image.Fill(color.White)
	return &paddle
}

func (p *Paddle) Update() node.Message {
	//fmt.Println("Paddle update")
	if p == nil {
		return 0
	}
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
	return 0
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	if p == nil {
		return
	}
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(p.Rect.X), float64(p.Rect.Y)) // ??
	screen.DrawImage(p.Image, &DIO)
	//fmt.Printf("%+v\n", p)
}

// PADDLE ^^^^^

// BRICK vvvvvv

// const (
// 	BrikcTypeNone           = 0
// 	BrickTypeBasic          = 1
// 	BrickTypeIndestructible = 2
// )

type Brick struct {
	Image     *ebiten.Image
	Rect      Rect
	Health    int
	BrickType int
	//Entity
}

func (b *Brick) Update() node.Message {
	if b == nil {
		return 0
	}
	// TODO: collision with Ball (destruction)
	return 0
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

// BRICK ^^^^^^

// BALL vvvvvv
type Ball struct {
	Image *ebiten.Image
	Rect  Rect
	//Speed          int
	SpeedBase      float64
	SpeedX, SpeedY float64
	Grabbed        bool
	// TODO: direction
	//Entity
}

func (b *Ball) Update() node.Message {
	if b == nil {
		return 0
	}
	//b.Rect.Y = b.Rect.Y + b.Speed
	if !b.Grabbed {
		b.Rect.Y = b.Rect.Y + b.SpeedY
		b.Rect.X = b.Rect.X + b.SpeedX
	}
	return 0
}

func (b *Ball) Draw(screen *ebiten.Image) {
	if b == nil {
		return
	}
	DIO := ebiten.DrawImageOptions{}
	DIO.GeoM.Translate(float64(b.Rect.X), float64(b.Rect.Y)) // ??
	screen.DrawImage(b.Image, &DIO)
}

// Sets a Ball SpeedX, SpeedY for a given angle (degrees).
func (b *Ball) CalcXYForAngle(angle float64) {
	if b == nil {
		return
	}
	radian := angle * (math.Pi / 180)
	b.SpeedX = b.SpeedBase * math.Sin(radian)
	// flip Y component to correct for game space coordinate system
	b.SpeedY = -(b.SpeedBase * math.Cos(radian))
}

// NewBall Returns a pointer to a new Ball.
func NewBall(x, y, speedBase, angle float64, grabbed bool) *Ball {
	ball := Ball{}
	ball.Grabbed = grabbed
	ball.SpeedBase = speedBase
	//ball.SpeedX, ball.SpeedY = speedXYForAngle(angle)
	rect := Rect{}
	rect.X = x
	rect.Y = y
	rect.W = config.BallSize
	rect.H = config.BallSize
	ball.Rect = rect
	ball.Image = ebiten.NewImage(int(ball.Rect.W), int(ball.Rect.H))
	ball.Image.Fill(color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
	ball.CalcXYForAngle(angle)
	return &ball
}

// BALL ^^^^^^
