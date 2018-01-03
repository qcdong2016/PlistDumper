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
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	return png.Encode(imgfile, img)
}

type ImageInfo struct {
	Rotated         bool   `plist:"rotated"`
	Frame           string `plist:"frame"`
	Offset          string `plist:"offset"`
	SourceColorRect string `plist:"sourceColorRect"`
	SourceSize      string `plist:"sourceSize"`
}

type Meta struct {
	Texture string `plist:"textureFileName"`
}

type ImagePack struct {
	Frames map[string]*ImageInfo `plist:"frames"`
	Meta   *Meta                 `plist:"metadata"`
}

func intArr(str string) []int {
	s := strings.Replace(str, "{", "", -1)
	s = strings.Replace(s, "}", "", -1)

	sA := strings.Split(s, ",")

	ret := make([]int, len(sA))
	for i, v := range sA {
		value, err := strconv.ParseFloat(v, 32)
		if err != nil {
			panic(err)
		}

		ret[i] = int(value)
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

	pack := ImagePack{}
	_, err := plist.Unmarshal(data, &pack)
	if err != nil {
		panic(err)
	}

	bigImage, err := LoadImage(pack.Meta.Texture)
	if err != nil {
		panic(err)
	}

	basename := filepath.Base(plistFile) + ".dir"

	err = os.Mkdir(basename, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for k, v := range pack.Frames {
		fmt.Println(k)

		s := intArr(v.Frame)
		var sub image.Image
		if v.Rotated {
			s[3], s[2] = s[2], s[3]
		}

		w, h := s[2], s[3]

		sub = SubImage(bigImage, s[0], s[1], w, h)

		if v.Rotated {
			w, h = h, w
			sub = RotateImage(sub)
		}

		ss := intArr(v.SourceSize)

		of := intArr(v.Offset)
		x, y := of[0], of[1]

		var box image.Rectangle
		box = image.Rect((ss[0]-w)/2+x, (ss[1]-h)/2-y, (ss[0]+w)/2+x, (ss[1]+h)/2-y)
		dst := image.NewRGBA(image.Rect(0, 0, ss[0], ss[1]))
		draw.Draw(dst, box, sub, image.Point{0, 0}, draw.Src)
		SaveImage(path.Join(basename, k), dst)
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
		dumpPlist(os.Args[1])
	}
	fmt.Println("https://github.com/qcdong2016/PlistDumper")
}
