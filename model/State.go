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
	elem  EntityI
	index int
}

type ByIndexDescending []Entry

func (s ByIndexDescending) Len() int {
	return len(s)
}

func (s ByIndexDescending) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByIndexDescending) Less(i, j int) bool {
	return s[i].index > s[j].index
}