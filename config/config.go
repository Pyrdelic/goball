package config

const (
	PlayAreaWidth  = 320 // in-game resolution
	PlayAreaHeight = 240 // in-game resolution

	StartingLives = 3

	BrickColumnCount    = 16
	BrickRowCount       = 6
	BrickCount          = BrickColumnCount * BrickRowCount
	BrickHeight         = 10
	BrickWidth          = PlayAreaWidth / BrickColumnCount
	BrickHitScore       = 10
	BrickDestroyedScore = 100

	PaddleStartingWidth = PlayAreaWidth / 6

	BallStartingSpeed  = 2.0
	BallSize           = 5.0
	BallMaxCount       = 4
	BallStartingAngle  = 360.0 - 12.75
	BallSpeedIncrement = 1.10

	ButtonHeight = 50
	ButtonWidth  = 50

	MenuTitleTextX = PlayAreaWidth / 8
	MenuTitleTextY = PlayAreaHeight / 8

	HiScoreTopCount = 10
	HiScoreNameLen  = 3

	PowerUpGravity   = 0.15
	PowerUpWidth     = 20
	PowerUpHeight    = 20
	PowerUpSpeedMult = 2.0

	GodMode = false
)
