package model

import (
	"time"
	"math"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/config"
)

type Weapon struct {
	Id               int
	MagazineCurrent  int
	Bullets          int
}

const (
	KNIFE   = iota
	RIFLE
	SHOTGUN
	HANDGUN
)

func NewWeapon(weaponNum int) (weapon Weapon){
	weapon.Id = weaponNum
	weapon.MagazineCurrent = weapon.GetMagazineCapacity()
	weapon.Bullets = weapon.GetCapacity()
	return weapon
}

func GetWeaponRef(weaponNum int) (Weapon) {
	return Weapon{Id:weaponNum}
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
			Start:     player.Pos.Add(pixel.V(config.GunPosX, config.GunPosY).Rotated(player.Dir - math.Pi/2)),
			Angle:     player.Dir + (float64(angle) * (math.Pi / 180)),
			StartTime: timestamp,
			Weapon:    weapon.Id,
		}

		angle += weapon.GetSpread()
		weapon.MagazineCurrent--
	}

	return shoots
}

func (weapon Weapon) GetReloadSpeed() time.Duration {
	switch weapon.Id {
	case RIFLE:
		return time.Second / 2
	case HANDGUN:
		return time.Second / 2
	case SHOTGUN:
		return time.Second / 2
	default:
		return time.Second / 2
	}
}

func (weapon Weapon) GetShootDelay() time.Duration {
	switch weapon.Id {
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
	switch weapon.Id {
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
	switch weapon.Id {
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
	switch weapon.Id {
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

func (weapon Weapon) GetProjectileSpeed() float64 {
	switch weapon.Id {
	default:
		return 1000
	}
}

func (weapon Weapon) GetBulletsPerShot() int {
	switch weapon.Id {
	case SHOTGUN:
		return 5
	default:
		return 1
	}
}

func (weapon Weapon) GetSpread() int {
	switch weapon.Id {
	case SHOTGUN:
		return 5
	default:
		return 0

	}
}

func (weapon Weapon) GetCapacity() int {
	switch weapon.Id {
	default:
		return 150
	}
}

func (weapon Weapon) GetName() string {
	switch weapon.Id {
	case RIFLE:
		return "rifle"
	case HANDGUN:
		return "handgun"
	case SHOTGUN:
		return "shotgun"
	case KNIFE:
		return "knife"
	default:
		return "none"
	}
}