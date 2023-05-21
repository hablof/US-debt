package imgmerge

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Grid holds the data for each grid
type Grid struct {
	ImageFilePath   string
	BackgroundColor color.Color
	OffsetX         int
	OffsetY         int
	Grids           []*Grid
}

// MergeImage is the struct that is responsible for merging the given images
type MergeImage struct {
	Row             []*Grid
	BaseDir         string
	GridSizeFromNth int
}

// New returns a new *MergeImage instance
func New(row []*Grid, opts ...func(*MergeImage)) *MergeImage {
	mi := &MergeImage{
		Row: row,
	}

	for _, option := range opts {
		option(mi)
	}

	return mi
}

// OptBaseDir is an functional option to set the BaseDir field
func OptBaseDir(dir string) func(*MergeImage) {
	return func(mi *MergeImage) {
		mi.BaseDir = dir
	}
}

func (m *MergeImage) readGridImage(grid *Grid) (image.Image, error) {
	imgPath := grid.ImageFilePath

	if m.BaseDir != "" {
		imgPath = path.Join(m.BaseDir, grid.ImageFilePath)
	}

	return m.ReadImageFile(imgPath)
}

func (m *MergeImage) readGridsImages() ([]image.Image, error) {
	var images []image.Image

	for _, grid := range m.Row {
		img, err := m.readGridImage(grid)
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return images, nil
}

func (m *MergeImage) ReadImageFile(path string) (image.Image, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	imgFile, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	var img image.Image
	splittedPath := strings.Split(path, ".")
	ext := splittedPath[len(splittedPath)-1]

	if ext == "jpg" || ext == "jpeg" {
		img, err = jpeg.Decode(imgFile)
	} else {
		img, err = png.Decode(imgFile)
	}

	if err != nil {
		return nil, err
	}

	return img, nil
}

func (m *MergeImage) mergeGrids(images []image.Image) (*image.RGBA, error) {
	var canvas *image.RGBA

	canvasBoundX := 0 //m.ImageCountDX * imageBoundX
	for _, img := range images {
		canvasBoundX += img.Bounds().Dx()
	}

	canvasBoundY := images[0].Bounds().Dy()

	canvasMaxPoint := image.Point{canvasBoundX, canvasBoundY}
	canvasRect := image.Rectangle{image.Point{0, 0}, canvasMaxPoint}
	canvas = image.NewRGBA(canvasRect)

	// draw grids one by one
	x := 0
	for i, grid := range m.Row {
		img := images[i]

		minPoint := image.Point{x, 0}
		maxPoint := minPoint.Add(image.Point{img.Bounds().Dx(), img.Bounds().Dy()})
		nextGridRect := image.Rectangle{minPoint, maxPoint}

		x += img.Bounds().Dx()

		if grid.BackgroundColor != nil {
			draw.Draw(canvas, nextGridRect, &image.Uniform{grid.BackgroundColor}, image.Point{}, draw.Src)
			draw.Draw(canvas, nextGridRect, img, image.Point{}, draw.Over)
		} else {
			draw.Draw(canvas, nextGridRect, img, image.Point{}, draw.Src)
		}

		if grid.Grids == nil {
			continue
		}

		// draw top layer grids
		for _, grid := range grid.Grids {
			img, err := m.readGridImage(grid)
			if err != nil {
				return nil, err
			}

			gridRect := nextGridRect.Bounds().Add(image.Point{img.Bounds().Dx(), 0})
			draw.Draw(canvas, gridRect, img, image.Point{}, draw.Over)
		}
	}

	return canvas, nil
}

// Merge reads the contents of the given file paths, merges them according to given configuration
func (m *MergeImage) Merge() (*image.RGBA, error) {
	images, err := m.readGridsImages()
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, errors.New("there is no image to merge")
	}

	return m.mergeGrids(images)
}
