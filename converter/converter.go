package converter

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/MarioCdeS/imgtoascii/converter/config"
)

const (
	ramp10 = " .:-=+*#%@"
	ramp70 = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
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

type internalConfig struct {
	*config.Config
	colWidth  int
	rowHeight int
	outRows   int
	outRamp   []rune
}

func Run(cfg *config.Config) *Error {
	img, errLoad := loadImage(cfg.ImagePath)

	if errLoad != nil {
		return &Error{"unable to load image", errLoad}
	}

	intCfg, err := calculateInternalConfig(cfg, img.Bounds())

	if err != nil {
		return err
	}

	convertToASCII(img, intCfg)
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

func calculateInternalConfig(cfg *config.Config, imgBounds image.Rectangle) (*internalConfig, *Error) {
	imgWidth := imgBounds.Max.X - imgBounds.Min.X

	if cfg.OutCols > imgWidth {
		return nil, &Error{
			fmt.Sprintf("image size is too small for the specified number of output columns (%d)", cfg.OutCols),
			nil,
		}
	}

	imgHeight := imgBounds.Max.Y - imgBounds.Min.Y
	colWidth := imgWidth / cfg.OutCols
	rowHeight := int(float64(colWidth) * cfg.ColRowRatio)
	outRows := imgHeight / rowHeight

	if outRows > imgHeight {
		return nil, &Error{
			fmt.Sprintf("image size is too small for the calculated number of output rows (%d)", outRows),
			nil,
		}
	}

	var outRampStr string

	if cfg.Ramp == config.Ramp10 {
		outRampStr = ramp10
	} else {
		outRampStr = ramp70
	}

	return &internalConfig{
		cfg,
		colWidth,
		rowHeight,
		outRows,
		[]rune(outRampStr),
	}, nil
}

func convertToASCII(img image.Image, cfg *internalConfig) {
	line := make([]rune, cfg.OutCols)

	for j := 0; j < cfg.outRows; j++ {
		for i := 0; i < cfg.OutCols; i++ {
			rect := image.Rect(i*cfg.colWidth, j*cfg.rowHeight, (i+1)*cfg.colWidth, (j+1)*cfg.rowHeight)
			idx := int(math.Floor((pixelsGrayAverage(img, &rect) / math.MaxUint8) * float64(len(cfg.outRamp)-1)))
			line[i] = cfg.outRamp[idx]
		}

		fmt.Println(string(line))
	}
}

func pixelsGrayAverage(img image.Image, rect *image.Rectangle) float64 {
	var total float64

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			total += float64(color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y)
		}
	}

	return total / float64((rect.Max.X-rect.Min.X)*(rect.Max.Y-rect.Min.Y))
}
