package client

import (
	"errors"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"time"
	"sort"
	"regexp"
	"strconv"
)

type Animation struct {
	Sprites []*pixel.Sprite
	Cur     int
	Tick    *time.Ticker
}

func (a Animation) Start(s time.Duration) Animation {
	a.Tick = time.NewTicker(time.Second / s)
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

func (a *Animation) ChangeAnimation(other Animation) (e error){
	if len(other.Sprites) <= 0 {
		e = errors.New("need non empty animation")
		return
	}
	a.Sprites = other.Sprites
	a.Cur = 0
	return
}

func LoadAnimations(path string, prefix string) map[string]Animation {
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
			for k, v := range LoadAnimations(path+"/"+elem.Name(), prefix+del+elem.Name()) {
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

type ByString []os.FileInfo

func (s ByString) Len() int{
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

func loadAnimation(path string) (Animation, error) {
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	if len(elems) <= 0 || elems[0].IsDir(){
		return Animation{}, errors.New("can only load files")
	}
	res := make([]*pixel.Sprite, len(elems))
	sort.Sort(ByString(elems))
	i := 0
	for _, elem := range elems {
		if elem.IsDir() {
			return Animation{}, errors.New("can only load files")
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

func Prefix(aps ...string) (res string){
	if len(aps) > 0 {
		res = aps[0]
	}

	for _, ap := range aps[1:] {
		res += "." + ap
	}
	return
}