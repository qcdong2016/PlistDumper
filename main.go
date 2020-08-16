package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type Frame struct {
	Rect         image.Rectangle
	Offset       image.Point
	OriginalSize image.Point
	Rotated      bool
}

func LoadImage(path string) (img image.Image, err error) {
	return imaging.Open(path)
}

func SaveImage(path string, img image.Image) (err error) {

	if filepath.Ext(path) == "" {
		path = path + ".png"
	}

	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	return png.Encode(imgfile, img)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func GetFiles(dir string, allow []string) []string {

	allowMap := map[string]bool{}
	if allow != nil {
		for _, v := range allow {
			allowMap[v] = true
		}
	}

	ret := []string{}
	filepath.Walk(dir, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		ext := path.Ext(fpath)
		if allowMap[ext] {
			ret = append(ret, filepath.ToSlash(fpath))
		}

		return nil
	})

	return ret
}

type DumpContext struct {
	FileName    string
	FileContent []byte
	Frames      map[string]Frame
	ImageFile   string
}

func dumpFrames(frames map[string]Frame, textureFileName, outdir string) error {
	textureImage, err := LoadImage(textureFileName)
	if err != nil {
		return fmt.Errorf("open image error:" + textureFileName)
	}

	if !IsDir(outdir) {
		err = os.Mkdir(outdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for k, v := range frames {
		fmt.Println(k)

		var subImage image.Image

		w, h := v.Rect.Size().X, v.Rect.Size().Y
		ox, oy := v.Offset.X, v.Offset.Y
		ow, oh := v.OriginalSize.X, v.OriginalSize.Y
		x, y := v.Rect.Min.X, v.Rect.Min.Y

		if v.Rotated {
			subImage = imaging.Crop(textureImage, image.Rect(x, y, x+h, y+w))
			subImage = imaging.Rotate90(subImage)
		} else {
			subImage = imaging.Crop(textureImage, image.Rect(x, y, x+w, y+h))
		}

		destImage := image.NewRGBA(image.Rect(0, 0, ow, oh))
		newImage := imaging.Paste(destImage, subImage, image.Point{(ow-w)/2 + ox, (oh-h)/2 - oy})

		SaveImage(path.Join(outdir, k), newImage)
	}

	return nil
}

func dumpByFileName(filename string) {

	c := DumpContext{
		FileName: filename,
		Frames:   map[string]Frame{},
	}

	data, _ := ioutil.ReadFile(c.FileName)
	c.FileContent = data

	var err error

	ext := path.Ext(filename)
	switch ext {
	case ".plist":
		err = dumpPlist(&c)
	case ".json":
		err = dumpJson(&c)
	case ".fnt":
		err = dumpFnt(&c)
	default:
		return
	}

	if err != nil {
		panic(err)
	}

	err = dumpFrames(c.Frames, c.ImageFile, c.FileName+".dir")
	if err != nil {
		panic(err)
	}
}

func doDump(path string) {
	allfiles := []string{}

	if IsDir(path) {
		files := GetFiles(path, []string{".json", ".plist", ".fnt"})
		allfiles = append(allfiles, files...)
	}

	fmt.Println(fmt.Sprintf("开始导出：共（%d）个", len(allfiles)))
	for i, v := range allfiles {
		fmt.Println(fmt.Sprintf("导出 %d/%d %s", i+1, len(allfiles), v))
		dumpByFileName(v)
	}
}

func main() {

	if len(os.Args) == 1 {
		doDump("./")
	} else {
		doDump(os.Args[1])
	}

	fmt.Printf("\n")
	fmt.Printf("好用请给个Star，谢谢.\n")
	fmt.Printf("https://github.com/qcdong2016/PlistDumper.git\n")
}
