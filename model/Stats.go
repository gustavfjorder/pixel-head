package model

type Stats struct {
	Health int
	Speed float32
}

const (
	Human  = iota
	zombie
)

func NewStats(being int) (s Stats) {
	switch being {
	case Human:
		s = Stats{
			Health:100,
			Speed:2,
		}
	case zombie:
		s = Stats{
			Health:20,
			Speed:1,
		}
	}
	return
}