package sudoku

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

var (
	cellImages  []image.Image
	hCellImages []image.Image
)

func init() {
	cellImages = make([]image.Image, 10)
	hCellImages = make([]image.Image, 10)
	for i := 0; i < len(cellImages); i++ {
		cellImages[i] = loadImageOrPanic(fmt.Sprintf("img/%d.png", i))
		hCellImages[i] = loadImageOrPanic(fmt.Sprintf("img/h_%d.png", i))
	}
}

func loadImageOrPanic(fileName string) image.Image {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func CreateBoardImage(board Board, hIdx int) image.Image {
	cw := cellImages[0].Bounds().Dx()
	ch := cellImages[0].Bounds().Dy()
	cg := 2
	bgColor := image.NewUniform(color.RGBA{0x2C, 0x28, 0x2F, 0xFF})
	boardImg := image.NewRGBA(image.Rect(0, 0, (cw+cg)*9+cg, (ch+cg)*9+cg))
	draw.Draw(boardImg, boardImg.Rect, bgColor, image.Point{}, draw.Src)
	for i, v := range board {
		cx := i % 9
		cy := i / 9
		cellImage := cellImages[v]
		if i == hIdx {
			cellImage = hCellImages[v]
		}

		dx := cx*(cw+cg) + (cx/3)*cg
		dy := cy*(ch+cg) + (cy/3)*cg
		dstRect := image.Rect(dx, dy, dx+cw, dy+ch)

		draw.Draw(boardImg, dstRect, cellImage, image.Point{}, draw.Src)
	}
	return boardImg
}
