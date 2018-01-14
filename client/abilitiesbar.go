package client

import (
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/faiface/pixel"
	"fmt"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)



func (ah AnimationHandler) DrawAbilities() {
	//myWep
	me := ah.me
	win := ah.win
	wep, err := me.Weapon()
	if err != nil {
		return
	}
	var (
		//abilities bar
		abilitiesBar = ah.animations[Prefix("abilities","abilitiesBar")].CurrentSprite()

		//knife - used as reference
		knife = ah.animations[Prefix("abilities","knife")].CurrentSprite()

		//dimensions of abilities bar

		abPos = pixel.V(win.Bounds().Max.X/2, win.Bounds().Min.Y+abilitiesBar.Picture().Bounds().Max.Y/2).Add(me.Pos).Sub(win.Bounds().Center())

		//the factor by which we scale the weapon icons so they fit in the abilities bar
		//fractionOfAbilitiesBar is the height of the icon as a fraction of the height of the abilities bar
		scalefactor = pixel.V(abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*32/40, abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*32/40)
		scaled      = pixel.IM.ScaledXY(pixel.ZV, scalefactor)

		//locations of weapon icons
		knifeLocation   = pixel.Vec{abPos.X - (abilitiesBar.Picture().Bounds().Max.X / 2.8), abPos.Y}
		handgunLocation = pixel.Vec{abPos.X - (abilitiesBar.Picture().Bounds().Max.X / 8.8), abPos.Y}
		rifleLocation   = pixel.Vec{abPos.X + (abilitiesBar.Picture().Bounds().Max.X / 8.5), abPos.Y}
		shotgunLocation = pixel.Vec{abPos.X + (abilitiesBar.Picture().Bounds().Max.X / 2.8), abPos.Y}
		myWep           = me.WeaponType
	)

	ah.animations[abilitySpriteName(me, model.KNIFE)].CurrentSprite().Draw(win, scaled.Moved(knifeLocation))
	ah.animations[abilitySpriteName(me, model.HANDGUN)].CurrentSprite().Draw(win, scaled.Moved(handgunLocation))
	ah.animations[abilitySpriteName(me, model.RIFLE)].CurrentSprite().Draw(win, scaled.Moved(rifleLocation))
	ah.animations[abilitySpriteName(me, model.SHOTGUN)].CurrentSprite().Draw(win, scaled.Moved(shotgunLocation))

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	weaponText := text.New(pixel.V(abPos.X+abilitiesBar.Picture().Bounds().Max.X/2, abPos.Y), basicAtlas)
	bulletsText := text.New(pixel.V(abPos.X-abilitiesBar.Picture().Bounds().Max.X/2, abPos.Y), basicAtlas)
	var bulletTextSize float64
	fmt.Fprintln(weaponText, myWep.Name())
	if myWep != model.KNIFE {
		s := fmt.Sprint(wep.Bullets, wep.MagazineCurrent)
		bulletTextSize = bulletsText.LineHeight * 1.3 * float64(len(s))
		fmt.Fprintln(bulletsText, s)
	}
	weaponText.Draw(win, pixel.IM.Scaled(weaponText.Orig, 2))
	bulletsText.Draw(win, pixel.IM.Scaled(bulletsText.Orig, 2).Moved(pixel.V(-bulletTextSize, 0)))

	abilitiesBar.Draw(win, pixel.IM.Moved(pixel.V(abPos.X, abPos.Y)))

}

func abilitySpriteName(me model.Player, weapon model.WeaponType) string {
	s := weapon.Name()
	if me.WeaponType == weapon {
		s += "Selected"
	} else if !me.IsAvailable(weapon) {
		s += "Dark"
	}
	return Prefix("abilities", s)
}


func (ah AnimationHandler) DrawHealthbar() {
	var (
		win = ah.win
		me = ah.me
		//load sprites
		healthgraphic    = ah.animations[Prefix("health","health")].CurrentSprite()
		healthBackground = ah.animations[Prefix("health","healthbardark")].CurrentSprite()
		healthBarFrame   = ah.animations[Prefix("health","healthbarframe")].CurrentSprite()

		//scalefactor on x and y axis so health width = ability width and health height is half of ability height
		xScalefactor = ah.animations[Prefix("abilities","abilitiesBar")].CurrentSprite().Picture().Bounds().Max.X / healthBarFrame.Picture().Bounds().Max.X
		yScalefactor = ah.animations[Prefix("abilities","abilitiesBar")].CurrentSprite().Picture().Bounds().Max.Y / healthBarFrame.Picture().Bounds().Max.Y / 2

		//frame and background are scaled to have same width as ability bar
		// health is scaled according to the health of the player
		scaled         = pixel.IM.ScaledXY(pixel.ZV, pixel.V(xScalefactor, yScalefactor))
		healthscaled   = pixel.IM.ScaledXY(pixel.ZV, pixel.V(xScalefactor*float64(me.Health)/float64(me.GetMaxHealth()), yScalefactor))
		frameLocation  = pixel.Vec{win.Bounds().Max.X / 2, win.Bounds().Min.Y + ah.animations[Prefix("abilities","abilitiesBar")].CurrentSprite().Picture().Bounds().Max.Y + healthgraphic.Picture().Bounds().Max.Y/4}
		healthfraction = float64(me.GetMaxHealth()-me.Health) / float64(me.GetMaxHealth())
		healthLocation = pixel.Vec{
			win.Bounds().Max.X/2 + (healthgraphic.Picture().Bounds().Max.X/2*xScalefactor)*float64(healthfraction) - 2,
			win.Bounds().Min.Y + ah.animations[Prefix("abilities","abilitiesBar")].CurrentSprite().Picture().Bounds().Max.Y + healthgraphic.Picture().Bounds().Max.Y/4}
	)
	pos := me.Pos.Sub(win.Bounds().Center())
	frameLocation = frameLocation.Add(pos)
	healthLocation = healthLocation.Add(pos)

	//Draw background, health and frame
	healthBackground.Draw(win, scaled.Moved(frameLocation))

	if healthfraction < 0.95 {
		healthgraphic.Draw(win, healthscaled.Moved(healthLocation))
	}
	healthBarFrame.Draw(win, scaled.Moved(frameLocation))
}
