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
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	return png.Encode(imgfile, img)
}

type ImageInfoStd struct {
	Rotated         bool   `plist:"rotated"`
	Frame           string `plist:"frame"`
	Offset          string `plist:"offset"`
	SourceColorRect string `plist:"sourceColorRect"`
	SourceSize      string `plist:"sourceSize"`
}

type MetaStd struct {
	Texture string `plist:"textureFileName"`
}

type ImagePackStd struct {
	Frames map[string]*ImageInfoStd `plist:"frames"`
	Meta   *MetaStd                 `plist:"metadata"`
}

type ImageInfoCocos2dx struct {
	Aliases      []interface{} `plist:"aliases"`
	SpriteOffset string        `plist:"spriteOffset"`
	SpriteSize   string        `plist:"spriteSize"`
	SourceSize   string        `plist:"spriteSourceSize"`
	TextureRect  string        `plist:"textureRect"`
	Rotated      bool          `plist:"textureRotated"`
}

type MetaCocos2dx struct {
	Format      int    `plist:"format"`
	RealTexture string `plist:"realTextureFileName"`
	Size        string `plist:"size"`
	SmartUpdate string `plist:"smartupdate"`
	Texture     string `plist:"textureFileName"`
}

type ImagePackCocos2dx struct {
	Frames map[string]*ImageInfoCocos2dx `plist:"frames"`
	Meta   *MetaCocos2dx                 `plist:"metadata"`
}

func IsEmpty(str string) bool {
	if str == "" {
		return true
	} else {
		return false
	}
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
	if !IsEmpty(str) {
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

func dumpPlistCocos2dx(plistFile string) {
	fmt.Println(">> ", plistFile)
	data, _ := ioutil.ReadFile(plistFile)

	pack := ImagePackCocos2dx{}
	_, err := plist.Unmarshal(data, &pack)
	if err != nil {
		panic(err)
	}

	bigImage, err := LoadImage(pack.Meta.Texture)
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

	for key, value := range pack.Frames {
		fmt.Println(key)

		s := intArr(value.TextureRect)
		var subImage image.Image
		x, y := s[0], s[1]
		width, height := s[2], s[3]

		if value.Rotated {
			subImage = SubImage(bigImage, x, y, height, width)
			subImage = RotateImage(subImage)
		} else {
			subImage = SubImage(bigImage, x, y, width, height)
		}

		spriteOffset := intArr(value.SpriteOffset)
		spriteOffsetX, spriteOffsetY := spriteOffset[0], spriteOffset[1]

		spriteSize := intArr(value.SpriteSize)
		spriteWidth, spriteHeight := spriteSize[0], spriteSize[1]

		var imgRect image.Rectangle
		imgRect = image.Rect((spriteWidth-width)/2+spriteOffsetX, (spriteHeight-height)/2+spriteOffsetY,
			(spriteWidth-width)/2+spriteOffsetX+width, (spriteHeight-height)/2+spriteOffsetY+height)
		dest := image.NewRGBA(image.Rect(0, 0, spriteWidth, spriteHeight))

		draw.Draw(dest, imgRect, subImage, image.Point{0, 0}, draw.Src)

		SaveImage(path.Join(basename, key), dest)
	}
}

func dumpPlistStd(plistFile string) {
	fmt.Println(">> ", plistFile)
	data, _ := ioutil.ReadFile(plistFile)

	pack := ImagePackStd{}
	_, err := plist.Unmarshal(data, &pack)
	if err != nil {
		panic(err)
	}

	bigImage, err := LoadImage(pack.Meta.Texture)
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

	for key, value := range pack.Frames {
		fmt.Println(key)

		s := intArr(value.Frame)
		var subImage image.Image
		x, y := s[0], s[1]
		width, height := s[2], s[3]

		if value.Rotated {
			subImage = SubImage(bigImage, x, y, height, width)
			subImage = RotateImage(subImage)
		} else {
			subImage = SubImage(bigImage, x, y, width, height)
		}

		spriteOffset := intArr(value.Offset)
		spriteOffsetX, spriteOffsetY := spriteOffset[0], spriteOffset[1]

		spriteSize := intArr(value.SourceSize)
		spriteWidth, spriteHeight := spriteSize[0], spriteSize[1]

		var imgRect image.Rectangle
		imgRect = image.Rect((spriteWidth-width)/2+spriteOffsetX, (spriteHeight-height)/2+spriteOffsetY,
			(spriteWidth-width)/2+spriteOffsetX+width, (spriteHeight-height)/2+spriteOffsetY+height)
		dest := image.NewRGBA(image.Rect(0, 0, spriteWidth, spriteHeight))

		draw.Draw(dest, imgRect, subImage, image.Point{0, 0}, draw.Src)

		SaveImage(path.Join(basename, key), dest)
	}
}

func dumpPlist(format string, fpath string) {
	if format == "std" {
		// "std"
		dumpPlistStd(fpath)
	} else {
		// "cocos2dx"
		dumpPlistCocos2dx(fpath)
	}
}

func isValidFormat(plistFormat string) bool {
	if plistFormat == "std" {
		// "std"
		return true
	} else if plistFormat == "cocos2dx" {
		// "cocos2dx"
		return true
	} else {
		// unknown
		return false
	}
}

func printUsage(appName string) {
	fmt.Printf("\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("  %s [cocos2dx] [plistfile]    # dump a cocos2dx plist file\n\n", appName)

	fmt.Printf("  %s                           # dump all cocos2dx plist file\n", appName)
	fmt.Printf("  %s abc.plist                 # dump a cocos2dx plist file\n\n", appName)

	fmt.Printf("  %s cocos2dx                  # dump all cocos2dx plist file\n", appName)
	fmt.Printf("  %s std                       # dump all standard plist file\n\n", appName)

	fmt.Printf("  %s cocos2dx abc.plist        # dump a cocos2dx plist file\n", appName)
	fmt.Printf("  %s std abc.plist             # dump a standard plist file\n\n", appName)
}

func main() {
	// Default plist file format
	format := "cocos2dx"

	if len(os.Args) == 1 {
		filepath.Walk("./", func(fpath string, f os.FileInfo, err error) error {
			if f == nil || f.IsDir() {
				return nil
			}

			ext := path.Ext(fpath)
			if ext == ".plist" {
				dumpPlist(format, fpath)
			}

			return nil
		})
	} else if len(os.Args) == 2 {
		validFormat := isValidFormat(os.Args[1])
		if !validFormat {
			fpath := os.Args[1]

			ext := path.Ext(fpath)
			if ext == ".plist" {
				dumpPlist(format, fpath)
			} else {
				printUsage(os.Args[0])
			}
		} else {
			format = os.Args[1]
			filepath.Walk("./", func(fpath string, f os.FileInfo, err error) error {
				if f == nil || f.IsDir() {
					return nil
				}

				ext := path.Ext(fpath)
				if ext == ".plist" {
					dumpPlist(format, fpath)
				}

				return nil
			})
		}
	} else if len(os.Args) == 3 {
		format = os.Args[1]
		fpath := os.Args[2]

		ext := path.Ext(fpath)
		if ext == ".plist" {
			dumpPlist(format, fpath)
		} else {
			printUsage(os.Args[0])
		}
	} else {
		printUsage(os.Args[0])
	}

	fmt.Printf("\n")
	fmt.Printf("https://github.com/shines77/PlistDumper [Original author qcdong2016]\n")
}
