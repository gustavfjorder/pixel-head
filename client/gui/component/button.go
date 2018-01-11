package component

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Button struct {
	width   int // internal bricks
	sprite  pixel.Sprite
	pic     pixel.Picture
	rects   []pixel.Rect
	batch   *pixel.Batch
	data    SpriteData
	Pressed bool
	Text    string
	bounds  pixel.Rect
}

func NewButton(width int) Button {
	return Button{
		width: width,
	}
}

func (b *Button) loadSprite() {
	b.pic, b.rects, b.data = LoadSprite("assets/gui/Yellow")
	b.batch = pixel.NewBatch(&pixel.TrianglesData{}, b.pic)
}

func (b *Button) Draw(target pixel.Target, pos pixel.Vec, center ...bool) {
	b.loadSprite()

	topRight := pixel.V(float64(b.width + 2) * b.data.Width, b.data.Height)
	b.bounds = pixel.Rect{
		Min: pos,
		Max: pos.Add(topRight),
	}

	if len(center) > 0 && center[0] {
		b.bounds = b.bounds.Moved(b.bounds.Center().Sub(b.bounds.Max))
	}

	adder := 0
	if b.Pressed {
		adder = 3
	}

	for column := 0; column < b.width + 2; column++ {
		var rect pixel.Rect

		if column == 0 {
			rect = b.rects[54 + adder]
		} else if column == b.width + 1 {
			rect = b.rects[56 + adder]
		} else {
			rect = b.rects[55 + adder]
		}

		place := b.bounds.Min.Add(pixel.V(
			float64(column) * rect.Size().X,
			0,
		)).Add(pixel.V(b.data.Width / 2, b.data.Height / 2))

		sprite := pixel.NewSprite(b.pic, rect)
		sprite.Draw(b.batch, pixel.IM/*.Scaled(pixel.ZV, 2)*/.Moved(place))
	}

	b.batch.Draw(target)

	// Draw text if any
	if b.Text != "" {
		txtComp := NewTextWithContent(b.Text)
		txtComp.Size = 10

		txtComp.Draw(target, b.bounds.Center(), true)
	}
}


/**
 * Event handling
 */
func (b *Button) OnLeftMouseClick(win *pixelgl.Window, handler func()) {
	b.OnClick(win, func(button pixelgl.Button) {
		if button == pixelgl.MouseButtonLeft {
			handler()
		}
	})
}

func (b *Button) OnRightMouseClick(win *pixelgl.Window, handler func()) {
	b.OnClick(win, func(button pixelgl.Button) {
		if button == pixelgl.MouseButtonRight {
			handler()
		}
	})
}

func (b *Button) OnClick(win *pixelgl.Window, handler func(button pixelgl.Button)) {
	mouse := win.MousePosition()
	if ! b.bounds.Contains(mouse) {
		return
	}

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		b.Pressed = true
		handler(pixelgl.MouseButtonLeft)
	} else if win.JustPressed(pixelgl.MouseButtonRight) {
		handler(pixelgl.MouseButtonRight)
	}

	if win.JustReleased(pixelgl.MouseButtonLeft) {
		b.Pressed = false
	}
}
