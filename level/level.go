package level

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/node"
)

type Level struct {
	Bricks       [config.BrickRowCount][config.BrickColumnCount]*entities.Brick
	TotalHealth  int
	Paddle       *entities.Paddle
	Balls        [config.BallMaxCount]*entities.Ball
	BallCount    int
	CurrLevelNum int
	// TODO: Paddle here
	// TODO: Balls here
}

// Detects a general collision between two Rects
func isColliding(a *entities.Rect, b *entities.Rect) bool {
	// x axis
	if !(a.X+a.W < b.X || b.X+b.W < a.X) {

		// y axis
		if !(a.Y+a.H < b.Y || b.Y+b.H < a.Y) {
			return true
		}
	}
	return false
}

func (l *Level) LoadFromFile(path string) {
	//l.Bricks = nil
	// open the level file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	iRow := 0
	for scanner.Scan() {
		if !(iRow < config.BrickRowCount) {
			break
		}
		line := scanner.Text()
		fmt.Println(line, len(line), utf8.RuneCountInString(line))
		for iColumn, runeCharacter := range []rune(line) {
			if !(iColumn < config.BrickColumnCount) {
				break // max column count reached
			}
			switch runeCharacter {
			case '0':
				// no brick
				l.Bricks[iRow][iColumn] = nil
				fmt.Println("No brick")
			case '1':
				// basic brick

				fmt.Println("Basic brick")
				brick := entities.Brick{}
				brick.Health = 1
				brick.BrickType = 1
				brick.Rect.X = float64(iColumn * config.BrickWidth)
				brick.Rect.Y = float64(iRow * config.BrickHeight)
				brick.Rect.W = config.BrickWidth
				brick.Rect.H = config.BrickHeight
				brick.Image = ebiten.NewImage(
					int(brick.Rect.W),
					int(brick.Rect.H))
				brick.Image.Fill(color.RGBA{
					R: uint8(64),
					G: uint8(255),
					B: uint8(64),
					A: uint8(255)})
				l.TotalHealth += brick.Health
				l.Bricks[iRow][iColumn] = &brick

				// default to no brick
			default:
				// default to no brick
				fmt.Println("No brick")
				l.Bricks[iRow][iColumn] = nil
			}
		}
		iRow++
	}
}

// returns a pointer to a new level, based on the level number.
func NewLevel(levelNumber int) *Level {
	level := Level{}
	levelPath := fmt.Sprintf("levels/level%d.txt", levelNumber)
	fmt.Println("Loading level from file:", levelPath)
	level.LoadFromFile(levelPath)
	return &level
}

func (l *Level) PrintLevel() {
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			var brickType int
			if l.Bricks[iRow][iColumn] == nil {
				brickType = 0
			} else {
				brickType = l.Bricks[iRow][iColumn].BrickType
			}
			fmt.Printf("%d", brickType)
		}
		fmt.Println()
	}
}

