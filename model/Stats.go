package model

type Stats struct {
	Health    int
	Being	  int
}

const (
	HUMAN  = iota
	ZOMBIE
)

func NewStats(being int) (s Stats) {
	s.Being=being
	s.Health = s.GetMaxHealth()
	return
}

func (s Stats) GetMaxHealth() int{
	switch s.Being {
	case HUMAN:
		return 100
	case ZOMBIE:
		return 20
	}
	return 0
}

func (s Stats) GetMoveSpeed() float64{
	switch s.Being {
	case HUMAN:
		return 10
	case ZOMBIE:
		return 2
	}
	return 0
}

func (s Stats) GetPower() int {
	switch s.Being {
	case HUMAN:
		return 5
	case ZOMBIE:
		return 3
	}
	return 0
}