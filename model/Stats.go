package model

type Stats struct {
	Health    int
	MoveSpeed float64
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
		}
	case zombie:
		s = Stats{
			Health:    20,
			MoveSpeed: 1,
		}
	}
	return
}