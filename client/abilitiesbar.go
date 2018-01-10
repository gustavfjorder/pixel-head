package client

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/config"
	"fmt"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

//ability sprites
var ab = Load(config.Conf.AbilityPath, "", IMG)

func DrawAbilities(win *pixelgl.Window, me *model.Player) {
	if me == nil {
		return
	}
	var (
		//abilities bar
		abilitiesBar = ab["abilitiesBar.png"].Sprites[0]

		//knife
		knife         = ab["knife.png"].Sprites[0]
		knifeSelected = ab["knifeSelected.png"].Sprites[0]

		//handgun
		handgun         = ab["handgun.png"].Sprites[0]
		handgunDark     = ab["handgunDark.png"].Sprites[0]
		handgunSelected = ab["handgunSelected.png"].Sprites[0]

		//rifle
		rifle         = ab["rifle.png"].Sprites[0]
		rifleDark     = ab["rifleDark.png"].Sprites[0]
		rifleSelected = ab["rifleSelected.png"].Sprites[0]

		//shotgun
		shotgun         = ab["shotgun.png"].Sprites[0]
		shotgunDark     = ab["shotgunDark.png"].Sprites[0]
		shotgunSelected = ab["shotgunSelected.png"].Sprites[0]

		//dimensions of abilities bar
		abilitiesBarPosX = win.Bounds().Max.X / 2
		abilitiesBarPosY = win.Bounds().Min.Y + abilitiesBar.Picture().Bounds().Max.Y/2

		//the factor by which we scale the weapon icons so they fit in the abilities bar
		//fractionOfAbilitiesBar is the height of the icon as a fraction of the height of the abilities bar
		scalefactor = pixel.V(abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*32/40, abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*32/40)
		scaled      = pixel.IM.ScaledXY(pixel.ZV, scalefactor)

		//locations of weapon icons
		knifeLocation   = pixel.Vec{abilitiesBarPosX - (abilitiesBar.Picture().Bounds().Max.X / 2.8), abilitiesBarPosY}
		handgunLocation = pixel.Vec{abilitiesBarPosX - (abilitiesBar.Picture().Bounds().Max.X / 8.8), abilitiesBarPosY}
		rifleLocation   = pixel.Vec{abilitiesBarPosX + (abilitiesBar.Picture().Bounds().Max.X / 8.5), abilitiesBarPosY}
		shotgunLocation = pixel.Vec{abilitiesBarPosX + (abilitiesBar.Picture().Bounds().Max.X / 2.8), abilitiesBarPosY}

		//myWep
		myWep = me.GetWeapon().Id
	)
	fmt.Println("weapon1212:", me.Weapon)
	if myWep != model.KNIFE {
		knife.Draw(win, scaled.Moved(knifeLocation))
	} else {
		knifeSelected.Draw(win, scaled.Moved(knifeLocation))

	}

	if !me.IsAvailable(model.HANDGUN) {
		handgunDark.Draw(win, scaled.Moved(handgunLocation))
	} else if myWep == model.HANDGUN {
		handgunSelected.Draw(win, scaled.Moved(handgunLocation))
	} else {
		handgun.Draw(win, scaled.Moved(handgunLocation))
	}

	if !me.IsAvailable(model.RIFLE) {
		rifleDark.Draw(win, scaled.Moved(rifleLocation))
	} else if myWep == model.RIFLE {
		rifleSelected.Draw(win, scaled.Moved(rifleLocation))
	} else {
		rifle.Draw(win, scaled.Moved(rifleLocation))
	}

	if !me.IsAvailable(model.SHOTGUN) {
		shotgunDark.Draw(win, scaled.Moved(shotgunLocation))
	} else if myWep == model.SHOTGUN {
		shotgunSelected.Draw(win, scaled.Moved(shotgunLocation))
	} else {
		shotgun.Draw(win, scaled.Moved(shotgunLocation))
	}
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	weaponText := text.New(pixel.V(abilitiesBarPosX+abilitiesBar.Picture().Bounds().Max.X/2, abilitiesBarPosY*0.9), basicAtlas)
	bulletsText := text.New(pixel.V(abilitiesBarPosX-abilitiesBar.Picture().Bounds().Max.X/2, abilitiesBarPosY*0.9), basicAtlas)
	bulletTextSize := 0.0
	fmt.Fprintln(weaponText, model.GetWeaponRef(myWep).GetName())
	if myWep != model.KNIFE {
		s := fmt.Sprint(me.GetWeapon().Bullets, me.GetWeapon().MagazineCurrent)
		bulletTextSize = bulletsText.LineHeight * 1.3 * float64(len(s))
		fmt.Fprintln(bulletsText, s)
	}
	weaponText.Draw(win, pixel.IM.Scaled(weaponText.Orig, 2))
	bulletsText.Draw(win, pixel.IM.Scaled(bulletsText.Orig, 2).Moved(pixel.V(-bulletTextSize, 0)))

	abilitiesBar.Draw(win, pixel.IM.Moved(pixel.V(abilitiesBarPosX, abilitiesBarPosY)))

}

//load health sprites
var hp = Load(config.Conf.HealthPath, "", IMG)

func DrawHealthbar(win *pixelgl.Window, me *model.Player) {
	if me == nil {
		return
	}
	var (
		//load sprites
		healthgraphic = hp["health.png"].Sprites[0]
		healthBackground = hp["healthbardark.png"].Sprites[0]
		healthBarFrame = hp["healthbarframe.png"].Sprites[0]

		//scalefactor on x and y axis so health width = ability width and health height is half of ability height
		x_scalefactor = ab["abilitiesBar.png"].Sprites[0].Picture().Bounds().Max.X/healthBarFrame.Picture().Bounds().Max.X
		y_scalefactor = ab["abilitiesBar.png"].Sprites[0].Picture().Bounds().Max.Y/healthBarFrame.Picture().Bounds().Max.Y/2

		//frame and background are scaled to have same width as ability bar
		// health is scaled according to the health of the player
		scaled      = pixel.IM.ScaledXY(pixel.ZV, pixel.V(x_scalefactor,y_scalefactor))
		healthscaled= pixel.IM.ScaledXY(pixel.ZV, pixel.V(x_scalefactor*float64(me.Health)/me.GetMaxHealth(),y_scalefactor))
		frameLocation = pixel.Vec{win.Bounds().Max.X/2,win.Bounds().Min.Y+ab["abilitiesBar.png"].Sprites[0].Picture().Bounds().Max.Y+healthgraphic.Picture().Bounds().Max.Y/4}
		healthfraction=float64((me.GetMaxHealth()-float64(me.Health))/(0.0001+me.GetMaxHealth()))
		healthLocation = pixel.Vec{
			win.Bounds().Max.X/2+(healthgraphic.Picture().Bounds().Max.X/2*x_scalefactor)*healthfraction-2,
			win.Bounds().Min.Y+ab["abilitiesBar.png"].Sprites[0].Picture().Bounds().Max.Y+healthgraphic.Picture().Bounds().Max.Y/4}

	)
	//draw background, health and frame
	healthBackground.Draw(win,scaled.Moved(frameLocation))

	fmt.Println("healthfraction987:",healthfraction)
	if healthfraction<0.95 {
		healthgraphic.Draw(win, healthscaled.Moved(healthLocation))
	}
	healthBarFrame.Draw(win, scaled.Moved(frameLocation))
}
