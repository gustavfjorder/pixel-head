package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
)

func main() {
	pixelgl.Run(run)

}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	player1 := newplayer()
	player2 := newplayer()

	count := 0

	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		player1.idle[count].Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		player2.idle[count].Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(200, 100))))
		count = (count + 1) % 20
		win.Update()
	}

}
