package model

import (
	"time"
	"math"
	"github.com/pkg/errors"
)

type WeaponType int

const (
	KNIFE   WeaponType = iota
	RIFLE
	SHOTGUN
	HANDGUN
	nWeapon
)

type Weapon struct {
	WeaponType
	MagazineCurrent int
	Bullets         int
}

func NewWeapon(weaponNum WeaponType) (weapon Weapon) {
	if weaponNum >= nWeapon {
		panic(errors.New("No such weapon"))
	}
	weapon.WeaponType = weaponNum
	weapon.MagazineCurrent = weapon.MagazineCapacity()
	weapon.Bullets = weapon.Capacity()
	return weapon
}

//Returns true if the magazine was reloaded
func (weapon *Weapon) RefillMag() bool {
	if weapon.MagazineCurrent >= weapon.MagazineCapacity() {
		return false
	}
	dBullet := minInt(weapon.MagazineCapacity(), weapon.Bullets)
	weapon.MagazineCurrent = dBullet
	weapon.Bullets -= dBullet
	return dBullet > 0
}

func (weapon *Weapon) GenerateShoots(player Player) []Shot {
	if weapon.MagazineCurrent <= 0 {
		return []Shot{}
	}
	offset := 0.0
	if weapon.BulletsPerShot()%2 == 0 {
		offset = weapon.Spread() * 0.5
	}
	angle := offset - weapon.Spread()*float64(weapon.BulletsPerShot()/2)
	shoots := make([]Shot, weapon.BulletsPerShot())

	for i := 0; i < weapon.BulletsPerShot(); i++ {
		shoots[i] = NewShot(player, angle)
		angle += weapon.Spread()
	}
	weapon.MagazineCurrent--
	return shoots
}

func (weaponType WeaponType) ReloadSpeed() time.Duration {
	return time.Second / 5
}

func (weaponType WeaponType) ShootDelay() time.Duration {
	switch weaponType {
	case RIFLE:
		return time.Second / 5
	case HANDGUN:
		return time.Second / 2
	case SHOTGUN:
		return time.Second / 4
	default:
		return 0
	}
}

func (weaponType WeaponType) MagazineCapacity() int {
	switch weaponType {
	case RIFLE:
		return 30
	case HANDGUN:
		return 10
	case SHOTGUN:
		return math.MaxInt32
	default:
		return 0
	}
}

func (weaponType WeaponType) Power() int {
	switch weaponType {
	case RIFLE:
		return 10
	case HANDGUN:
		return 10
	case SHOTGUN:
		return 5
	case KNIFE:
		return 100
	default:
		return 0
	}
}

func (weaponType WeaponType) Range() float64 {
	switch weaponType {
	case RIFLE:
		return 500
	case HANDGUN:
		return 500
	case SHOTGUN:
		return 300
	case KNIFE:
		return 100
	default:
		return 0

	}
}

//Units per second
func (weaponType WeaponType) ProjectileSpeed() (speed float64) {
	switch weaponType {
	default:
		speed = 1000
	}

	return
}

func (weaponType WeaponType) BulletsPerShot() int {
	switch weaponType {
	case SHOTGUN:
		return 5
	default:
		return 1
	}
}

func (weaponType WeaponType) Spread() float64 {
	switch weaponType {
	case SHOTGUN:
		return math.Pi / 40
	default:
		return 0

	}
}

func (weaponType WeaponType) Capacity() int {
	switch weaponType {
	case SHOTGUN:
		return 20
	default:
		return 150
	}
}

func (weaponType WeaponType) Name() string {
	switch weaponType {
	case RIFLE:
		return "rifle"
	case HANDGUN:
		return "handgun"
	case SHOTGUN:
		return "shotgun"
	case KNIFE:
		return "knife"
	default:
		return ""
	}
}
