package model

type Stats struct {
	Health    int
	MoveSpeed float64
	Power     int
	Being	  int
}

const (
	HUMAN  = iota
	ZOMBIE
)

func NewStats(being int) (s Stats) {
	switch being {
	case HUMAN:
		s = Stats{
			Health:    100,
			MoveSpeed: 10,
			Power:     5,
		}
	case ZOMBIE:
		s = Stats{
			Health:    20,
			MoveSpeed: 2,
			Power:     3,
		}
	}
	s.Being=being
	return
}

func (s Stats) GetMaxHealth() float64{
	switch s.Being {
	case HUMAN:
		return 100
	case ZOMBIE:
		return 20
	}
	return 0
}