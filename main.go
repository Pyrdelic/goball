package main

import (
	//"crypto/rand"

	"fmt"
	"image/color"
	"log"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pyrdelic/goball/entities"
)

// global variables
var (
	paddleHeight int = 600
)

type Game struct {
	// TODO: entities into their respective structs
	//paddle *ebiten.Image
	paddle       *ebiten.Image
	paddleStruct entities.Paddle
	bricks       []*ebiten.Image
	brickStructs [10]entities.Brick
}

func (g *Game) Update() error {
	g.paddleStruct.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	// paddleDIO := ebiten.DrawImageOptions{}
	// cursorX, _ := ebiten.CursorPosition()
	// paddleDIO.GeoM.Translate(float64(cursorX), 60)
	// ebiten.SetCursorMode(ebiten.CursorModeHidden)
	// screen.DrawImage(g.paddle, &paddleDIO)

	for i := range g.bricks {
		brickDIO := ebiten.DrawImageOptions{}
		//randomInt = rand.intN
		brickDIO.GeoM.Translate(float64(rand.Intn(640)), float64(rand.Intn(480)))
		screen.DrawImage(g.bricks[i], &brickDIO)
	}
	for i := range g.brickStructs {
		g.brickStructs[i].Draw(screen)
	}
	//fmt.Printf("before draw():\t%p\n", screen)
	g.paddleStruct.Draw(screen)
	// DIO := ebiten.DrawImageOptions{}
	// DIO.GeoM.Translate(float64(g.paddleStruct.X), float64(g.paddleStruct.Y))
	// screen.DrawImage(g.paddleStruct.Image, &DIO)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GO-BALL")

	game := Game{}
	game.paddleStruct.X, _ = ebiten.CursorPosition()
	game.paddleStruct.Y = 150
	game.paddleStruct.W = 80
	game.paddleStruct.H = 10
	game.paddleStruct.Image = ebiten.NewImage(game.paddleStruct.W, game.paddleStruct.H)
	//fmt.Println("dimensions: ", game.paddleStruct.W, game.paddleStruct.H) // works
	game.paddleStruct.Image.Fill(color.White)

	// init stuff before running the game
	//game.paddle = ebiten.NewImage(20, 20)
	//game.paddle.Fill(color.White)
	ebiten.SetVsyncEnabled(false)

	// init bricks
	for i := 0; i < 10; i++ {
		brick := ebiten.NewImage(10, 10)
		brick.Fill(color.RGBA{R: uint8(i * 5), G: 127, B: 127, A: 255})
		game.bricks = append(game.bricks, brick)
	}

	for i := range 10 {
		println("loop: ", i)
	}

	for i := 0; i < 10; i++ {
		game.brickStructs[i].X = 0 + i*10
		game.brickStructs[i].Y = 20
		game.brickStructs[i].W = 10
		game.brickStructs[i].H = 10
		game.brickStructs[i].Image = ebiten.NewImage(game.brickStructs[i].W, game.brickStructs[i].H)
		game.brickStructs[i].Image.Fill(color.RGBA{
			R: uint8(i * 5),
			G: uint8(127),
			B: uint8(127),
			A: uint8(255)})
	}

	fmt.Println("Run game")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
