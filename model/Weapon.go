package model

// todo: add ammunition handling
type Weapon struct {
	Id       int
	Name     string
	Power    float64
	Range    float64
	Speed    float64
	Magazine int
	Capacity int
	Bullets  int
}

const (
	Knife   = iota
	Rifle
	Shotgun
	Handgun
)

var Weapons = map[int]Weapon{
	Rifle:
	{
		Id:       Rifle,
		Name:     "rifle",
		Power:    20,
		Range:    1000,
		Speed:    4,
		Magazine: 30,
		Capacity: 150,
		Bullets: 0,
	},
	Knife:
	{
		Id:       Knife,
		Name:     "knife",
		Power:    20,
		Range:    20,
		Speed:    4,
		Magazine: -1,
		Capacity: -1,
		Bullets: -1,
	},
	Shotgun:
	{
		Id:       Shotgun,
		Name:     "shotgun",
		Power:    20,
		Range:    1000,
		Speed:    4,
		Magazine: 3,
		Capacity: 24,
		Bullets: 0,
	},
	Handgun:
	{
		Id:       Handgun,
		Name:     "handgun",
		Power:    20,
		Range:    1000,
		Speed:    4,
		Magazine: 10,
		Capacity: 50,
		Bullets: 0,
	},
}

func (weapon *Weapon) RefillMag(){
	weapon.Magazine=Weapons[weapon.Id].Magazine
}
