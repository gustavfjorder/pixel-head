package components

import (
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel"
	"fmt"
	"image/color"
	"golang.org/x/image/colornames"
)

type Text struct {
	atlas   *text.Atlas
	content []string
	Color   color.Color
	Size    float64
}

func NewText() Text {
	return NewTextWithContents([]string{})
}

func NewTextWithContent(content string) Text {
	return NewTextWithContents([]string{content})
}

func NewTextWithContents(content []string) Text {
	return Text{
		Color:   colornames.Black,
		Size:    12,
		content: content,
	}
}

func (t *Text) loadFontFace() {
	face, err := LoadTTF("assets/gui/kenvector_future.ttf", t.Size)
	if err != nil {
		panic(err)
	}

	t.atlas = text.NewAtlas(face, text.ASCII)
}

func (t *Text) Draw(target pixel.Target, pos pixel.Vec, center ...bool) {
	t.loadFontFace()

	txt := text.New(pixel.ZV, t.atlas)

	txt.Color = t.Color

	for _, str := range t.content {
		//txt.Dot.X -= txt.BoundsOf(str).W() // Right align
		//txt.Dot.X -= txt.BoundsOf(str).W() / 2 // Center align
		fmt.Fprintln(txt, str)
	}

	if len(center) > 0 && center[0] {
		pos = pos.Sub(txt.Bounds().Center())
	}

	txt.Draw(target, pixel.IM.Moved(pos))
}

func (t *Text) Write(txt string) {
	t.content = append(t.content, txt)
}

