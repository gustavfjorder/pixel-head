package component

import (
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/client"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"encoding/xml"
	"golang.org/x/image/font"
	"github.com/gustavfjorder/pixel-head/assets"
	"bytes"
)

//type Data struct {
//	// Have to specify where to find episodes since this
//	// doesn't match the xml tags of the data that needs to go into it
//	Textures []Texture `xml:"SubTexture"`
//}
//
//type Texture struct {
//	Name   string  `xml:"name,attr"`
//	X      float64 `xml:"x,attr"`
//	Y      float64 `xml:"y,attr"`
//	Width  float64 `xml:"width,attr"`
//	Height float64 `xml:"height,attr"`
//}

type SpriteData struct {
	Width   float64
	Height  float64
	Margin  float64
	Columns int
	Rows    int
}

func LoadSprite(file string) (pixel.Picture, []pixel.Rect, SpriteData) {
	pic, _ := client.LoadPicture(file + ".png")

	var data SpriteData
	loadXml(file + ".xml", &data)

	var sprites []pixel.Rect
	for y := pic.Bounds().Min.Y; y < pic.Bounds().Max.Y; y += data.Height + data.Margin {
		for x := pic.Bounds().Min.X; x < pic.Bounds().Max.X; x += data.Width + data.Margin {
			sprites = append(sprites, pixel.R(x, y, x + data.Width, y + data.Height))
		}
	}

	return pic, sprites, data
}

func LoadTTF(path string, size float64) (font.Face, error) {
	data, err := assets.Asset(path)
	if err != nil {
		return nil, err
	}
	file := bytes.NewReader(data)

	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	dataFont, err := truetype.Parse(dataBytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(dataFont, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func loadXml(file string, structure interface{}) (*interface{}, error) {
	data, err := assets.Asset(file)
	if err != nil {
		return nil, err
	}
	xmlFile := bytes.NewReader(data)

	dataBytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(dataBytes, &structure)
	if err != nil {
		return nil, err
	}

	return &structure, nil
}

