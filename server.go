package main

import "github.com/gustavfjorder/pixel-head/server"

const MaxRooms = 10

func main() {
	server.NewLounge(MaxRooms)
}
