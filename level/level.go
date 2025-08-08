package level

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pyrdelic/goball/brick"
	"github.com/pyrdelic/goball/config"
	"github.com/pyrdelic/goball/entities"
	"github.com/pyrdelic/goball/node"
	"github.com/pyrdelic/goball/powerup"
)

type Level struct {
	Bricks       [config.BrickRowCount][config.BrickColumnCount]*brick.Brick
	TotalHealth  int
	Paddle       *entities.Paddle
	Balls        []*entities.Ball
	BallCount    int
	CurrLevelNum int
	Lives        int
	PowerUps     []*powerup.PowerUp
}

// MESSAGES
const (
	GameOver int = iota + 1
	Pause
)

// var messageName = map[node.Message]string{
// 	GameOver: "Game over.",
// 	Pause:    "Pause",
// }

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
		//fmt.Println(line, len(line), utf8.RuneCountInString(line))
		for iColumn, runeCharacter := range []rune(line) {
			if !(iColumn < config.BrickColumnCount) {
				break // max column count reached
			}
			var healthToAdd int = 0
			l.Bricks[iRow][iColumn], healthToAdd = brick.NewBrick(
				float64(iColumn*config.BrickWidth),
				float64(iRow*config.BrickHeight),
				runeCharacter,
			)
			l.TotalHealth += healthToAdd
		}
		iRow++
	}
	fmt.Println("l.TotalHealth:", l.TotalHealth)
}

// returns a pointer to a new level, based on the level number.
func NewLevel(levelNumber int) *Level {
	level := Level{}
	level.BallCount = 1
	level.Lives = config.StartingLives
	level.CurrLevelNum = levelNumber
	levelPath := fmt.Sprintf("levels/level%d.txt", levelNumber)
	fmt.Println("Loading level from file:", levelPath)
	level.LoadFromFile(levelPath)
	level.Paddle = entities.NewPaddle()
	// TODO: addball function
	level.Balls = append(level.Balls, entities.NewBall(
		config.PlayAreaWidth/2.0,
		config.PlayAreaHeight/2.0,
		config.BallStartingSpeed,
		config.BallStartingAngle,
		true,
	))
	// level.Balls[0] = entities.NewBall(
	// 	config.PlayAreaWidth/2.0,
	// 	config.PlayAreaHeight/2.0,
	// 	config.BallStartingSpeed,
	// 	config.BallStartingAngle,
	// 	true,
	// )
	return &level
}

func (l *Level) PrintLevel() {
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			var brickType rune = rune('0')
			if l.Bricks[iRow][iColumn] == nil {
				brickType = rune('0')
			} else {
				brickType = l.Bricks[iRow][iColumn].BrickType
			}
			fmt.Print(string(brickType))
		}
		fmt.Println()
	}
}

// hitBrick damages brick(s) and returns the gained score
func (l *Level) hitBrick(brickRow, brickCol int) int {
	if l.Bricks[brickRow][brickCol] == nil || l.Bricks[brickRow][brickCol].Health <= 0 {
		return 0
	}

	switch l.Bricks[brickRow][brickCol].BrickType {
	case brick.BrickTypeBasic:

		l.Bricks[brickRow][brickCol].Health--
		l.TotalHealth--

		return config.BrickHitScore
	case brick.BrickTypeSteel:
		return 0
	case brick.BrickTypeFire:
		score := 0

		l.Bricks[brickRow][brickCol].Health--
		l.TotalHealth--

		score += config.BrickHitScore
		// one block wide cross-shaped area effect (explosion)
		if brickRow > 0 {
			score += l.hitBrick(brickRow-1, brickCol)
		}
		if brickRow < config.BrickRowCount-1 {
			score += l.hitBrick(brickRow+1, brickCol)
		}
		if brickCol > 0 {
			score += l.hitBrick(brickRow, brickCol-1)
		}
		if brickCol < config.BrickColumnCount-1 {
			score += l.hitBrick(brickRow, brickCol+1)
		}
		return score
	default:
		return 0
	}
	return 0
}

