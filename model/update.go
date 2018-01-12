package model

type Updates State

func (updates Updates) Empty() bool {
	return !(len(updates.Zombies) > 0 || len(updates.Players) > 0 || len(updates.Barrels) > 0 || len(updates.Shots) > 0)
}

func (updates *Updates) Clear() {
	updates.Zombies = make([]Zombie, 0)
	updates.Players = make([]Player, 0)
	updates.Shots = make([]Shot, 0)
	updates.Barrels = make([]Barrel, 0)
}

func (updates Updates) Size() int{
	return len(updates.Zombies) + len(updates.Players) + len(updates.Barrels) + len(updates.Shots)
}

func (updates *Updates) Add(elems ...interface{}) {
	for _, elem := range elems {
		switch elem.(type) {
		case Zombie:
			updates.Zombies = append(updates.Zombies, elem.(Zombie))
		case Player:
			updates.Players = append(updates.Players, elem.(Player))
		case Shot:
			updates.Shots = append(updates.Shots, elem.(Shot))
		case Barrel:
			updates.Barrels = append(updates.Barrels, elem.(Barrel))
		}
	}
}
