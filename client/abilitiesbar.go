package client

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/config"
)

var ab = Load(config.Conf.AbilityPath, "", IMG)

func DrawAbilities(win *pixelgl.Window, me model.Player){
	var(
		//abilities bar
		abilitiesBar = ab["abilitiesBar.png"].Sprites[0]

		//knife

		knife = ab["knife.png"].Sprites[0]

		//handgun
		handgun = ab["handgun.png"].Sprites[0]

		handgunDark = ab["handgunDark.png"].Sprites[0]

		//rifle
		rifle = ab["rifle.png"].Sprites[0]
		rifleDark = ab["rifleDark.png"].Sprites[0]

		//shotgun
		shotgun = ab["shotgun.png"].Sprites[0]
		shotgunDark = ab["shotgunDark.png"].Sprites[0]
	)
	abilitiesX:=win.Bounds().Max.X/2
	abilitiesY:=win.Bounds().Min.Y+abilitiesBar.Picture().Bounds().Max.Y/2
	scalefactor:=pixel.V(abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*2/3)
	scaled:=pixel.IM.ScaledXY(pixel.ZV,scalefactor)
	knifeLocation:=pixel.Vec{abilitiesX-(abilitiesBar.Picture().Bounds().Max.X/2.8),abilitiesY}
	handgunLocation:=pixel.Vec{abilitiesX-(abilitiesBar.Picture().Bounds().Max.X/8.8),abilitiesY}
	rifleLocation:=pixel.Vec{abilitiesX+(abilitiesBar.Picture().Bounds().Max.X/8.8),abilitiesY}
	shotgunLocation:=pixel.Vec{abilitiesX+(abilitiesBar.Picture().Bounds().Max.X/2.8),abilitiesY}

	if me.IsAvailable(model.KNIFE)|| !me.IsAvailable(model.KNIFE) {
		knife.Draw(win,scaled.Moved(knifeLocation))
	}

	if me.IsAvailable(model.HANDGUN){
		handgun.Draw(win,scaled.Moved(handgunLocation))
	} else{
		handgunDark.Draw(win,scaled.Moved(handgunLocation))
	}

	if me.IsAvailable(model.RIFLE){
		rifle.Draw(win,scaled.Moved(rifleLocation))
	}else{
		rifleDark.Draw(win,scaled.Moved(rifleLocation))
	}

	if me.IsAvailable(model.SHOTGUN){
		shotgun.Draw(win,scaled.Moved(shotgunLocation))
	}else{
		shotgunDark.Draw(win,scaled.Moved(shotgunLocation))
	}
	abilitiesBar.Draw(win, pixel.IM.Moved(pixel.V(abilitiesX,abilitiesY)))

}
