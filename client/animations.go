package main

import (
	"github.com/faiface/pixel"
	"strconv"
)

type player struct {
	idle []*pixel.Sprite
	/*move  []*pixel.Sprite
	shoot []*pixel.Sprite*/
}

func newplayer() player {
	p := player{idle: make([]*pixel.Sprite, 20) /*,move:make([]*pixel.Sprite, 20),shoot:make([]*pixel.Sprite, 3)*/}

	for i := 0; i < len(p.idle); i++ {
		sprite, _ := loadPicture("sprites/survivor/rifle/idle/survivor-idle_rifle_" + strconv.Itoa(i) + ".png")

		p.idle[i] = pixel.NewSprite(sprite, sprite.Bounds())
	}
	return p
}
