package model

type Level struct {
	NumberOfZombies int
}

var Levels = map[int]Level {
	0: {10},
	1: {20},
	2: {30},
	3: {40},
	4: {50},
}
