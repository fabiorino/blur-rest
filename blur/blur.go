package blur

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/esimov/stackblur-go"
)

func Blur(imgSrc *os.File, imgDest *os.File, radius int) error {
	// Decode source image
	src, _, err := image.Decode(imgSrc)
	if err != nil {
		return err
	}

	// Blur frame by frame
	imgs := make([]image.Image, radius)
	done := make(chan struct{}, radius)
	for i := 0; i < radius; i++ {
		go func() {
			imgs[i] = stackblur.Process(src, uint32(src.Bounds().Dx()), uint32(src.Bounds().Dy()), uint32(i), done)
		}()
	}

	// Wait for all the frames to be processed
	for i := 0; i < radius; i++ {
		<-done
	}

	// Encode gif
	return encodeGIF(imgs, imgDest)
}

func encodeGIF(imgs []image.Image, imgDest *os.File) error {
	outGif := &gif.GIF{}
	for _, inPng := range imgs {
		inGif := image.NewPaletted(inPng.Bounds(), palette.Plan9)
		draw.Draw(inGif, inPng.Bounds(), inPng, image.Point{}, draw.Src)
		outGif.Image = append(outGif.Image, inGif)
		outGif.Delay = append(outGif.Delay, 0)
	}
	return gif.EncodeAll(imgDest, outGif)
}
