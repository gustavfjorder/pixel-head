package component

import (
	"github.com/faiface/pixel"
	"fmt"
)

type Box struct {
	Component
}

func NewBox(width, height float64) *Box {
	box := &Box{
		Component: Component{
			columns: width + 2,
			rows:    height + 2,
		},
	}

	box.loadSprite("assets/gui/Grey")

	return box
}

func (b *Box) Render() ComponentInterface {
	fmt.Println("Render")
	for row := 0.0; row < b.rows + 2; row++ {
		for column := 0.0; column < b.columns + 2; column++ {
			var rect pixel.Rect

			if row == 0 { // bottom row

				if column == 0 {
					rect = b.Rects[42]
				} else if column == b.columns + 1 {
					rect = b.Rects[44]
				} else {
					rect = b.Rects[43]
				}

			} else if row == b.rows + 1 { // top row

				if column == 0 {
					rect = b.Rects[54]
				} else if column == b.columns+1 {
					rect = b.Rects[56]
				} else {
					rect = b.Rects[55]
				}

			} else { // middle rows

				if column == 0 {
					rect = b.Rects[48]
				} else if column == b.columns + 1 {
					rect = b.Rects[50]
				} else {
					rect = b.Rects[49]
				}

			}

			// Yellow sprite
			//if row == 0 { // bottom row
			//
			//	if column == 0 {
			//		rect = b.Rects[36]
			//	} else if column == b.columns + 1 {
			//		rect = b.Rects[38]
			//	} else {
			//		rect = b.Rects[37]
			//	}
			//
			//} else if row == b.rows + 1 { // top row
			//
			//	if column == 0 {
			//		rect = b.Rects[48]
			//	} else if column == b.columns+1 {
			//		rect = b.Rects[50]
			//	} else {
			//		rect = b.Rects[49]
			//	}
			//
			//} else { // middle rows
			//
			//	if column == 0 {
			//		rect = b.Rects[42]
			//	} else if column == b.columns + 1 {
			//		rect = b.Rects[44]
			//	} else {
			//		rect = b.Rects[43]
			//	}
			//
			//}

			place := b.bounds.Min.Add(pixel.V(
				float64(column) * rect.Size().X,
				float64(row) * rect.Size().Y,
			)).Add(pixel.V(b.data.Width / 2, b.data.Height / 2))

			sprite := pixel.NewSprite(b.pic, rect)
			sprite.Draw(b.Batch, pixel.IM.Moved(place))
		}
	}

	return b
}
