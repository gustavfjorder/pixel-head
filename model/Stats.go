package model

type Stats struct {
	Health    int
	MoveSpeed float64
	Power     int
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
	return
}