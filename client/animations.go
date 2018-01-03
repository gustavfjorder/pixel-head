package main

import (
	"github.com/faiface/pixel"
	"strconv"
)

type player struct {
	idle  []*pixel.Sprite
	move  []*pixel.Sprite
	shoot []*pixel.Sprite
}

func newplayer(weapon string) player {

	p := player{idle: make([]*pixel.Sprite, 20), move: make([]*pixel.Sprite, 20), shoot: make([]*pixel.Sprite, 3)}

	for i := 0; i < len(p.idle); i++ {
		sprite, _ := loadPicture("sprites/survivor/" + weapon + "/idle/survivor-idle_rifle_" + strconv.Itoa(i) + ".png")
		p.idle[i] = pixel.NewSprite(sprite, sprite.Bounds())
	}
	for i := 0; i < len(p.move); i++ {
		sprite, _ := loadPicture("sprites/survivor/" + weapon + "/move/survivor-move_rifle_" + strconv.Itoa(i) + ".png")
		p.move[i] = pixel.NewSprite(sprite, sprite.Bounds())
	}
	for i := 0; i < len(p.shoot); i++ {
		sprite, _ := loadPicture("sprites/survivor/" + weapon + "/shoot/survivor-shoot_rifle_" + strconv.Itoa(i) + ".png")
		p.shoot[i] = pixel.NewSprite(sprite, sprite.Bounds())
	}

	return p
}
