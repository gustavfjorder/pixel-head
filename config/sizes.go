package config

const (
	//Human sprite settings
	HumanWidth = 40.0

	//Zombie sprite settings
	ZombieWidth = 30.0

	BulletWidth = 3.0

	//Don't touch below this point
	HumanPicWidth        = 148.0
	HumanHeightWidthFrac = 228.0 / 123.0
	HumanScale           = HumanWidth / HumanPicWidth
	HumaneHeight         = HumanWidth * HumanHeightWidthFrac

	ZombiePicWidth        = 31.0
	ZombieHeightWidthFrac = 72.0 / 31.0
	ZombieScale           = ZombieWidth / ZombiePicWidth
	ZombieHeight          = ZombieWidth * ZombieHeightWidthFrac

	BulletPicWidth        = 192.0
	BulletHeightWidthFrac = 511.0 / 192.0
	BulletScale           = BulletWidth / BulletPicWidth
	BulletHeight          = BulletWidth * BulletHeightWidthFrac

	//Made to fit with shotgun/rifle
	GunPosX = HumanWidth * 0.34
	GunPosY = HumaneHeight * 0.5
)
