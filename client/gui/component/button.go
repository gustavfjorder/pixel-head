package component

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Button struct {
	Component
	Clickable

	text  *Text
}

func NewButton(columns float64) *Button {
	button := &Button{
		Component: Component{
			columns: columns + 2,
			rows:    1,
		},
	}

	button.loadSprite("assets/gui/Yellow")

	return button
}

func (b *Button) Render() ComponentInterface {
	pressedState := 0
	if b.Pressed {
		pressedState = 3
	}

	start := float64(1)
	for column := start; column <= b.columns; column++ {
		var rect pixel.Rect

		if column == start {
			rect = b.Rects[54 + pressedState]
		} else if column == b.columns {
			rect = b.Rects[56 + pressedState]
		} else {
			rect = b.Rects[55 + pressedState]
		}

		place := b.bounds.Min.Add(pixel.V(
			(column - 1) * rect.Size().X,
			0,
		)).Add(pixel.V(b.data.Width / 2, b.data.Height / 2))

		sprite := pixel.NewSprite(b.pic, rect)
		sprite.Draw(b.Batch, pixel.IM/*.Scaled(pixel.ZV, 2)*/.Moved(place))
	}

	return b
}

func (b *Button) Text(str string) {
	b.text = NewTextWithContent(str)
	b.text.Pos(pixel.V(
		b.bounds.W() / 2,
		b.bounds.H() / 2,
	))

	b.ClearChildren()
	b.Child(b.text)
}

func (b *Button) DetermineEvent(win *pixelgl.Window) {
	if ! b.bounds.Contains(win.MousePosition()) {
		return
	}

	b.DistributeEvent(win)
}
