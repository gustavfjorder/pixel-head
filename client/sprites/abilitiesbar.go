package sprites

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/client"
)
var(
	//abilities bar
	abilitiesBar,_  = client.LoadPicture("client/sprites/abilities/abilitiesBar.png")
	abilitiesSprite = pixel.NewSprite(abilitiesBar,abilitiesBar.Bounds())

	//knife
	knifeIcon,_ = client.LoadPicture("client/sprites/abilities/knife.png")
	knifeSprite = pixel.NewSprite(knifeIcon,knifeIcon.Bounds())

	//handgun
	handgunIcon,_ = client.LoadPicture("client/sprites/abilities/handgun.png")
	handgunSprite = pixel.NewSprite(handgunIcon,handgunIcon.Bounds())

	handgunDarkIcon,_ = client.LoadPicture("client/sprites/abilities/handgunDark.png")
	handgunDarkSprite = pixel.NewSprite(handgunDarkIcon,handgunDarkIcon.Bounds())

	//rifle
	rifleIcon,_ = client.LoadPicture("client/sprites/abilities/rifle.png")
	rifleSprite = pixel.NewSprite(rifleIcon,rifleIcon.Bounds())

	rifleDarkIcon,_ = client.LoadPicture("client/sprites/abilities/rifleDark.png")
	rifleDarkSprite = pixel.NewSprite(rifleIcon,rifleIcon.Bounds())

	//shotgun
	shotgunIcon,_ = client.LoadPicture("client/sprites/abilities/shotgun.png")
	shotgunSprite = pixel.NewSprite(shotgunIcon,shotgunIcon.Bounds())

	shotgunDarkIcon,_ = client.LoadPicture("client/sprites/abilities/shotgunDark.png")
	shotgunDarkSprite = pixel.NewSprite(shotgunIcon,shotgunIcon.Bounds())
)

func drawAbilities(win *pixelgl.Window, me model.Player){
	abilitiesSprite.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Max.X/2,win.Bounds().Min.Y+abilitiesBar.Bounds().Max.Y/2)))

	if me.IsAvailable(model.Knife) {
		scalefactor:=pixel.V(abilitiesBar.Bounds().Max.Y/knifeIcon.Bounds().Max.Y*2/3,abilitiesBar.Bounds().Max.Y/knifeIcon.Bounds().Max.Y*2/3)
		movelocation:=pixel.V((win.Bounds().Max.X-abilitiesBar.Bounds().Max.X)-(abilitiesBar.Bounds().Max.X/4),abilitiesBar.Bounds().Max.Y+knifeIcon.Bounds().Max.Y*1.1)

		knifeSprite.Draw(win, pixel.IM.Moved(movelocation).ScaledXY(win.Bounds().Center(),scalefactor))
	}

	if me.IsAvailable(model.Handgun){
		handgunSprite.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Bounds().Max.X)-(abilitiesBar.Bounds().Max.X/4)-50,abilitiesBar.Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3,abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3)))
	} else{
		handgunDarkSprite.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Bounds().Max.X)-(abilitiesBar.Bounds().Max.X/4)-50,abilitiesBar.Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3,abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3)))
	}

	if me.IsAvailable(model.Rifle){
		rifleSprite.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Bounds().Max.X)-(abilitiesBar.Bounds().Max.X/4),abilitiesBar.Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Bounds().Max.Y/rifleIcon.Bounds().Max.Y*2/3,abilitiesBar.Bounds().Max.Y/rifleIcon.Bounds().Max.Y*2/3)))
	}else{
		rifleDarkSprite.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Bounds().Max.X)-(abilitiesBar.Bounds().Max.X/4)-50,abilitiesBar.Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3,abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3)))
	}

	if me.IsAvailable(model.Shotgun){
		shotgunSprite.Draw(win,pixel.IM.Moved(pixel.V(abilitiesBar.Bounds().Max.X-15500, abilitiesBar.Bounds().Max.Y/2-2)))
	}else{
		shotgunDarkSprite.Draw(win, pixel.IM.Moved(pixel.V((win.Bounds().Max.X-abilitiesBar.Bounds().Max.X)-(abilitiesBar.Bounds().Max.X/4)-50,abilitiesBar.Bounds().Max.Y-525/*+handgunIcon.Bounds().Max.Y*1.1*/)).ScaledXY(win.Bounds().Center(),pixel.V(abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3,abilitiesBar.Bounds().Max.Y/handgunIcon.Bounds().Max.Y*2/3)))
	}
}
