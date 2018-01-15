package client

import (
	"os"
	"regexp"
	"strconv"
	"io/ioutil"
	"github.com/faiface/pixel"
	"sort"
	"image"
	"strings"
	_ "image/png"
	."github.com/gustavfjorder/pixel-head/client/animation"
)

func LoadAll(paths ...string) (res map[string]Animation) {
	res = make(map[string]Animation)
	for _, path := range paths {
		for _, animation := range Load(path, "") {
			res[animation.Prefix()] = animation
		}
	}
	return res
}

func Load(path, prefix string) (res []Animation){
	res = make([]Animation,0)
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, elem := range elems {
		del := "."
		if len(prefix) <= 0 {
			del = ""
		}
		animationType, present := AnimationIndex[elem.Name()]
		if elem.IsDir() && present {
			sprites, names := LoadSprites(path + "/" + elem.Name())
			if len(sprites) <= 0 {
				continue
			} else if animationType == Still {
				for i, sprite := range sprites {
					res = append(res, NewAnimation(elem.Name() + "." + names[i], []*pixel.Sprite{sprite}, animationType))
				}
			} else {
				res = append(res, NewAnimation(prefix+del+elem.Name(), sprites,animationType))
			}
		} else {
			loaded := Load(path + "/" + elem.Name(), prefix + del + elem.Name())
			res = append(res, loaded...)
		}
	}
	return
}

func LoadSprites(path string) (res []*pixel.Sprite, names []string) {
	elems, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	res = make([]*pixel.Sprite, 0, len(elems))
	names = make([]string, 0, len(elems))
	sort.Sort(ByString(elems))
	for _, elem := range elems {
		if elem.IsDir() {
			continue
		}
		img, err := LoadPicture(path + "/" + elem.Name())
		if err != nil {
			continue
		}
		res = append(res, pixel.NewSprite(img, img.Bounds()))
		i := strings.Index(elem.Name(), ".")
		names = append(names,elem.Name()[:i] )
	}
	return
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

func LoadSpriteSheet(deltax float64, deltay float64, total int, path string) (sprites []*pixel.Sprite) {
	pic, err := LoadPicture(path)
	if err != nil {
		panic(err)
	}
	sprites = make([]*pixel.Sprite, total)
	index := 0
	for y := pic.Bounds().Max.Y - deltay; y >= pic.Bounds().Min.Y; y = y - deltay {
		for x := pic.Bounds().Min.X; x <= pic.Bounds().Max.X; x = x + deltax {
			sprites[index] = pixel.NewSprite(pic, pixel.R(x, y, x+deltax, y+deltay))
			index++
			if index >= total {
				return sprites
			}
		}
	}
	return sprites
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

var AnimationIndex = map[string]Type{
	"idle":        NonBlocking,
	"meleeattack": Blocking,
	"move":        NonBlocking,
	"reload":      Blocking,
	"shoot":       Blocking,
	"attack01":    Blocking,
	"attack02":    Blocking,
	"attack03":    Blocking,
	"death01":     Terminal,
	"death02":     Terminal,
	"eating":      Blocking,
	"run":         NonBlocking,
	"saunter":     NonBlocking,
	"walk":        NonBlocking,
	"abilities":   Still,
	"barrel":      Still,
	"lootbox":     Still,
	"bullet":      Still,
	"health":      Still,
}
