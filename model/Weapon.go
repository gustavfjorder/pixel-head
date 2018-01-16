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

type WeaponI interface {
	RefillMag() bool
	ReloadSpeed() time.Duration
	ShootDelay() time.Duration
	Power() int
	Range() float64
	ProjectileSpeed() float64
	BulletsPerShot() int
	Spread() float64
	Capacity() int
	Name() string
	Type() WeaponType
	AddBullets(n int)
	RemoveBullets(n int)
	GetMagazine() int
	GetBullets() int
	IncLevel()
	Shoot(Player, *Game)
}

type WeaponBase struct {
	WeaponI
	WeaponType
	MagazineCurrent int
	Bullets         int
	Level           int
}

func NewWeapon(weaponNum WeaponType) (weapon WeaponI) {
	if weaponNum >= nWeapon {
		panic(errors.New("No such weapon"))
	}
	base := WeaponBase{
		WeaponType:      weaponNum,
	}
	base.MagazineCurrent = base.MagazineCapacity()
	base.Bullets = base.Capacity()
	switch weaponNum {
	case SHOTGUN:
		weapon = &Shotgun{base}
	case RIFLE:
		weapon = &Rifle{base}
	case HANDGUN:
		weapon = &Handgun{base}
	default:
		weapon = &Knife{base}
	}
	return weapon
}

//Returns true if the magazine was reloaded
func (weapon *WeaponBase) RefillMag() bool {
	if weapon.MagazineCurrent >= weapon.MagazineCapacity() {
		return false
	}
	dBullet := MinInt(weapon.MagazineCapacity(), weapon.Bullets)
	weapon.MagazineCurrent = dBullet
	weapon.Bullets -= dBullet
	return dBullet > 0
}

func (weapon *WeaponBase) Shoot(player Player, game *Game){
	if weapon.GetMagazine() <= 0 {
		return
	}
	game.Add(NewShot(player))
	weapon.MagazineCurrent -= 1
}

func (weapon *Shotgun) Shoot(player Player, game *Game){
	if weapon.GetMagazine() <= 0 {
		return
	}
	offset := 0.0
	if weapon.BulletsPerShot()%2 == 0 {
		offset = weapon.Spread() * 0.5
	}
	angle := offset - weapon.Spread()*float64(weapon.BulletsPerShot()/2)
	shots := make([]Shot, weapon.BulletsPerShot())

	for i := 0; i < weapon.BulletsPerShot(); i++ {
		shots[i] = NewShot(player, angle)
		angle += weapon.Spread()
	}
	weapon.WeaponBase.MagazineCurrent -= 1
	for _, shot := range shots {
		game.Add(shot)
	}
}

func (knife *Knife) Shoot(player Player, game *Game) {
	for _, zombie := range game.State.Zombies{
		dAngle := math.Mod(math.Abs(player.Dir - angle(player.Pos, zombie.GetPos())),math.Pi)
		if dAngle <= knife.AttackCone() &&
			player.Pos.Sub(zombie.GetPos()).Len() <= knife.Range(){
				zombie.SubHealth(knife.Power())
		}
	}
}

func (weapon WeaponBase) GetMagazine() int {
	return weapon.MagazineCurrent
}

func (weapon WeaponBase) GetBullets() int {
	return weapon.Bullets
}

func (weapon *WeaponBase) IncLevel() {
	weapon.Level++
}

func (weapon *WeaponBase) AddBullets(n int) {
	weapon.Bullets += n
}

func (weapon WeaponBase) Type() WeaponType {
	return weapon.WeaponType
}

func (weapon WeaponBase) ReloadSpeed() time.Duration {
	return time.Second / 2
}

func (weapon WeaponBase) ShootDelay() time.Duration {
	return time.Second / 4
}

func (s Shotgun) ShootDelay() time.Duration {
	return time.Second / 4
}

func (r Rifle) ShootDelay() time.Duration {
	return time.Second / time.Duration(r.Level*4 + 1)
}

func (h Handgun) ShootDelay() time.Duration {
	return time.Second / 2
}

func (weapon WeaponBase) MagazineCapacity() int {
	switch weapon.Type() {
	case RIFLE:
		return 100
	case HANDGUN:
		return 10
	case SHOTGUN:
		return 10
	default:
		return 0
	}
}

func (weapon WeaponBase) Power() int {
	switch weapon.Type() {
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

func (weapon WeaponBase) Range() float64 {
	switch weapon.Type() {
	case RIFLE:
		return 500
	case HANDGUN:
		return 500
	case SHOTGUN:
		return 300
	case KNIFE:
		return 70
	default:
		return 0

	}
}

//Units per second
func (weapon WeaponBase) ProjectileSpeed() (speed float64) {
	switch weapon.Type() {
	default:
		speed = 1000
	}

	return
}

func (weapon WeaponBase) BulletsPerShot() int {
	return 1
}

func (s Shotgun) BulletsPerShot() int {
	return 5 + s.Level
}

func (weapon WeaponBase) Spread() float64 {
	switch weapon.Type() {
	case SHOTGUN:
		return math.Pi / 40
	default:
		return 0

	}
}

func (weapon WeaponBase) Capacity() int {
	switch weapon.Type() {
	case SHOTGUN:
		return 50
	case RIFLE:
		return 500
	default:
		return 50
	}
}

func (weapon *WeaponBase) RemoveBullets(n int) {
	weapon.MagazineCurrent -= n
}

func (weapon WeaponType) Name() string {
	switch weapon {
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
func (weapon WeaponBase) Name() string {
	return weapon.Type().Name()
}

type Shotgun struct {
	WeaponBase
}

type Rifle struct {
	WeaponBase
}

type Handgun struct {
	WeaponBase
}

type Knife struct {
	WeaponBase
}

func (knife Knife) AttackCone() float64 {
	return math.Pi * 2 /3
}