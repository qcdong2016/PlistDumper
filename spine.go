package main

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"strings"
)

const (
	statNone      = 0
	statAtlasPart = 1
	statFramePart = 2
)

type SpineAtlasPart struct {
	Name  string
	Parts map[string]*FramePart
}

type FramePart struct {
	Name   string
	Rotate string
	XY     string
	Size   string
	Orig   string
	Offset string
	Index  string
}

type SpineAtlas struct {
	Parts map[string]*SpineAtlasPart
}

func dumpSpine(c *DumpContext) error {
	reader := bufio.NewReader(bytes.NewReader(c.FileContent))

	state := 0

	atlas := SpineAtlas{}
	atlas.Parts = map[string]*SpineAtlasPart{}

	var currentAtlasPart *SpineAtlasPart
	var currentFramePart *FramePart

	for {
		linebytes, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := string(linebytes)

		if line == "" {
			state = statAtlasPart
			currentAtlasPart = nil
			continue
		}

		if state == statAtlasPart {
			if currentAtlasPart == nil {
				currentAtlasPart = &SpineAtlasPart{
					Name:  strings.TrimSpace(line),
					Parts: map[string]*FramePart{},
				}
				atlas.Parts[currentAtlasPart.Name] = currentAtlasPart
				continue
			}

			if strings.Index(line, ":") == -1 {
				state = statFramePart
			}
		}

		if state == statFramePart {
			if strings.HasPrefix(line, " ") {
				line = strings.TrimSpace(line)
				kv := strings.Split(line, ":")
				v := strings.TrimSpace(kv[1])
				switch kv[0] {
				case "rotate":
					currentFramePart.Rotate = v
				case "xy":
					currentFramePart.XY = v
				case "size":
					currentFramePart.Size = v
				case "orig":
					currentFramePart.Orig = v
				case "offset":
					currentFramePart.Offset = v
				case "index":
					currentFramePart.Index = v
				default:
					return errors.New("error parse line " + line)
				}
			} else {
				currentFramePart = nil
			}

			if currentFramePart == nil {
				currentFramePart = &FramePart{
					Name: strings.TrimSpace(line),
				}
				currentAtlasPart.Parts[currentFramePart.Name] = currentFramePart
				continue
			}
		}
	}

	for _, sp := range atlas.Parts {
		part := c.AppendPart()
		for _, spf := range sp.Parts {
			xy := intArr(spf.XY)
			orig := intArr(spf.Orig)
			offset := intArr(spf.Offset)
			size := intArr(spf.Size)

			ro := false
			if spf.Rotate == "true" {
				ro = true
			}

			part.ImageFile = sp.Name
			part.Frames[spf.Name] = &Frame{
				Rotated:      ifelse(ro, 270, 0),
				OriginalSize: image.Pt(orig[0], orig[1]),
				Offset:       image.Pt(offset[0], offset[1]),
				Rect:         image.Rect(xy[0], xy[1], xy[0]+size[0], xy[1]+size[1]),
			}
		}
	}

	return nil
}
