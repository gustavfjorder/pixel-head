package client

import (
	"os"
	"regexp"
	"strconv"
	"io/ioutil"
	"github.com/faiface/pixel"
	"sort"
	"image"
	"errors"
)

const (
	ANIM = iota
	IMG
)

func Load(path string, prefix string, op int) map[string]Animation {
	res := make(map[string]Animation)
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, elem := range elems {
		del := "."
		if len(prefix) <= 0 {
			del = ""
		}
		if elem.IsDir() {

			for k, v := range Load(path+"/"+elem.Name(), prefix+del+elem.Name(), op) {
				res[k] = v
			}
		} else {
			if op == ANIM {
				anim, err := LoadAnimation(path)
				if err == nil {
					anim.Prefix = prefix
					res[prefix] = anim
				}
				break
			} else if op == IMG {
				pic, err := LoadPicture(path + "/" + elem.Name())
				if err == nil {
					res[del+elem.Name()] = Animation{
						Prefix:  del + elem.Name(),
						Sprites: []*pixel.Sprite{pixel.NewSprite(pic, pic.Bounds())},
					}
				}
			}
		}
	}
	return res
}

type ByString []os.FileInfo

func (s ByString) Len() int {
	return len(s)
}

func (s ByString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByString) Less(i, j int) bool {
	r := regexp.MustCompile("[0-9]+")
	si, _ := strconv.Atoi(r.FindString(s[i].Name()))
	sj, _ := strconv.Atoi(r.FindString(s[j].Name()))
	return si < sj
}

func LoadAnimation(path string) (Animation, error) {
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	if len(elems) <= 0 || elems[0].IsDir() {
		return Animation{}, errors.New("can only load files")
	}
	res := make([]*pixel.Sprite, len(elems))
	sort.Sort(ByString(elems))
	i := 0
	for _, elem := range elems {
		if elem.IsDir() {
			return Animation{}, errors.New("can only load files")
		}
		img, err := LoadPicture(path + "/" + elem.Name())
		if err != nil {
			panic(err)
		}

		res[i] = pixel.NewSprite(img, img.Bounds())
		i++

	}
	return Animation{
		Sprites:  res,
		Cur:      0,
		NextAnim: &Animation{},
	}, nil
}

func LoadPicture(path string) (pixel.Picture, error) {
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

func Prefix(aps ...string) (res string) {
	if len(aps) > 0 {
		res = aps[0]
	}

	for _, ap := range aps[1:] {
		res += "." + ap
	}
	return
}

func LoadSpriteSheet(deltax float64, deltay float64, total int, path string) (anim Animation) {
	pic, err := LoadPicture(path)
	if err != nil {
		panic(err)
	}
	sprites := make([]*pixel.Sprite, total)
	index := 0
	for y := pic.Bounds().Max.Y - deltay; y >= pic.Bounds().Min.Y; y = y - deltay {
		for x := pic.Bounds().Min.X; x <= pic.Bounds().Max.X; x = x + deltax {
			sprites[index] = pixel.NewSprite(pic, pixel.R(x, y, x+deltax, y+deltay))
			index++
			if index >= total {
				goto loopdone
			}
		}
	}
loopdone:
	anim.NextAnim = &Animation{}
	anim.Sprites = sprites
	return

}
