package client

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/config"
	"fmt"
)
var ab = Load(config.Conf.AbilityPath, "", IMG)

func DrawAbilities(win *pixelgl.Window, me model.Player){
	fmt.Print(ab)
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
	abilitiesBar.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Max.X/2,win.Bounds().Min.Y+abilitiesBar.Picture().Bounds().Max.Y/2)))

	if me.IsAvailable(model.Knife) {
		scalefactor:=pixel.V(abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/knife.Picture().Bounds().Max.Y*2/3)
		movelocation:=pixel.V((win.Bounds().Max.X-abilitiesBar.Picture().Bounds().Max.X)-(abilitiesBar.Picture().Bounds().Max.X/4),abilitiesBar.Picture().Bounds().Max.Y+knife.Picture().Bounds().Max.Y*1.1)

		knife.Draw(win, pixel.IM.Moved(movelocation).ScaledXY(win.Bounds().Center(),scalefactor))
	}

	if me.IsAvailable(model.Handgun){
		handgun.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Picture().Bounds().Max.X)-(abilitiesBar.Picture().Bounds().Max.X/4)-50,abilitiesBar.Picture().Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3)))
	} else{
		handgunDark.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Picture().Bounds().Max.X)-(abilitiesBar.Picture().Bounds().Max.X/4)-50,abilitiesBar.Picture().Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3)))
	}

	if me.IsAvailable(model.Rifle){
		rifle.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Picture().Bounds().Max.X)-(abilitiesBar.Picture().Bounds().Max.X/4),abilitiesBar.Picture().Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Picture().Bounds().Max.Y/rifle.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/rifle.Picture().Bounds().Max.Y*2/3)))
	}else{
		rifleDark.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Picture().Bounds().Max.X)-(abilitiesBar.Picture().Bounds().Max.X/4)-50,abilitiesBar.Picture().Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3)))
	}

	if me.IsAvailable(model.Shotgun){
		shotgun.Draw(win,pixel.IM.Moved(pixel.V(abilitiesBar.Picture().Bounds().Max.X-15500, abilitiesBar.Picture().Bounds().Max.Y/2-2)))
	}else{
		shotgunDark.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Picture().Bounds().Max.X)-(abilitiesBar.Picture().Bounds().Max.X/4)-50,abilitiesBar.Picture().Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3,abilitiesBar.Picture().Bounds().Max.Y/handgun.Picture().Bounds().Max.Y*2/3)))
	}
}
