package model

import "time"

type Level struct {
	NumberOfZombiesPerWave int
	NumberOfWaves          int
	TimeBetweenWaves       time.Duration
	TimeBetweenZombies	   time.Duration
}

var Levels = map[int]Level {
	0: {10000,1,time.Second*5,time.Second/100},
	1: {20,2,time.Second*5,time.Second/2},
	2: {30,3,time.Second*5,time.Second/3},
	3: {40,4,time.Second*5,time.Second/4},
	4: {50,5,time.Second*5,time.Second/5},
	5: {60,6,time.Second*5,time.Second/6},
	6: {70,7,time.Second*5,time.Second/7},
	7: {80,8,time.Second*5,time.Second/8},
	8: {90,9,time.Second*5,time.Second/9},
	9: {100,10,time.Second*5,time.Second/10},
	10: {110,11,time.Second*5,time.Second/11},

}

