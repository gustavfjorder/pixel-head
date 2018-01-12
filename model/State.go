package model

import "time"

type State struct {
	Timestamp time.Duration
	Players   []Player
	Zombies   []Zombie
	Shots     []Shoot
	Barrels   []Barrel
}