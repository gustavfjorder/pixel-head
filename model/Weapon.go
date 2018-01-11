package model

import (
	"time"
	"math"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/pkg/errors"
)


type WeaponType int
const (
	KNIFE  WeaponType = iota
	RIFLE
	SHOTGUN
	HANDGUN
	nWeapon
)


type Weapon struct {
	weaponType      WeaponType
	MagazineCurrent int
	Bullets         int
}

func NewWeapon(weaponNum WeaponType) (weapon Weapon){
	if weaponNum >= nWeapon{
		panic(errors.New("No such weapon"))
	}
	weapon.weaponType = weaponNum
	weapon.MagazineCurrent = weapon.GetMagazineCapacity()
	weapon.Bullets = weapon.GetCapacity()
	return weapon
}

func GetWeaponRef(weaponType WeaponType) (Weapon) {
	return Weapon{weaponType: weaponType}
}

//Returns true if the magazine was reloaded
func (weapon *Weapon) RefillMag() bool {
	if weapon.MagazineCurrent >= weapon.GetMagazineCapacity(){
		return false
	}
	dBullet := minInt(weapon.GetMagazineCapacity(), weapon.Bullets)
	weapon.MagazineCurrent = dBullet
	weapon.Bullets -= dBullet
	return dBullet > 0
}

func (weapon *Weapon) GenerateShoots(timestamp time.Duration, player Player) []Shoot {
	shotsPerSideOfDirection := int(math.Floor(float64(weapon.GetBulletsPerShot() / 2)))
	angle := -(shotsPerSideOfDirection * weapon.GetBulletsPerShot())
	shoots := make([]Shoot, int(math.Min(float64(weapon.GetBulletsPerShot()), float64(weapon.MagazineCurrent))))

	for i := 0; i < weapon.GetBulletsPerShot() && weapon.MagazineCurrent > 0; i++ {
		shoots[i] = Shoot{
			Start:      player.Pos.Add(pixel.V(config.GunPosX, config.GunPosY).Rotated(player.Dir - math.Pi/2)),
			Angle:      player.Dir + (float64(angle) * (math.Pi / 180)),
			StartTime:  timestamp,
			WeaponType: weapon.weaponType,
		}

		angle += weapon.GetSpread()
		weapon.MagazineCurrent--
	}

	return shoots
}

func (weapon Weapon) GetReloadSpeed() time.Duration {
	return time.Second / 5
}

func (weapon Weapon) GetShootDelay() time.Duration {
	switch weapon.weaponType {
	case RIFLE:
		return time.Second / 5
	case HANDGUN:
		return time.Second / 2
	case SHOTGUN:
		return time.Second / 2
	default:
		return 0
	}
}

func (weapon Weapon) GetMagazineCapacity() int {
	switch weapon.weaponType {
	case RIFLE:
		return 30
	case HANDGUN:
		return 10
	case SHOTGUN:
		return 20
	default:
		return 0
	}
}

func (weapon Weapon) GetPower() int{
	switch weapon.weaponType {
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

func (weapon Weapon) GetRange() float64 {
	switch weapon.weaponType {
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
func (weapon Weapon) GetProjectileSpeed() (speed float64) {
	switch weapon.weaponType {
	default:
		speed = 1000
	}

	return
}

func (weapon Weapon) GetBulletsPerShot() int {
	switch weapon.weaponType {
	case SHOTGUN:
		return 5
	default:
		return 1
	}
}

func (weapon Weapon) GetSpread() int {
	switch weapon.weaponType {
	case SHOTGUN:
		return 5
	default:
		return 0

	}
}

func (weapon Weapon) GetCapacity() int {
	switch weapon.weaponType {
	default:
		return 150
	}
}

func (weaponType WeaponType) GetName() string {
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

func (weapon Weapon) WeaponType() WeaponType{
	return weapon.weaponType
}