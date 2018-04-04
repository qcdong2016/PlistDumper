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

func SubImage(srcImage image.Image, x, y, w, h int) image.Image {
	destRect := image.Rect(0, 0, w, h)
	destIamge := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(destIamge, destRect, srcImage, image.Point{x, y}, draw.Src)
	return destIamge
}

func RotateImage(srcImage image.Image) image.Image {
	width := srcImage.Bounds().Max.X
	height := srcImage.Bounds().Max.Y
	destIamge := image.NewRGBA(image.Rect(0, 0, height, width))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			destIamge.Set(y, (width-1)-x, srcImage.At(x, y))
		}
	}

	return destIamge
}

func dumpPlistCocos2dx(plistFile string) {
	fmt.Println("")
	fmt.Println(">>", plistFile)
	fmt.Println("")

	data, _ := ioutil.ReadFile(plistFile)

	pack := ImagePackCocos2dx{}
	_, err := plist.Unmarshal(data, &pack)
	if err != nil {
		panic(err)
	}

	textureImage, err := LoadImage(pack.Meta.Texture)
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

		textureRect := intArr(value.TextureRect)
		textureLeft, textureTop := textureRect[0], textureRect[1]
		textureWidth, textureHeight := textureRect[2], textureRect[3]

		spriteOffset := intArr(value.SpriteOffset)
		spriteOffsetX, spriteOffsetY := spriteOffset[0], spriteOffset[1]

		spriteSize := intArr(value.SpriteSize)
		spriteWidth, spriteHeight := spriteSize[0], spriteSize[1]

		sourceSpriteSize := intArr(value.SourceSize)
		sourceSpriteWidth, sourceSpriteHeight := sourceSpriteSize[0], sourceSpriteSize[1]

		if textureWidth != spriteWidth || textureHeight != spriteHeight {
			fmt.Printf("\n")
			fmt.Printf("Error: TextureSize is not equal to SpriteSize!\n")
			fmt.Printf("textureWidth = %d, textureHeight = %d\n", textureWidth, textureHeight)
			fmt.Printf("spriteWidth  = %d, spriteHeight  = %d\n", spriteWidth, spriteHeight)
			fmt.Printf("\n")
		}

		if sourceSpriteWidth < textureWidth || sourceSpriteHeight < textureHeight {
			fmt.Printf("\n")
			fmt.Printf("Error: SourceSpriteSize is smaller than TextureSize!\n")
			fmt.Printf("textureWidth      = %d, textureHeight      = %d\n",
				textureWidth, textureHeight)
			fmt.Printf("sourceSpriteWidth = %d, sourceSpriteHeight = %d\n",
				sourceSpriteWidth, sourceSpriteHeight)
			fmt.Printf("\n")
		}

		var subImage image.Image
		if value.Rotated {
			subImage = SubImage(textureImage, textureLeft, textureTop, textureHeight, textureWidth)
			subImage = RotateImage(subImage)
		} else {
			subImage = SubImage(textureImage, textureLeft, textureTop, textureWidth, textureHeight)
		}

		var destRect image.Rectangle
		destRect = image.Rect(0-spriteOffsetX, 0-spriteOffsetY, 0-spriteOffsetX+spriteWidth, 0-spriteOffsetY+spriteHeight)
		/*
			if spriteOffsetX < 0 && spriteOffsetY < 0 {
				destRect = image.Rect(sourceSpriteWidth+spriteOffsetX-spriteWidth, sourceSpriteHeight+spriteOffsetY-spriteHeight,
					sourceSpriteWidth+spriteOffsetX, sourceSpriteHeight+spriteOffsetY)
			} else if spriteOffsetX >= 0 && spriteOffsetY < 0 {
				destRect = image.Rect(spriteOffsetX, sourceSpriteHeight+spriteOffsetY-spriteHeight, spriteOffsetX+spriteWidth, sourceSpriteHeight+spriteOffsetY)
			} else if spriteOffsetX < 0 && spriteOffsetY >= 0 {
				destRect = image.Rect(sourceSpriteWidth+spriteOffsetX-spriteWidth, spriteOffsetY, sourceSpriteWidth+spriteOffsetX, spriteOffsetY+spriteHeight)
			} else {
				destRect = image.Rect(spriteOffsetX, spriteOffsetY, spriteOffsetX+spriteWidth, spriteOffsetY+spriteHeight)
			}
		//*/

		// Create the destination sprite image [Output]
		destImage := image.NewRGBA(image.Rect(0, 0, sourceSpriteWidth, sourceSpriteHeight))

		// Copy image to destination sprite image
		draw.Draw(destImage, destRect, subImage, image.Point{0, 0}, draw.Src)

		// Save the destination sprite image
		SaveImage(path.Join(basename, key), destImage)
	}
}

func dumpPlistStd(plistFile string) {
	fmt.Println("")
	fmt.Println(">>", plistFile)
	fmt.Println("")

	data, _ := ioutil.ReadFile(plistFile)

	pack := ImagePackStd{}
	_, err := plist.Unmarshal(data, &pack)
	if err != nil {
		panic(err)
	}

	textureImage, err := LoadImage(pack.Meta.Texture)
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
			subImage = SubImage(textureImage, x, y, height, width)
			subImage = RotateImage(subImage)
		} else {
			subImage = SubImage(textureImage, x, y, width, height)
		}

		spriteOffset := intArr(value.Offset)
		spriteOffsetX, spriteOffsetY := spriteOffset[0], spriteOffset[1]

		spriteSize := intArr(value.SourceSize)
		spriteWidth, spriteHeight := spriteSize[0], spriteSize[1]

		var destRect image.Rectangle
		destRect = image.Rect((spriteWidth-width)/2+spriteOffsetX, (spriteHeight-height)/2+spriteOffsetY,
			(spriteWidth-width)/2+spriteOffsetX+width, (spriteHeight-height)/2+spriteOffsetY+height)

		// Create the destination sprite image [Output]
		destImage := image.NewRGBA(image.Rect(0, 0, spriteWidth, spriteHeight))

		// Copy image to destination sprite image
		draw.Draw(destImage, destRect, subImage, image.Point{0, 0}, draw.Src)

		SaveImage(path.Join(basename, key), destImage)
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
	fmt.Printf("https://github.com/qcdong2016/PlistDumper [Modified by shines77]\n")
}
