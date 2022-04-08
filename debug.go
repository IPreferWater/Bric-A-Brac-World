package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func addPointForDebugWailing(angleRadian, x, y, charSize, distanceInsectBack, wideWailing float64,screen *ebiten.Image) {
	img := ebiten.NewImage(5, 5)
		img.Fill(color.White)
		opt := &ebiten.DrawImageOptions{}
	// -charSize/2 because we need to recenter the point in the middle of the image
	var x1 = (math.Cos(angleRadian)*distanceInsectBack) -  (math.Sin(angleRadian) * -wideWailing) + x - charSize/2
	var y1 = (math.Sin(angleRadian)*distanceInsectBack) + (math.Cos(angleRadian) * -wideWailing) + y
	opt.GeoM.Translate(x1, y1)
	screen.DrawImage(img, opt)

	var x2 = (math.Cos(angleRadian)*distanceInsectBack) -  (math.Sin(angleRadian) * wideWailing) + x - charSize/2
	var y2 = (math.Sin(angleRadian)*distanceInsectBack) + (math.Cos(angleRadian) * wideWailing) + y
	img2 := ebiten.NewImage(5, 5)
	img2.Fill(color.RGBA{0xff, 0, 0, 0xff})
	opt2 := &ebiten.DrawImageOptions{}
	opt2.GeoM.Translate(x2, y2)
	screen.DrawImage(img2, opt2)
}
