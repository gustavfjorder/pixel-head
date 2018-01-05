package model

type Stats struct {
	Health    int
	MoveSpeed float64
	Power     int
}

const (
	Human  = iota
	zombie
)

func NewStats(being int) (s Stats) {
	switch being {
	case Human:
		s = Stats{
			Health:    100,
			MoveSpeed: 5,
			Power:     5,
		}
	case zombie:
		s = Stats{
			Health:    20,
			MoveSpeed: 1,
			Power:     3,
		}
	}
	return
}