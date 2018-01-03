package main

import (
	"errors"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
)

type Animation struct {
	Sprites []*pixel.Sprite
	Cur     int
}

func loadAnimations(path string, prefix string) map[string]Animation {
	res := make(map[string]Animation)
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	anim, err := loadAnimation(path)
	if err == nil {
		res[prefix] = anim
	}
	for _, elem := range elems {
		if elem.IsDir() {
			for k, v := range loadAnimations(path+"/"+elem.Name(), prefix+"."+elem.Name()) {
				res[k] = v
			}
		}
	}
	return res
}

func loadAnimation(path string) (Animation, error) {
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	res := make([]*pixel.Sprite, len(elems))
	filesPresent := false
	i := 0
	for _, elem := range elems {
		if !elem.IsDir() {
			img, err := loadPicture(path + "/" + elem.Name())
			if err != nil {
				panic(err)
			}
			res[i] = pixel.NewSprite(img, img.Bounds())
			filesPresent = true
			i++
		}
	}
	if !filesPresent {
		return Animation{Sprites: nil, Cur: 0}, errors.New("No files were found")
	}
	return Animation{Sprites: res, Cur: 0}, nil
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

func (a *Animation) Next() (s *pixel.Sprite) {
	s = a.Sprites[a.Cur]
	a.Cur = (a.Cur + 1) % len(a.Sprites)
	return
}
