package images

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/image/draw"
)

var DownloadClient = &http.Client{
	Timeout: 5 * time.Second,
}

func DownloadImage(url string) (image.Image, error) {
	resp, err := DownloadClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Infof("Downloaded image from %s", url)

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Errorf("Error decoding image: %s", err)
		return nil, err
	}

	log.Infof("Decoded image from %s", url)

	return img, nil
}

func CropToCenterSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Determine the size of the square and the starting points for the crop
	size := min(width, height) // Size of the square will be the smaller dimension
	startX := bounds.Min.X + (width-size)/2
	startY := bounds.Min.Y + (height-size)/2

	// Define the rectangle that represents the square to crop to
	cropRect := image.Rect(startX, startY, startX+size, startY+size)

	// Create a new image to hold the cropped square
	squareImg := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
		for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
			// Map the pixels from the original image to the new square image
			squareImg.Set(x-startX, y-startY, img.At(x, y))
		}
	}
	return squareImg
}

func ResizeImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func FitImageToWidth(src image.Image, targetWidth int) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newHeight := (targetWidth * height) / width

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}
