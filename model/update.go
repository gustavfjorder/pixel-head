package model

import (
	"github.com/faiface/pixel"
	"time"
)

type Updates struct {
	Removed   []Entity
	Added     []EntityI
	Timestamp time.Duration
}

type EntityType int

const (
	ShotE    EntityType = iota
	BarrelE
	LootboxE
	ZombieE
	PlayerE
)

type Entity struct {
	ID string
	EntityType
}

type EntityI interface {
	ID() string
	EntityType() EntityType
	GetPos() pixel.Vec
	GetHitbox() float64
	GetDir() float64
}

func (updates *Updates) Clear() {
	updates.Removed = make([]Entity, 0)
	updates.Added = make([]EntityI, 0)
}

func (updates *Updates) Add(entities ...EntityI) {
	updates.Added = append(updates.Added, entities...)
}

func (updates *Updates) Remove(entities ...EntityI) {
	for _, entity := range entities {
		updates.Removed = append(updates.Removed, Entity{ID: entity.ID(), EntityType: entity.EntityType()})
	}
}

func (updates Updates) Empty() bool {
	return len(updates.Removed) <= 0 && len(updates.Added) <= 0
}
