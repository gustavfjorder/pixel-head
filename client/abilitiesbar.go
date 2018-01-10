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
		knifeSelected = ab["knifeSelected.png"].Sprites[0]

		//handgun
		handgun = ab["handgun.png"].Sprites[0]
		handgunDark = ab["handgunDark.png"].Sprites[0]
		handgunSelected = ab["handgunSelected.png"].Sprites[0]


		//rifle
		rifle = ab["rifle.png"].Sprites[0]
		rifleDark = ab["rifleDark.png"].Sprites[0]
		rifleSelected = ab["rifleSelected.png"].Sprites[0]


		//shotgun
		shotgun = ab["shotgun.png"].Sprites[0]
		shotgunDark = ab["shotgunDark.png"].Sprites[0]
		shotgunSelected = ab["shotgunSelected.png"].Sprites[0]

	)
	abilitiesX:=win.Bounds().Max.X/2
	abilitiesY:=win.Bounds().Min.Y+abilitiesBar.Picture().Bounds().Max.Y/2
	scalefactor:=pixel.V(abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*2/3)
	scaled:=pixel.IM.ScaledXY(pixel.ZV,scalefactor)
	knifeLocation:=pixel.Vec{abilitiesX-(abilitiesBar.Picture().Bounds().Max.X/2.8),abilitiesY}
	handgunLocation:=pixel.Vec{abilitiesX-(abilitiesBar.Picture().Bounds().Max.X/8.8),abilitiesY}
	rifleLocation:=pixel.Vec{abilitiesX+(abilitiesBar.Picture().Bounds().Max.X/8.8),abilitiesY}
	shotgunLocation:=pixel.Vec{abilitiesX+(abilitiesBar.Picture().Bounds().Max.X/2.8),abilitiesY}
	myWep:=me.Weapon
	if myWep != model.KNIFE {
		knife.Draw(win,scaled.Moved(knifeLocation))
	} else{
		knifeSelected.Draw(win,scaled.Moved(knifeLocation))
	}

	if me.IsAvailable(model.HANDGUN) {
		handgunDark.Draw(win, scaled.Moved(handgunLocation))
	} else if me.Weapon==model.HANDGUN{
		handgunSelected.Draw(win,scaled.Moved(handgunLocation))
	} else {
		handgun.Draw(win,scaled.Moved(handgunLocation))
	}

	if !me.IsAvailable(model.RIFLE){
		rifleDark.Draw(win,scaled.Moved(rifleLocation))
	}else if myWep==model.RIFLE{
		rifleSelected.Draw(win,scaled.Moved(rifleLocation))
	}else{
		rifle.Draw(win,scaled.Moved(rifleLocation))
	}

	if !me.IsAvailable(model.SHOTGUN){
		shotgunDark.Draw(win,scaled.Moved(shotgunLocation))
	}else if myWep==model.SHOTGUN{
		shotgunSelected.Draw(win,scaled.Moved(shotgunLocation))
	}else{
		shotgun.Draw(win,scaled.Moved(shotgunLocation))
	}
	abilitiesBar.Draw(win, pixel.IM.Moved(pixel.V(abilitiesX,abilitiesY)))

}
