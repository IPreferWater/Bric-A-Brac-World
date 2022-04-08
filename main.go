package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 768
	screenHeight = 480
)

const (
	tileSize = 48
)

var (
	tilesImage  *ebiten.Image
	insectImage *ebiten.Image
)

func init() {
	tilesEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/tiles/room.png")
	if err != nil {
		log.Fatal(err)
	}

	insectImageEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/insect.png")
	if err != nil {
		log.Fatal(err)
	}

	tilesImage = tilesEbitenImage
	insectImage = insectImageEbitenImage
}

type Game struct {
	layers        [][]int
	mapCenterX    float64
	mapCenterY    float64
	mainCharacter mainCharacter
	state         State
	size          float64
	worldSpeed    float64
	frame         int
}

func (g *Game) Update() error {
	checkAction(g)

	g.frame++
	weaveUpdate := false
	//TODO remove weaveUpdate, do it at each frame
	if g.frame%1 == 0 {
		weaveUpdate = true
		g.frame = 0
	}

	m := g.mainCharacter

	//stoped to weave, do something
	if !m.weave.isWeaving && len(m.weave.weavePoints) > 0 {
		g.mainCharacter.weave.weavePoints = nil
	}

	if m.weave.isWeaving {
		if weaveUpdate {

			indexlastWeavePoint := len(m.weave.weavePoints) - 1

			//TODO not beautifull ...
			// we need to set the starting point to set where be begin to draw the path
			// after that, we can simply update the drawTo if same angle
			if indexlastWeavePoint < 2 {
				g.mainCharacter.weave.weavePoints = append(m.weave.weavePoints, weavePoint{
					x:     m.position.x,
					y:     m.position.y,
					angle: m.angle,
				})
				return nil
			}
			previous := m.weave.weavePoints[indexlastWeavePoint]
			//if same angle, just continue the vector
			if previous.angle == m.angle {
				g.mainCharacter.weave.weavePoints[indexlastWeavePoint].x = m.position.x
				g.mainCharacter.weave.weavePoints[indexlastWeavePoint].y = m.position.y
				return nil
			}

			// if angle changed, just add the new point
			g.mainCharacter.weave.weavePoints = append(m.weave.weavePoints, weavePoint{
				x:     m.position.x,
				y:     m.position.y,
				angle: m.angle,
			})
		}

	}
	return nil
}

//TODO
/*func getRotationCoordonates(x,y,angle, xDistanceFromOrigin, yDistanceFromOrigin, recenterValue float64)(float64,float64) {
	angleRadian := angle * (math.Pi / 180)
	return (math.Cos(angleRadian) * xDistanceFromOrigin) - (math.Sin(angleRadian) * yDistanceFromOrigin) + x - recenterValue, (math.Sin(angleRadian) * xDistanceFromOrigin) + (math.Cos(angleRadian) * yDistanceFromOrigin) + y
}*/
func drawWeaving(wide float64, weaves []weavePoint, screen *ebiten.Image) {
	arr := make([]coordinate, 0)

	charSize := float64(40)
	var pathW vector.Path
	wideWailing := float64(10)
		distanceInsectBack := float64(-20)
	for index, c := range weaves {

		//formula https://gamefromscratch.com/gamedev-math-recipes-rotating-one-point-around-another-point/
		angleRadian := c.angle * (math.Pi / 180)
	//	addPointForDebugWailing(angleRadian, c.x, c.y, charSize, distanceInsectBack, wideWailing, screen)
		// -charSize/2 because we need to recenter the point in the middle of the image
	/*	x1, y1 := getRotationCoordonates(c.x,c.y,c.angle,distanceInsectBack,-wideWailing, -charSize/2)
		x2, y2 := getRotationCoordonates(c.x,c.y,c.angle,distanceInsectBack,wideWailing, -charSize/2)*/
		x1 := (math.Cos(angleRadian) * distanceInsectBack) - (math.Sin(angleRadian) * -wideWailing) + c.x - charSize/2
		y1 := (math.Sin(angleRadian) * distanceInsectBack) + (math.Cos(angleRadian) * -wideWailing) + c.y
		x2 := (math.Cos(angleRadian) * distanceInsectBack) - (math.Sin(angleRadian) * wideWailing) + c.x - charSize/2
		y2 := (math.Sin(angleRadian) * distanceInsectBack) + (math.Cos(angleRadian) * wideWailing) + c.y
		if index == 0 {
			pathW.MoveTo(float32(x1), float32(y1))
			pathW.LineTo(float32(x1), float32(y1))

			arr = append(arr, coordinate{x: x2, y: y2})
			continue
		}

		pathW.LineTo(float32(x1), float32(y1))
		arr = append(arr, coordinate{x: x2, y: y2})
	}

	/*pathW.MoveTo(80, 170)
	pathW.LineTo(100, 170)
	pathW.QuadTo(150, 157.5, 100, 120)
	pathW.LineTo(90, 130)
	pathW.QuadTo(140, 157.5, 100, 160)
	pathW.LineTo(100, 160)
	pathW.LineTo(80, 170)*/

	for i := len(arr) - 1; i >= 0; i-- {
		pathW.LineTo(float32(arr[i].x), float32(arr[i].y))
	}

	emptyImage := ebiten.NewImage(3, 3)
	emptyImage.Fill(color.White)
	emptySubImage := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	vs, is := pathW.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		//rgba(22, 160, 133, 1)
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0x16 / float32(0xff)
		vs[i].ColorG = 0xa0 / float32(0xff)
		vs[i].ColorB = 0x85 / float32(0xff)
		vs[i].ColorA = 0.45
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)

	/*for _, c := range arr {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(c.x,c.y)
		screen.DrawImage(emptyImage, opt)
	}*/

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 13,
		G: 17,
		B: 23,
		A: 0,
	})
	characterDraw, charOpts := drawCharacter(g)
	screen.DrawImage(characterDraw, charOpts)

	square := ebiten.NewImage(16, 16)
	opChar := &ebiten.DrawImageOptions{}
	square.Fill(color.White)
	opChar.GeoM.Translate(float64(g.mainCharacter.position.x), float64(g.mainCharacter.position.y))
	//weaves
	if len(g.mainCharacter.weave.weavePoints) > 0 {
		drawWeaving(10, g.mainCharacter.weave.weavePoints, screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n , charX %f charY %f, charAngle %f", ebiten.CurrentTPS(), g.mainCharacter.position.x, g.mainCharacter.position.y, g.mainCharacter.angle))
}

func drawCharacter(g *Game) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Translate(-(g.mainCharacter.size / 2), -(g.mainCharacter.size / 2))
	opChar.GeoM.Rotate(g.mainCharacter.angle * 2 * math.Pi / 360)
	opChar.GeoM.Translate(float64(g.mainCharacter.position.x)-g.mainCharacter.size/2, float64(g.mainCharacter.position.y))

	return insectImage, opChar
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	g := &Game{
		state: WaitPlayerAction,
		mainCharacter: mainCharacter{
			position: coordinate{
				x: screenWidth / 2,
				y: screenHeight / 2,
			},
			speed:      1.5,
			angle:      0,
			angleSpeed: 3,
			size:       40,
			weave: weave{
				isWeaving:   false,
				weavePoints: make([]weavePoint, 0),
			},
		},

		worldSpeed: 1,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bric-A-Brac-World")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
