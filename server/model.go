package server

import (
	"github.com/faiface/pixel"
)

type Stats struct {
	Health int
	Speed float32
}

type Player struct {
	Id  string
	Pos pixel.Vec
	Stats Stats
}

type Zoombie struct{
	Pos pixel.Vec
	Stats Stats
}

type Shoot struct{
	Start pixel.Vec
	Angle float64
	Time float64
	Speed float32
}


type Request struct{
	Id string
	CurrentWep byte
	Dir float64
	Move bool
	Shoot bool
}

func handleRequest(request Request){

}