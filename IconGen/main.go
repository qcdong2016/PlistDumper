package main

import (
	"encoding/json"
	"fmt"
	"gg"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	toml "github.com/pelletier/go-toml"
)

type Info struct {
	Version int    `json:"version"`
	Author  string `json:"author"`
}

type Meta struct {
	Floader           string
	InfoFile          string
	RoundCornerRadius float32
}

type ImageInfo struct {
	Size     string `json:"size"`
	Filename string `json:"filename"`
	Scale    string `json:"scale"`
	Idiom    string `json:"idiom"`
}

type IconSet struct {
	Info   *Info        `json:"info"`
	Meta   *Meta        `json:"-"`
	Images []*ImageInfo `json:"images"`
}

type Config struct {
	IconSet []*IconSet
}

func LoadImage(filename string, roundCornerRadius float64) (image.Image, error) {
	im, err := gg.LoadImage(filename)
	if err != nil {
		return nil, err
	}

	if roundCornerRadius != 0 {
		w, h := im.Bounds().Max.X, im.Bounds().Max.Y

		var r float64

		if roundCornerRadius < 1 {
			r = float64(w) * roundCornerRadius
		} else {
			r = roundCornerRadius
		}

		dc := gg.NewContext(w, h)
		dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), r)
		dc.Clip()
		dc.DrawImage(im, 0, 0)

		return dc.Image(), nil
	}

	return im, nil
}

func getSize(info *ImageInfo) (int, int) {
	arr := strings.Split(info.Size, "x")
	w, err := strconv.ParseFloat(arr[0], 32)
	if err != nil {
		panic(err)
	}

	h, err := strconv.ParseFloat(arr[1], 32)
	if err != nil {
		panic(err)
	}

	if info.Scale != "" {

		scale, err := strconv.Atoi(info.Scale[:1])
		if err != nil {
			panic(err)
		}

		return int(w * float64(scale)), int(h * float64(scale))
	}
	return int(w), int(h)
}

func saveImage(img image.Image, w, h int, path string) {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	dst := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)

	imgfile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer imgfile.Close()
	if err := png.Encode(imgfile, dst); err != nil {
		panic(err)
	}
	fmt.Println(path, w, h)
}

func main() {
	bytes, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	config := Config{}

	if err := toml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	filename := os.Args[1]

	for _, imset := range config.IconSet {
		img, err := LoadImage(filename, float64(imset.Meta.RoundCornerRadius))
		if err != nil {
			panic(err)
		}

		for _, imginfo := range imset.Images {
			w, h := getSize(imginfo)
			saveImage(img, w, h, filepath.Join(imset.Meta.Floader, imginfo.Filename))
		}

		if imset.Meta.InfoFile != "" {
			bytes, err := json.MarshalIndent(imset, "", "\t")
			if err != nil {
				panic(err)
			}

			outfile := filepath.Join(imset.Meta.Floader, imset.Meta.InfoFile)
			if err := ioutil.WriteFile(outfile, bytes, os.ModePerm); err != nil {
				panic(err)
			}
		}
	}

}
