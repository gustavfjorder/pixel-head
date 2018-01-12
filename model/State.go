package model

import (
	"time"
)

type State struct {
	Players   []Player
	Zombies   []Zombie
	Shots     []Shot
	Barrels   []Barrel
}

var Timestamp time.Duration

type Entry struct {
	elem  interface{}
	index int
}
