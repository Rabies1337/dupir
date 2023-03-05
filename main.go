package main

import (
	"bytes"
	"fmt"
	"github.com/corona10/goimagehash"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

type Image struct {
	path string
	hash goimagehash.ImageHash
}

func main() {
	if len(os.Args) == 0 {
		log.Println("No args")
		return
	}

	entries, err := os.ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}

	images := make([]Image, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || entry.IsDir() || info.Size() == 0 {
			continue
		}

		format := strings.Split(info.Name(), ".")
		if len(format) <= 1 || (format[len(format)-1] != "jpg" && format[len(format)-1] != "jpeg" && format[len(format)-1] != "png") {
			continue
		}

		imagePath := os.Args[1] + "/" + info.Name()
		imageFormat := format[len(format)-1]
		b, err := os.ReadFile(imagePath)
		if err != nil {
			continue
		}

		var img image.Image
		switch imageFormat {
		case "png":
			img, _ = png.Decode(bytes.NewReader(b))
			break

		case "jpg":
		case "jpeg":
			img, _ = jpeg.Decode(bytes.NewReader(b))
			break

		default:
			continue
		}

		if img != nil {
			hash, _ := goimagehash.AverageHash(img)
			images = append(images, Image{
				path: imagePath,
				hash: *hash,
			})
		}
	}

	log.Println(fmt.Sprintf("Loaded (%v)", len(images)))
	for _, img := range images {
		smilierImage, at := GetSmilierImage(img, images)
		if smilierImage != nil {
			images = append(images[:at], images[at+1:]...)
			if err := os.Remove(smilierImage.path); err != nil {
				log.Println(fmt.Sprintf("Deletion failed (%v)", smilierImage.path))
			}

			log.Println(fmt.Sprintf("Remove duplicate (%v)", smilierImage.path))
		}
	}
}

func GetSmilierImage(base Image, images []Image) (*Image, int) {
	for i, img := range images {
		if base.path != img.path {
			dist, err := base.hash.Distance(&img.hash)
			if err != nil || dist > 0 {
				continue
			}
			return &img, i
		}
	}
	return nil, -1
}
