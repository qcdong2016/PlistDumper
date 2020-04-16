package main

import (
	"image"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func parseIntMap(str string) map[string]int {
	ret := map[string]int{}
	arr := strings.Split(str, " ")
	for _, kv := range arr {
		if kv == "" {
			continue
		}
		kvarr := strings.Split(kv, "=")

		i, err := strconv.Atoi(kvarr[1])
		if err != nil {
			panic(err)
		}
		ret[kvarr[0]] = i
	}
	return ret
}

func dumpFnt(c *DumpContext) error {
	content := strings.ReplaceAll(string(c.FileContent), "\r\n", "\n")

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		n := strings.Index(line, " ")

		switch line[:n] {
		case "page":
			regPage := regexp.MustCompile(`page\s+id=0\s+file="?([^"]+)"?`)
			match := regPage.FindStringSubmatch(line)
			c.ImageFile = filepath.Join(filepath.Dir(c.FileName), match[1])
		case "char":
			m := parseIntMap(line[n+1:])
			k := string(rune(m["id"]))

			switch k {
			case ":", "/", "\\", " ":
				k = strconv.Itoa(m["id"])
			}

			imgname := k + ".png"
			c.Frames[imgname] = Frame{
				Rect:         image.Rect(m["x"], m["y"], m["x"]+m["width"], m["y"]+m["height"]),
				OriginalSize: image.Point{m["width"], m["height"]},
				Offset:       image.Point{0, 0},
				Rotated:      false,
			}
		}
	}

	return nil
}
