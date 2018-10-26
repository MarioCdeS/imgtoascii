package converter

import (
	"fmt"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/MarioCdeS/imgtoascii/converter/config"
)

const (
	ramp70 = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;" +
		":,\"^`'. "
	ramp10 = " .:-=+*#%@"
)

type Error struct {
	Msg   string
	Cause error
}

func (e *Error) Error() string {
	if e == nil {
		return "unknown error"
	}

	return e.Msg
}

func Run(cfg *config.Config) *Error {
	fmt.Println(*cfg)

	img, err := loadImage(cfg.ImagePath)

	if err != nil {
		return &Error{"unable to load image", err}
	}

	fmt.Println(img.Bounds())

	return nil
}

func loadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)
	return img, err
}