func (l *Level) Update() {
	// TODO: Update every node in in the level
	fmt.Println("Update level")
	for i := range len(l.Balls) {
		if l.Balls[i] == nil {
			continue
		}
		alreadyBouncedBrick := false // prevents bounce cancellation if multiple collision
		for iRow := 0; iRow < config.BrickRowCount; iRow++ {
			for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
				if l.Bricks[iRow][iColumn] == nil {
					continue
				}
				if isColliding(&l.Balls[i].Rect, &l.Bricks[iRow][iColumn].Rect) {
					collidedBrick := l.Bricks[iRow][iColumn]
					// bounce if not already bounced (prevents bounce cancellation)
					if !alreadyBouncedBrick {
						// calculate collision lengts of x and y,
						// this determines if the collision is x or y sided
						// x
						var xCollisionLength, yCollisionLength float64
						if l.Balls[i].Rect.X < collidedBrick.Rect.X {
							xCollisionLength = l.Balls[i].Rect.X + l.Balls[i].Rect.W - collidedBrick.Rect.X
						} else {
							xCollisionLength = collidedBrick.Rect.X + collidedBrick.Rect.X - l.Balls[i].Rect.X
						}
						// y
						if l.Balls[i].Rect.Y < collidedBrick.Rect.Y {
							yCollisionLength = l.Balls[i].Rect.Y + l.Balls[i].Rect.H - collidedBrick.Rect.Y
						} else {
							yCollisionLength = collidedBrick.Rect.Y + collidedBrick.Rect.H - l.Balls[i].Rect.Y
						}

						if xCollisionLength >= yCollisionLength {
							// y-sided collision
							l.Balls[i].SpeedY = -l.Balls[i].SpeedY
							alreadyBouncedBrick = true
						} else {
							// x-sided collision
							l.Balls[i].SpeedX = -l.Balls[i].SpeedX
							alreadyBouncedBrick = true
						}

					}
					collidedBrick.Health--
					l.TotalHealth--
				}
			}
		}
	}

	// destroy bricks with 0 or less health
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			if l.Bricks[iRow][iColumn] == nil {
				continue
			}
			if l.Bricks[iRow][iColumn].Health <= 0 {
				l.Bricks[iRow][iColumn] = nil
			}
		}
	}

	// wall collisions & bounce
	for i := 0; i < len(l.Balls); i++ {
		if l.Balls[i] == nil {
			continue
		}
		// left wall
		if l.Balls[i].Rect.X <= 0 && l.Balls[i].SpeedX < 0 {
			l.Balls[i].SpeedX = -l.Balls[i].SpeedX
		}
		// right wall
		if l.Balls[i].Rect.X+l.Balls[i].Rect.W >= config.PlayAreaWidth &&
			l.Balls[i].SpeedX > 0 {
			l.Balls[i].SpeedX = -l.Balls[i].SpeedX
		}
		// ceiling
		if l.Balls[i].Rect.Y <= 0 && l.Balls[i].SpeedY < 0 {
			l.Balls[i].SpeedY = -l.Balls[i].SpeedY
		}
		// floor
		if l.Balls[i].Rect.Y+l.Balls[i].Rect.H >= config.PlayAreaHeight && l.Balls[i].SpeedY > 0 {
			// TODO: destroy ball
			if config.GodMode {
				// god mode just bounces off the floor too
				l.Balls[i].SpeedY = -l.Balls[i].SpeedY
			} else {
				l.Balls[i] = nil
				//fmt.Println("Ball destroyed")
				l.BallCount--
			}
		}
	}

	// Paddle collisions & bounce
	for i := 0; i < len(l.Balls); i++ {
		if l.Balls[i] == nil {
			continue
		}
		if !isColliding(&l.Balls[i].Rect, &l.Paddle.Rect) {
			continue
		}
		ballCenterX := l.Balls[i].Rect.X + l.Balls[i].Rect.W/2
		//fmt.Println("Ball centerX:", ballCenterX)
		segmentAngleDegrees := 22.5
		paddleSegmentLenX := l.Paddle.Rect.W / 6
		//fmt.Println("paddleSegmentLenx:", paddleSegmentLenX)

		if ballCenterX < l.Paddle.Rect.X+paddleSegmentLenX {
			//fmt.Println("multiball hit segment: 1")
			l.Balls[i].CalcXYForAngle(360.0 - segmentAngleDegrees*2 - segmentAngleDegrees/2.0)
		} else if ballCenterX < l.Paddle.Rect.X+paddleSegmentLenX*2 {
			//fmt.Println("multiball hit segment: 2")
			l.Balls[i].CalcXYForAngle(360.0 - segmentAngleDegrees - segmentAngleDegrees/2.0)
		} else if ballCenterX < l.Paddle.Rect.X+paddleSegmentLenX*3 {
			// fmt.Println("multiball hit segment: 3")
			l.Balls[i].CalcXYForAngle(360.0 - segmentAngleDegrees/2.0)
		} else if ballCenterX < l.Paddle.Rect.X+paddleSegmentLenX*4 {
			// fmt.Println("multiball hit segment: 4")
			l.Balls[i].CalcXYForAngle(segmentAngleDegrees / 2.0)
		} else if ballCenterX < l.Paddle.Rect.X+paddleSegmentLenX*5 {
			// fmt.Println("multiball hit segment: 5")
			l.Balls[i].CalcXYForAngle(segmentAngleDegrees + segmentAngleDegrees/2.0)
		} else {
			// fmt.Println("multiball hit segment: 6")
			l.Balls[i].CalcXYForAngle(segmentAngleDegrees*2 + segmentAngleDegrees/2.0)
		}
		// ensure that the ball bounces upwards
		if l.Balls[i].SpeedY > 0 {
			l.Balls[i].SpeedY = -l.Balls[i].SpeedY
		}

	}
	node.UpdateNode(l.Paddle)
	//g.paddle.Update()

	// update balls
	for i := 0; i < len(l.Balls); i++ {
		if l.Balls[i] == nil {
			continue
		}
		if l.Balls[i].Grabbed {
			l.Balls[i].Rect.X = l.Paddle.Rect.X
		} else {
			node.UpdateNode(l.Balls[i])
			//g.balls[i].Update()
		}
	}

	// level cleared, move to next
	if l.TotalHealth <= 0 {
		l.CurrLevelNum++
		// TODO: Change this to a scene system
		l = NewLevel(l.CurrLevelNum)
	}

	fmt.Println(l.Balls)
}

func (l *Level) Draw(screen *ebiten.Image) {
	// TODO: Draw every node in the level
	fmt.Println("Draw level")
}
