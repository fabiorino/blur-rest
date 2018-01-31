package blur

import (
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
)

// Blur blurs the input image and writes the output in the imgDest file
func Blur(imgSrc *os.File, imgDest *os.File, radius int) error {
	// Decode source image
	src, err := imaging.Open(imgSrc.Name())
	if err != nil {
		return err
	}

	// Blur image in using concurrency
	image := imaging.Blur(src, float64(radius))

	// Save the resulting image using JPEG format.
	return imaging.Encode(imgDest, image, imaging.PNG)
}
