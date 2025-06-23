package config

const (
	PlayAreaHeight      = 240 // in-game resolution
	PlayAreaWidth       = 320 // in-game resolution
	BrickColumnCount    = 16
	BrickRowCount       = 6
	BrickCount          = BrickColumnCount * BrickRowCount
	BrickHeight         = 10
	BrickWidth          = PlayAreaWidth / BrickColumnCount
	PaddleStartingWidth = PlayAreaWidth / 6
	BallStartingSpeed   = 2.0
)

type AppConfig struct {
	Test int
}

// var Config = AppConfig{
// 	Test: 1337,
// }
