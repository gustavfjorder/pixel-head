package component

import (
	"github.com/faiface/pixel"
)

type Box struct {
	width  int // internal bricks
	height int // internal bricks
	sprite pixel.Sprite
	pic    pixel.Picture
	rects  []pixel.Rect
	batch  *pixel.Batch
	data   SpriteData
	bounds pixel.Rect
}

func NewBox(width, height int) Box {
	return Box{
		width:  width,
		height: height,
	}
}

func (b *Box) loadSprite() {
	b.pic, b.rects, b.data = LoadSprite("assets/gui/Grey")
	b.batch = pixel.NewBatch(&pixel.TrianglesData{}, b.pic)
}

func (b *Box) Draw(target pixel.Target, pos pixel.Vec, center ...bool) {
	b.loadSprite()

	topRight := pixel.V(float64(b.width + 2) * b.data.Width, float64(b.height + 2) * b.data.Height)
	b.bounds = pixel.Rect{
		Min: pos,
		Max: pos.Add(topRight),
	}

	if len(center) > 0 && center[0] {
		b.bounds = b.bounds.Moved(b.bounds.Center().Sub(b.bounds.Max))
	}

	for row := 0; row < b.height + 2; row++ {
		for column := 0; column < b.width+2; column++ {
			var rect pixel.Rect

			if row == 0 { // bottom row

				if column == 0 {
					rect = b.rects[42]
				} else if column == b.width + 1 {
					rect = b.rects[44]
				} else {
					rect = b.rects[43]
				}

			} else if row == b.height + 1 { // top row

				if column == 0 {
					rect = b.rects[54]
				} else if column == b.width+1 {
					rect = b.rects[56]
				} else {
					rect = b.rects[55]
				}

			} else { // middle rows

				if column == 0 {
					rect = b.rects[48]
				} else if column == b.width + 1 {
					rect = b.rects[50]
				} else {
					rect = b.rects[49]
				}

			}

			// Yellow sprite
			//if row == 0 { // bottom row
			//
			//	if column == 0 {
			//		rect = b.rects[36]
			//	} else if column == b.width + 1 {
			//		rect = b.rects[38]
			//	} else {
			//		rect = b.rects[37]
			//	}
			//
			//} else if row == b.height + 1 { // top row
			//
			//	if column == 0 {
			//		rect = b.rects[48]
			//	} else if column == b.width+1 {
			//		rect = b.rects[50]
			//	} else {
			//		rect = b.rects[49]
			//	}
			//
			//} else { // middle rows
			//
			//	if column == 0 {
			//		rect = b.rects[42]
			//	} else if column == b.width + 1 {
			//		rect = b.rects[44]
			//	} else {
			//		rect = b.rects[43]
			//	}
			//
			//}

			place := b.bounds.Min.Add(pixel.V(
				float64(column) * rect.Size().X,
				float64(row) * rect.Size().Y,
			)).Add(pixel.V(b.data.Width / 2, b.data.Height / 2))

			sprite := pixel.NewSprite(b.pic, rect)
			sprite.Draw(b.batch, pixel.IM.Moved(place))
		}
	}

	b.batch.Draw(target)
}
