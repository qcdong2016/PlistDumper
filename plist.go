package main

import (
	"image"
	"image/draw"
	"path/filepath"
	"strconv"
	"strings"

	"howett.net/plist"
)

type FrameV0 struct {
	Height         int     `plist:"height"`
	Width          int     `plist:"width"`
	X              int     `plist:"x"`
	Y              int     `plist:"y"`
	OriginalWidth  int     `plist:"originalWidth"`
	OriginalHeight int     `plist:"originalHeight"`
	OffsetX        float32 `plist:"offsetX"`
	OffsetY        float32 `plist:"offsetY"`
}
type PlistV0 struct {
	Frames map[string]*FrameV0 `plist:"frames"`
}

type FrameV1 struct {
	Frame      string `plist:"frame"`
	Offset     string `plist:"offset"`
	SourceSize string `plist:"sourceSize"`
}
type PlistV1 struct {
	Frames map[string]*FrameV1 `plist:"frames"`
}

type FrameV2 struct {
	Rotated         bool   `plist:"rotated"`
	Frame           string `plist:"frame"`
	Offset          string `plist:"offset"`
	SourceColorRect string `plist:"sourceColorRect"`
	SourceSize      string `plist:"sourceSize"`
}
type PlistV2 struct {
	Frames map[string]*FrameV2 `plist:"frames"`
}

type FrameV3 struct {
	//Aliases      []interface{} `plist:"aliases"`
	SpriteOffset     string `plist:"spriteOffset"`
	SpriteSize       string `plist:"spriteSize"`
	SpriteSourceSize string `plist:"spriteSourceSize"`
	TextureRect      string `plist:"textureRect"`
	TextureRotated   bool   `plist:"textureRotated"`
}
type PlistV3 struct {
	Frames map[string]*FrameV3 `plist:"frames"`
}

type MetaData struct {
	Format      int    `plist:"format"`
	RealTexture string `plist:"realTextureFileName"`
	Size        string `plist:"size"`
	SmartUpdate string `plist:"smartupdate"`
	Texture     string `plist:"textureFileName"`
}

type Version struct {
	MetaData *MetaData `plist:"metadata"`
}

func intArr(str string) []int {
	ret := make([]int, 0)
	s := strings.Replace(str, "{", "", -1)
	s = strings.Replace(s, "}", "", -1)

	sA := strings.Split(s, ",")

	ret = make([]int, len(sA))
	for i, v := range sA {
		value, err := strconv.ParseFloat(v, 32)
		if err != nil {
			value, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				panic(err)
			}
			ret[i] = int(value)
		} else {
			ret[i] = int(value)
		}
	}

	return ret
}

func SubImage(src image.Image, x, y, w, h int) image.Image {
	r := image.Rect(0, 0, x+w, y+h)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, r, src, image.Point{x, y}, draw.Src)
	return dst
}

func RotateImage(src image.Image) image.Image {
	w := src.Bounds().Max.X
	h := src.Bounds().Max.Y
	dst := image.NewRGBA(image.Rect(0, 0, h, w))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dst.Set(y, w-x, src.At(x, y))
		}
	}

	return dst
}

func dumpPlist(c *DumpContext) error {

	version := Version{}
	_, err := plist.Unmarshal(c.FileContent, &version)
	if err != nil {
		return err
	}

	switch version.MetaData.Format {
	case 0:
		plistData := PlistV0{}
		_, err = plist.Unmarshal(c.FileContent, &plistData)
		if err != nil {
			return err
		}

		for k, v := range plistData.Frames {
			c.Frames[k] = Frame{
				Rect:         image.Rect(v.X, v.Y, v.X+v.Width, v.Y+v.Height),
				OriginalSize: image.Point{v.OriginalWidth, v.OriginalHeight},
				Offset:       image.Point{int(v.OffsetX), int(v.OffsetY)},
				Rotated:      false,
			}
		}
	case 1:

		plistData := PlistV1{}
		_, err = plist.Unmarshal(c.FileContent, &plistData)
		if err != nil {
			return err
		}
		for k, v := range plistData.Frames {
			f := intArr(v.Frame)
			o := intArr(v.Offset)
			s := intArr(v.SourceSize)
			c.Frames[k] = Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      false,
			}
		}
	case 2:

		plistData := PlistV2{}
		_, err = plist.Unmarshal(c.FileContent, &plistData)
		if err != nil {
			return err
		}
		for k, v := range plistData.Frames {
			f := intArr(v.Frame)
			o := intArr(v.Offset)
			s := intArr(v.SourceSize)
			c.Frames[k] = Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      v.Rotated,
			}
		}
	case 3:

		plistData := PlistV3{}
		_, err = plist.Unmarshal(c.FileContent, &plistData)
		if err != nil {
			return err
		}
		for k, v := range plistData.Frames {
			f := intArr(v.TextureRect)
			o := intArr(v.SpriteOffset)
			s := intArr(v.SpriteSourceSize)
			c.Frames[k] = Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      v.TextureRotated,
			}
		}
	}

	c.ImageFile = filepath.Join(filepath.Dir(c.FileName), version.MetaData.Texture)

	return nil
}
