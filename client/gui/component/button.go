package component

import (
	"github.com/faiface/pixel"
)

type Button struct {
	Component
	Clickable

	Text    string
}

func NewButton(width float64) *Button {
	button := &Button{
		Component: Component{
			columns: width + 2,
			rows:    1,
		},
	}

	button.loadSprite("assets/gui/Yellow")

	return button
}

func (b *Button) Render() ComponentInterface {
	adder := 0
	if b.Pressed {
		adder = 3
	}

	for column := 0.0; column < b.columns + 2.0; column++ {
		var rect pixel.Rect

		if column == 0 {
			rect = b.Rects[54 + adder]
		} else if column == b.columns + 1.0 {
			rect = b.Rects[56 + adder]
		} else {
			rect = b.Rects[55 + adder]
		}

		place := b.bounds.Min.Add(pixel.V(
			float64(column) * rect.Size().X,
			0,
		)).Add(pixel.V(b.data.Width / 2, b.data.Height / 2))

		sprite := pixel.NewSprite(b.pic, rect)
		sprite.Draw(b.Batch, pixel.IM/*.Scaled(pixel.ZV, 2)*/.Moved(place))
	}

	// Draw text if any
	if b.Text != "" {
		txtComp := NewTextWithContent(b.Text)
		txtComp.SetSize(10)

		b.Child(txtComp)
	}

	return b
}
