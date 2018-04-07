/*
 * from https://github.com/qcdong2016/PlistDumper
 *
 */
package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	plist "github.com/DHowett/go-plist"
)

type Frame struct {
	Rect         image.Rectangle
	Offset       image.Point
	OriginalSize image.Point
	Rotated      bool
}

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
	// Frames map[string]interface{} `plist:"frames"`
	MetaData *MetaData `plist:"metadata"`
}

func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, err = png.Decode(file)
	return
}

func SaveImage(path string, img image.Image) (err error) {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	return png.Encode(imgfile, img)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

func dumpPlist(plistFile string) {
	fmt.Println(">>", plistFile)

	data, _ := ioutil.ReadFile(plistFile)

	version := Version{}
	_, err := plist.Unmarshal(data, &version)
	if err != nil {
		panic(err)
	}

	frames := map[string]Frame{}
	switch version.MetaData.Format {
	case 0:
		plistData := PlistV0{}
		_, err = plist.Unmarshal(data, &plistData)
		if err != nil {
			panic(err)
		}

		for k, v := range plistData.Frames {
			frames[k] = Frame{
				Rect:         image.Rect(v.X, v.Y, v.X+v.Width, v.Y+v.Height),
				OriginalSize: image.Point{v.OriginalWidth, v.OriginalHeight},
				Offset:       image.Point{int(v.OffsetX), int(v.OffsetY)},
				Rotated:      false,
			}
		}
	case 1:

		plistData := PlistV1{}
		_, err = plist.Unmarshal(data, &plistData)
		if err != nil {
			panic(err)
		}
		for k, v := range plistData.Frames {
			f := intArr(v.Frame)
			o := intArr(v.Offset)
			s := intArr(v.SourceSize)
			frames[k] = Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      false,
			}
		}
	case 2:

		plistData := PlistV2{}
		_, err = plist.Unmarshal(data, &plistData)
		if err != nil {
			panic(err)
		}
		for k, v := range plistData.Frames {
			f := intArr(v.Frame)
			o := intArr(v.Offset)
			s := intArr(v.SourceSize)
			frames[k] = Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      v.Rotated,
			}
		}
	case 3:

		plistData := PlistV3{}
		_, err = plist.Unmarshal(data, &plistData)
		if err != nil {
			panic(err)
		}
		for k, v := range plistData.Frames {
			f := intArr(v.TextureRect)
			o := intArr(v.SpriteOffset)
			s := intArr(v.SpriteSourceSize)
			frames[k] = Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      v.TextureRotated,
			}
		}
	}

	textureImage, err := LoadImage(filepath.Join(filepath.Dir(plistFile), version.MetaData.Texture))
	if err != nil {
		panic(err)
	}

	basename := filepath.Base(plistFile) + ".dir"

	isExists, err := PathExists(basename)
	if err != nil {
		panic(err)
	}
	if !isExists {
		err = os.Mkdir(basename, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	for k, v := range frames {
		fmt.Println(k)

		var subImage image.Image

		w, h := v.Rect.Size().X, v.Rect.Size().Y
		ox, oy := v.Offset.X, v.Offset.Y
		ow, oh := v.OriginalSize.X, v.OriginalSize.Y

		if v.Rotated {
			subImage = SubImage(textureImage, v.Rect.Min.X, v.Rect.Min.Y, h, w)
			subImage = RotateImage(subImage)
		} else {
			subImage = SubImage(textureImage, v.Rect.Min.X, v.Rect.Min.Y, w, h)
		}

		var destRect image.Rectangle
		destRect = image.Rect((ow-w)/2+ox, (oh-h)/2+ox, (ow-w)/2+ox+w, (oh-h)/2+oy+h)

		// Create the destination sprite image [Output]
		destImage := image.NewRGBA(image.Rect(0, 0, ow, oh))

		// Copy image to destination sprite image
		draw.Draw(destImage, destRect, subImage, image.Point{0, 0}, draw.Src)

		SaveImage(path.Join(basename, k), destImage)
	}
}

func main() {

	if len(os.Args) == 1 {
		filepath.Walk("./", func(fpath string, f os.FileInfo, err error) error {
			if f == nil || f.IsDir() {
				return nil
			}

			ext := path.Ext(fpath)
			if ext == ".plist" {
				dumpPlist(fpath)
			}

			return nil
		})
	} else {
		fpath := os.Args[1]

		ext := path.Ext(fpath)
		if ext == ".plist" {
			dumpPlist(fpath)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("https://github.com/qcdong2016/PlistDumper (Thanks for shines77)\n")
}
