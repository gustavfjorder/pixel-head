package model

type Level struct {
	NumberOfZombies int
}

var Levels = map[int]Level {
	0: {10},
	1: {15},
	2: {20},
	3: {25},
	4: {30},
}
