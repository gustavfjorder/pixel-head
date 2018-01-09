package model

type Level struct {
	NumberOfZombies int
}

var Levels = map[int]Level {
	0: {100},
	1: {200},
	2: {300},
	3: {400},
	4: {500},
}