func (l *Level) Update() node.Message {
	if l == nil {
		return node.Message{
			TypeStr: "nil",
			Msg:     0,
		}
	}
	tickScore := 0
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return node.Message{
			TypeStr: "Level",
			Msg:     Pause,
		}
	}

	if l.BallCount <= 0 {
		// lose a life
		l.Lives--
		// Game over?
		if l.Lives < 0 {
			// TODO: messaging back to parent node
			return node.Message{
				TypeStr: "Level",
				Msg:     GameOver,
			}
		} else {
			l.Balls[0] = entities.NewBall(
				config.PlayAreaWidth/2.0,
				config.PlayAreaHeight/2.0,
				config.BallStartingSpeed,
				config.BallStartingAngle,
				true,
			)
			l.BallCount = 1
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for i := range len(l.Balls) {
			if l.Balls[i] == nil {
				continue
			}
			if l.Balls[i].Grabbed {
				l.Balls[i].Grabbed = false
				// make sure the ball launches upwards
				if l.Balls[i].SpeedY > 0 {
					l.Balls[i].SpeedY = -l.Balls[i].SpeedY
				}
			}
		}
	}

	//fmt.Println("Update level")

	var hitSpeedX, hitSpeedY float64
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
						hitSpeedX, hitSpeedY = l.Balls[i].SpeedX, l.Balls[i].SpeedY
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
					//collidedBrick.Health--
					// damage the brick

					// TODO: deprecate, move damage logic here
					// (area effects require access to all bricks)
					//scoreGained, damageDealt := collidedBrick.Hit()
					//damageDealt := 0

					scoreGained := l.hitBrick(iRow, iColumn)

					// switch collidedBrick.BrickType {
					// case brick.BrickTypeBasic:
					// 	collidedBrick.Health--
					// 	scoreGained += config.BrickHitScore
					// case brick.BrickTypeSteel:
					// 	break
					// case brick.BrickTypeFire:
					// 	collidedBrick.Health--
					// 	scoregained += config.BrickHitScore
					// }

					//l.TotalHealth -= damageDealt
					tickScore += scoreGained
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

			// destroy the brick
			if l.Bricks[iRow][iColumn].Health <= 0 &&
				l.Bricks[iRow][iColumn].Destructable {

				// spawn a powerup?
				fmt.Println("PowerUp spawned")
				l.PowerUps = append(l.PowerUps, powerup.NewPowerUp(
					int(l.Bricks[iRow][iColumn].Rect.X),
					int(l.Bricks[iRow][iColumn].Rect.Y),
					powerup.MultiBall,
					hitSpeedX*config.PowerUpSpeedMult,
					hitSpeedY*config.PowerUpSpeedMult,
				))

				// destroy the brick
				l.Bricks[iRow][iColumn] = nil
				tickScore += config.BrickDestroyedScore

			}
		}
	}

	// PowerUp hit detection
	for i := range len(l.PowerUps) {
		if l.PowerUps[i] == nil {
			continue
		}
		// left wall
		if l.PowerUps[i].Rect.X < 0 && l.PowerUps[i].SpeedX < 0 {
			l.PowerUps[i].SpeedX = -l.PowerUps[i].SpeedX
		}
		// right wall
		if l.PowerUps[i].Rect.X+l.PowerUps[i].Rect.W > config.PlayAreaWidth &&
			l.PowerUps[i].SpeedX > 0 {
			l.PowerUps[i].SpeedX = -l.PowerUps[i].SpeedX
		}

		// collect powerups with paddle
		if isColliding(&l.PowerUps[i].Rect, &l.Paddle.Rect) {
			switch l.PowerUps[i].PowerUpType {
			case powerup.MultiBall:
				// duplicate each ball, mirror SpeedX
				for iBall := range len(l.Balls) {
					if l.Balls[iBall] == nil {
						continue
					}
					newBall := entities.NewBall(
						l.Balls[iBall].Rect.X,
						l.Balls[iBall].Rect.Y,
						l.Balls[iBall].SpeedBase,
						config.BallStartingAngle,
						false,
					)
					fmt.Println("Got powerup!")

					newBall.SpeedX = -l.Balls[iBall].SpeedX
					newBall.SpeedY = l.Balls[iBall].SpeedY
					l.Balls = append(l.Balls, newBall)
					l.BallCount++
				}
			}
			l.PowerUps[i] = nil
		}
		fmt.Println(l.Balls)
		// floor
		if l.PowerUps[i] != nil &&
			l.PowerUps[i].Rect.Y > config.PlayAreaWidth {
			l.PowerUps[i] = nil
		}
	}
	//fmt.Println(l.PowerUps)
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
				fmt.Println("Ball destroyed") // TODO: not destroyed...
				fmt.Println(l.BallCount)
				l.BallCount--
				fmt.Println(l.BallCount)
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
		// speed up the ball
		l.Balls[i].SpeedBase *= config.BallSpeedIncrement

	}

	for i := range len(l.PowerUps) {
		if l.PowerUps[i] == nil {
			continue
		}
		node.Update(l.PowerUps[i])
	}

	node.Update(l.Paddle)
	//g.paddle.Update()

	// update balls
	for i := 0; i < len(l.Balls); i++ {
		if l.Balls[i] == nil {
			continue
		}
		if l.Balls[i].Grabbed {
			l.Balls[i].Rect.X = l.Paddle.Rect.X
		} else {
			node.Update(l.Balls[i])
			//g.balls[i].Update()
		}
	}

	// level cleared, move to next
	if l.TotalHealth <= 0 {

		// TODO: Change this to a scene system
		*l = *NewLevel(l.CurrLevelNum + 1)
	}

	//fmt.Println(l.Balls)
	return node.Message{
		TypeStr:  "Level",
		Msg:      0,
		IntExtra: tickScore,
	}
}

func (l *Level) Draw(screen *ebiten.Image) {
	if l == nil {
		return
	}
	// TODO: Draw every node in the level
	//fmt.Println("Draw level")

	// draw bricks
	for iRow := 0; iRow < config.BrickRowCount; iRow++ {
		for iColumn := 0; iColumn < config.BrickColumnCount; iColumn++ {
			//l.Bricks[iRow][iColumn].Draw(screen)
			node.Draw(l.Bricks[iRow][iColumn], screen)
		}
	}

	// draw powerups
	for i := range len(l.PowerUps) {
		if l.PowerUps[i] == nil {
			continue
		}
		node.Draw(l.PowerUps[i], screen)
	}
	// draw Paddle
	node.Draw(l.Paddle, screen)
	// draw Balls
	for i := range len(l.Balls) {
		node.Draw(l.Balls[i], screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"Ball count: %d",
		l.BallCount,
	))
}
