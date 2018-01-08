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
	"github.com/faiface/pixel/imdraw"
	"github.com/gustavfjorder/pixel-head/model"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"fmt"
)

type Animation struct {
	prefix   string
	Sprites  []*pixel.Sprite
	Cur      int
	Tick     *time.Ticker
	NextAnim *Animation
	Blocking bool
}

func (a Animation) Start(s time.Duration) (Animation) {
	a.Tick = time.NewTicker(time.Second / s)
	return a
}

func (a *Animation) Next() (s *pixel.Sprite) {
	s = a.Sprites[a.Cur]
	select {
	case <-a.Tick.C:
		a.Cur = (a.Cur + 1) % len(a.Sprites)
		if a.Cur <= 0 && a.NextAnim != nil && len(a.NextAnim.Sprites) > 0 {
			a.Blocking = a.NextAnim.Blocking
			a.Sprites = a.NextAnim.Sprites
			*a.NextAnim = Animation{}
		}
	default:
		break
	}
	return
}

func (a *Animation) ChangeAnimation(other Animation, blocking bool) (e error) {
	if len(other.Sprites) <= 0 {
		e = errors.New("need non empty animation")
		return
	}
	if a.Blocking {
		*a.NextAnim = other
		a.NextAnim.Blocking = blocking
	} else {
		a.Sprites = other.Sprites
		a.Blocking = blocking
		a.Cur = 0
	}
	return
}

func LoadMap(m model.Map) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	for _, w := range m.Walls {
		imd.Color = colornames.Black
		imd.EndShape = imdraw.SharpEndShape
		imd.Push(pixel.V(w.P.X, w.P.Y), pixel.V(w.Q.X, w.Q.Y))
		imd.Line(w.Thickness)
	}
	return imd
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
				anim.prefix = prefix
				res[prefix] = anim
			}
			break
		}
	}
	return res
}

func HandleAnimations(win *pixelgl.Window, state StateLock, anims map[string]Animation, currentAnims map[string]Animation){
	center := pixel.ZV
	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	for _, player := range state.State.Players {
		transformation := pixel.IM.Rotated(center, player.Dir).Scaled(center, 0.5).Moved(player.Pos)
		movement := ""
		blocking := false
		switch {
		case player.Reload:
			movement = "reload"
			blocking = true
		case player.Shoot:
			movement = "shoot"
			blocking = true
		case player.Melee:
			movement = "melee"
			blocking = true
		case player.Moved:
			movement = "moved"
		default:
			movement = "idle"
		}
		prefix := Prefix("survivor", player.Weapon.Name, movement)
		v, ok := currentAnims[player.Id]
		if !ok {
			v = anims[prefix]
			currentAnims[player.Id] = v
		}
		if v.prefix != prefix {
			v.ChangeAnimation(anims[prefix], blocking)
		}
		v.Next().Draw(win, transformation)
	}
	for _, zombie := range state.State.Zombies {
		v, ok := currentAnims[zombie.Id]
		transformation := pixel.IM.Rotated(center, zombie.Dir).Moved(zombie.Pos)
		if !ok {
			v = anims[Prefix("zombie","idle")]
			currentAnims[zombie.Id] = v
		}
		v.Next().Draw(win, transformation)
	}
	for _, shoot := range state.State.Shoots {
		fmt.Println(shoot.Weapon)
	}
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

func loadAnimation(path string) (Animation, error) {
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
		Sprites: res,
		Cur: 0,
		Tick: nil,
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
