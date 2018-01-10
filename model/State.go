package model

type State struct {
	Timestamp int64
	Players   []Player
	Zombies   []Zombie
	Shoots    []Shoot
}