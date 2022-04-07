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

func drawWeaving(xChar, yChar, angle float64, weaves []weavePoint, screen *ebiten.Image) {
	arr := make([]coordinate, 0)

	//charSize := 10
	var pathW vector.Path
	for index, c := range weaves {

		x1:=  c.x-4 + math.Cos(c.angle)/360 
		y1:=  c.y-4 + math.Sin(c.angle)/360 

		x2:=  c.x+4 + math.Cos(c.angle)/360 
		y2:=  c.y+4 + math.Sin(c.angle)/360 

	if index == 0 {
		pathW.MoveTo(float32(x1), float32(y1))
		pathW.LineTo(float32(x1), float32(y1))

		arr = append(arr, coordinate{x: x2 , y: y2})
		continue
	}
		pathW.LineTo(float32(x1), float32(y1))

		arr = append(arr, coordinate{x: x2 , y: y2})	
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
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0xdb / float32(0xff)
		vs[i].ColorG = 0x56 / float32(0xff)
		vs[i].ColorB = 0x20 / float32(0xff)
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
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
		drawWeaving(g.mainCharacter.position.x, g.mainCharacter.position.y, g.mainCharacter.angle, g.mainCharacter.weave.weavePoints, screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n , charX %f charY %f, charAngle %f", ebiten.CurrentTPS(), g.mainCharacter.position.x, g.mainCharacter.position.y, g.mainCharacter.angle))
}

func drawCharacter(g *Game) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Translate(-(g.mainCharacter.size / 2), -(g.mainCharacter.size / 2))
	opChar.GeoM.Rotate(g.mainCharacter.angle * 2 * math.Pi / 360)
	opChar.GeoM.Translate(float64(g.mainCharacter.position.x), float64(g.mainCharacter.position.y)-(g.mainCharacter.size/2))

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
			angleSpeed: 1,
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
