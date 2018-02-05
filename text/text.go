package text

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"sync"
)

const (
	Sans     string  = "sans.ttf"
	BaseFont float64 = 16.0
)

var (
	mutex     sync.Mutex
	fontFaces map[string]map[float64]*FontFace
	textCache map[string]*image.RGBA
)

func init() {
	fontFaces = map[string]map[float64]*FontFace{}
	textCache = map[string]*image.RGBA{}
}

type FontFace struct {
	face  font.Face
	mutex sync.Mutex
}

func getFontFace(path string, size float64) *FontFace {
	mutex.Lock()

	assert := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	size = math.Ceil(size * BaseFont)

	if _, ok := fontFaces[path]; !ok {
		fontFaces[path] = make(map[float64]*FontFace)
	}
	if _, ok := fontFaces[path][size]; !ok {

		data, eio := ioutil.ReadFile(path)
		assert(eio)

		font, ett := truetype.Parse(data)
		assert(ett)

		ff := new(FontFace)
		ff.face = truetype.NewFace(font, &truetype.Options{
			Size: size,
		})

		fontFaces[path][size] = ff
	}

	ff := fontFaces[path][size]
	mutex.Unlock()
	return ff
}

func Draw(quad color.Color, size float64, text string) *image.RGBA {

	ff := getFontFace(Sans, size)
	ff.mutex.Lock()

	dc := gg.NewContext(8, 8)
	dc.SetFontFace(ff.face)
	w, h := dc.MeasureString(text)

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h*1.1)))
	dc = gg.NewContextForRGBA(img)

	dc.SetColor(quad)
	dc.SetFontFace(ff.face)
	dc.DrawStringAnchored(text, w/2, h/2, 0.5, 0.35)

	ff.mutex.Unlock()
	return img
}

func DrawCache(quad color.Color, size float64, text string) *image.RGBA {
	key := fmt.Sprintf("%v:%f:%s", quad, size, text)
	if textCache[key] == nil {
		textCache[key] = Draw(quad, size, text)
	}
	return textCache[key]
}
