package component

import (
	"github.com/faiface/pixel"
)

type Box struct {
	Component
}

func NewBox(columns, rows float64) *Box {
	box := &Box{
		Component: Component{
			columns: columns + 2,
			rows:    rows + 2,
		},
	}

	box.loadSprite("assets/gui/Grey")

	return box
}

func (b *Box) Render() ComponentInterface {
	start := float64(1)

	for row := start; row <= b.rows; row++ {
		for column := start; column <= b.columns; column++ {
			var rect pixel.Rect

			if row == start { // bottom row

				if column == start {
					rect = b.Rects[42]
				} else if column == b.columns {
					rect = b.Rects[44]
				} else {
					rect = b.Rects[43]
				}

			} else if row == b.rows { // top row

				if column == start {
					rect = b.Rects[54]
				} else if column == b.columns {
					rect = b.Rects[56]
				} else {
					rect = b.Rects[55]
				}

			} else { // middle rows

				if column == start {
					rect = b.Rects[48]
				} else if column == b.columns {
					rect = b.Rects[50]
				} else {
					rect = b.Rects[49]
				}

			}

			// Yellow sprite
			//if row == start { // bottom row
			//
			//	if column == start {
			//		rect = b.Rects[36]
			//	} else if column == b.columns {
			//		rect = b.Rects[38]
			//	} else {
			//		rect = b.Rects[37]
			//	}
			//
			//} else if row == b.rows { // top row
			//
			//	if column == start {
			//		rect = b.Rects[48]
			//	} else if column == b.columns {
			//		rect = b.Rects[50]
			//	} else {
			//		rect = b.Rects[49]
			//	}
			//
			//} else { // middle rows
			//
			//	if column == start {
			//		rect = b.Rects[42]
			//	} else if column == b.columns {
			//		rect = b.Rects[44]
			//	} else {
			//		rect = b.Rects[43]
			//	}
			//
			//}

			place := b.bounds.Min.Add(pixel.V(
				(column - 1) * rect.Size().X,
				(row - 1) * rect.Size().Y,
			)).Add(pixel.V(b.data.Width / 2, b.data.Height / 2))

			sprite := pixel.NewSprite(b.pic, rect)
			sprite.Draw(b.Batch, pixel.IM.Moved(place))
		}
	}

	return b
}
