package imagehash

import (
	"image"

	"github.com/disintegration/imaging"
)

type Types string

const (
	DhashType = Types("dhash")
	HhashType = Types("hhash")
	VhashType = Types("vhash")
)

type Command struct {
	Type   Types
	Length int
}

func DhashBatch(img image.Image, requests ...*Command) ([][]byte, error) {
	biggestLength := 0
	for _, req := range requests {
		if biggestLength < req.Length {
			biggestLength = req.Length
		}
	}

	imgGray := imaging.Grayscale(img) // Grayscale image first for performance

	// Width and height of the scaled-down image
	width, height := biggestLength+1, biggestLength

	// Downscale the image by 'biggestLength' amount for a horizonal diff.
	res := imaging.Resize(imgGray, width, height, imaging.Lanczos)

	ret := make([][]byte, len(requests))

	var err error

	for i, req := range requests {
		switch req.Type {
		case DhashType:
			var horiz, vert []byte
			// Calculate both horizontal and vertical gradients
			horiz, err = horizontalGradient(res, req.Length)
			if err != nil {
				return nil, err
			}
			vert, err = verticalGradient(res, req.Length)

			// Return the concatenated horizontal and vertical hash
			ret[i] = append(horiz, vert...)
		case HhashType:
			ret[i], err = horizontalGradient(res, req.Length)
		case VhashType:
			ret[i], err = verticalGradient(res, req.Length)
		}

		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
