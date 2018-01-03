package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	"os"
	//"strconv"
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

	player := make([]*pixel.Sprite, 20)

	for i := 0; i < 20; i++ {
		sprite, _ := loadPicture("/sprites/survivor/rifle/idle/survivor-idle_rifle_1.png" /* + strconv.atoi(i) + ".png"*/)
		player[i] = pixel.NewSprite(sprite, sprite.Bounds())
	}

	count := 0

	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		player[count].Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		count = (count + 1) % 20
		win.Update()
	}

}
