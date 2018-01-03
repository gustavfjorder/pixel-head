package main

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	//"sort"
	"encoding/json"
	"time"
)

var config = LoadConfiguration("conf.json")

type Config struct {
	AnimationSpeed time.Duration
	Fps            time.Duration
}

type Animation struct {
	Sprites []*pixel.Sprite
	Cur     int
	Tick    *time.Ticker
}

func (a Animation) start() Animation {
	a.Tick = time.NewTicker(time.Second / config.AnimationSpeed)
	return a
}

func (a *Animation) Next() (s *pixel.Sprite) {
	s = a.Sprites[a.Cur]
	select {
	case <-a.Tick.C:
		a.Cur = (a.Cur + 1) % len(a.Sprites)
	default:
		break
	}
	return
}

func loadAnimations(path string, prefix string) map[string]Animation {
	res := make(map[string]Animation)
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, elem := range elems {
		if elem.IsDir() {
			del := "."
			if len(prefix) <= 0 {
				del = ""
			}
			for k, v := range loadAnimations(path+"/"+elem.Name(), prefix+del+elem.Name()) {
				res[k] = v
			}
		} else {
			anim, err := loadAnimation(path)
			if err == nil {
				res[prefix] = anim
			}
			break
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
	i := 0
	for _, elem := range elems {
		fmt.Println(elem.Name())
		if elem.IsDir() {
			return Animation{Sprites: nil, Cur: 0, Tick: nil}, errors.New("can only load files")
		}
		img, err := loadPicture(path + "/" + elem.Name())
		if err != nil {
			panic(err)
		}
		res[i] = pixel.NewSprite(img, img.Bounds())
		i++

	}
	return Animation{Sprites: res, Cur: 0, Tick: nil}, nil
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

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
