package component

import (
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel"
	"fmt"
	"image/color"
	"golang.org/x/image/colornames"
)

type Text struct {
	Component

	atlas   *text.Atlas
	txt		*text.Text
	content []string
	Color   color.Color
	size    float64
}

func NewText() *Text {
	return NewTextWithContents([]string{})
}

func NewTextWithContent(content string) *Text {
	return NewTextWithContents([]string{content})
}

func NewTextWithContents(content []string) *Text {
	txt := &Text{
		Color:   colornames.Black,
		size:    12,
		content: content,
	}

	txt.loadFontFace()

	txt.Center()

	return txt
}

func (t *Text) loadFontFace() {
	face, err := LoadTTF("assets/gui/kenvector_future.ttf", t.size)
	if err != nil {
		panic(err)
	}

	t.atlas = text.NewAtlas(face, text.ASCII)
	t.Text = text.New(pixel.ZV, t.atlas)
}

func (t *Text) Render() ComponentInterface {
	t.Text.Color = t.Color

	for _, str := range t.content {
		//t.Text.Dot.X -= t.Text.BoundsOf(str).W() // Right align
		//t.Text.Dot.X -= t.Text.BoundsOf(str).W() / 2 // Center align
		fmt.Fprintln(t.Text, str)
	}

	return t
}

func (t *Text) SetSize(size float64) {
	t.size = size
	t.loadFontFace()
}

func (t *Text) Write(txt string) {
	t.content = append(t.content, txt)
}

