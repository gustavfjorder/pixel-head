package model


type State struct {
	Timestamp int
	Players   []Player
	Zombies   []Zombie
	Shoots    []Shoot
}