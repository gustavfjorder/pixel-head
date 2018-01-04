package model

// todo: add ammunition handling
type Weapon struct {
	Id       int
	Name     string
	Power    float64
	Range    float64
	Speed    float64
	Magazine int
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
	},
	Knife:
	{
		Id:       Knife,
		Name:     "knife",
		Power:    20,
		Range:    20,
		Speed:    4,
		Magazine: 0,
	},
	Shotgun:
	{
		Id:       Shotgun,
		Name:     "shotgun",
		Power:    20,
		Range:    1000,
		Speed:    4,
		Magazine: 5,
	},
	Handgun:
	{
		Id:       Handgun,
		Name:     "handgun",
		Power:    20,
		Range:    1000,
		Speed:    4,
		Magazine: 10,
	},
}
