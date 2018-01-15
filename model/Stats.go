package model

import (
	"github.com/gustavfjorder/pixel-head/config"
)

type Stats struct {
	Health    int
	Being
}

type Being int
const (
	HUMAN Being = iota
	ZOMBIE
	FASTZOMBIE
	SLOWZOMBIE
	BOMBZOMBIE
	nBeing
)

func NewStats(being Being) (s Stats) {
	if being >= nBeing{
		panic("Invalid being")
	}
	s.Being=being
	s.Health = s.GetMaxHealth()
	return
}

func (being Being) GetMaxHealth() int{
	switch being {
	case HUMAN:
		return 100
	case ZOMBIE:
		return 20
	case FASTZOMBIE:
		return 5
	case SLOWZOMBIE:
		return 100
	case BOMBZOMBIE:
		return 20
	}
	return 0
}

//Number of units per second
func (being Being) GetMoveSpeed() (speed float64){
	switch being {
	case HUMAN:
		speed = 400
	case ZOMBIE:
		speed =  100
	case FASTZOMBIE:
		speed = 200
	case SLOWZOMBIE:
		speed = 50
	case BOMBZOMBIE:
		speed = 200
	}
	return speed * config.Conf.ServerHandleSpeed.Seconds()
}

func (being Being) GetPower() int {
	switch being {
	case HUMAN:
		return 5
	case ZOMBIE:
		return 3
	case FASTZOMBIE:
		return 1
	case SLOWZOMBIE:
		return 15
	case BOMBZOMBIE:
		return 3
	}
	return 0
}