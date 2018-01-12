package model

type Updates struct {
	Removed []Entity
}

type EntityType int

const (
	ShotE EntityType = iota
	BarrelE
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
}

func (updates *Updates) Clear() {
	updates.Removed = make([]Entity, 0)
}

func (updates *Updates) Remove(entities ...EntityI) {
	for _, entity  := range entities {
		updates.Removed = append(updates.Removed, Entity{entity.ID(), entity.EntityType()})
	}
}

func (updates Updates) Empty() bool {
	return len(updates.Removed) <= 0
}
//
//
//func (updates Updates) Size() int{
//	return len(updates.Zombies) + len(updates.Players) + len(updates.Barrels) + len(updates.Shots)
//}
//
//func (updates *Updates) Add(elems ...interface{}) {
//	for _, elem := range elems {
//		switch elem.(type) {
//		case Zombie:
//			updates.Zombies = append(updates.Zombies, elem.(Zombie))
//		case Player:
//			updates.Players = append(updates.Players, elem.(Player))
//		case Shot:
//			updates.Shots = append(updates.Shots, elem.(Shot))
//		case Barrel:
//			updates.Barrels = append(updates.Barrels, elem.(Barrel))
//		default:
//			fmt.Fprint(os.Stderr, "Invalid interface for add got:", elem)
//		}
//	}
//}
