package model

import (
	"time"
  )

type State struct {
	Players   []Player
	Zombies   []ZombieI
	Shots     []Shot
	Barrels   []Barrel
	Lootboxes []Lootbox
}

var Timestamp time.Duration

type Entry struct {
	elem  EntityI
	index int
}

func (state State) Compress() (compressed State) {
	compressed.Players = state.Players
	compressed.Zombies = state.Zombies
	return
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

func (state *State) HasLootboxAt(point Point) bool {
	result := false

	for _, box := range state.Lootboxes {
		if box.Pos.X == point.X && box.Pos.Y == point.Y {
			result = true
			break
		}
	}

	return result
}
