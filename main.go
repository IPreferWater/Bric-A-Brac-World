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

			//if the player didn't moove
			if previous.x == m.position.x && previous.y == m.position.y {
				return nil
			}

			if previous.angle == m.angle {
				g.mainCharacter.weave.weavePoints[indexlastWeavePoint].x = m.position.x
				g.mainCharacter.weave.weavePoints[indexlastWeavePoint].y = m.position.y
				return nil
			}

			// if angle changed, just add the new point
			//we add 2 points, one for the start of the new rectangle, one for the end
			g.mainCharacter.weave.weavePoints = append(m.weave.weavePoints, weavePoint{
				x:     m.position.x,
				y:     m.position.y,
				angle: m.angle,
			},weavePoint{
				x:     m.position.x,
				y:     m.position.y,
				angle: m.angle,
			})
		}

	}
	return nil
}

func drawWeaving(wide float64, weaves []weavePoint, screen *ebiten.Image) {
	arr := make([]coordinate, 0)

	charSize := float64(40)
	var pathW vector.Path
	wideWailing := float64(10)
		distanceInsectBack := float64(-2)
	for index, c := range weaves {

		//formula https://gamefromscratch.com/gamedev-math-recipes-rotating-one-point-around-another-point/
		angleRadian := c.angle * (math.Pi / 180)
		//addPointForDebugWailing(angleRadian, c.x, c.y, charSize, distanceInsectBack, wideWailing, screen)

		x1 := (math.Cos(angleRadian) * distanceInsectBack) - (math.Sin(angleRadian) * -wideWailing) + c.x - charSize/2
		y1 := (math.Sin(angleRadian) * distanceInsectBack) + (math.Cos(angleRadian) * -wideWailing) + c.y
		x2 := (math.Cos(angleRadian) * distanceInsectBack) - (math.Sin(angleRadian) * wideWailing) + c.x - charSize/2
		y2 := (math.Sin(angleRadian) * distanceInsectBack) + (math.Cos(angleRadian) * wideWailing) + c.y

		if index == 0 {
			pathW.MoveTo(float32(x1), float32(y1))
			pathW.LineTo(float32(x1), float32(y1))

			arr = append(arr, coordinate{x: x2, y: y2})

			
			x0 := (math.Cos(angleRadian) * distanceInsectBack) - (math.Sin(angleRadian)) + c.x - charSize/2
			y0 := (math.Sin(angleRadian) * distanceInsectBack) + (math.Cos(angleRadian)) + c.y
			getStartingPointImage(x0,y0, angleRadian, distanceInsectBack, charSize, screen)
			
			continue
		}

		pathW.LineTo(float32(x1), float32(y1))
		arr = append(arr, coordinate{x: x2, y: y2})
	}

	for i := len(arr) - 1; i >= 0; i-- {
		pathW.LineTo(float32(arr[i].x), float32(arr[i].y))
	}

	/*for index, w := range weaves {
		//√(x2−x1)2+(y2−y1)2
		if index ==0 {
			//skip
			continue
		}
		wr := arr[index-1]
		distance := math.Sqrt(math.Pow(wr.x-w.x,2)+math.Pow(wr.y-w.y,2))
		
		//fmt.Printf("i %d x1 %f y1 %f x2 %f y2 %f angle %f distance %f\n", index, w.x, w.y, arr[index].x, arr[index].y, w.angle, distance)
	}
	fmt.Println("***")*/

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
}



func getStartingPointImage(x,y,angle, distanceInsectBack, charSize float64, screen *ebiten.Image){
	purpleClr := color.RGBA{255, 0, 255, 255}
	
	radius64 := float64(20)
	minAngle := math.Acos(1 - 1/radius64)

	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius64 * math.Cos(angle)
		yDelta := radius64 * math.Sin(angle)

		x1 := int(math.Round(float64(x) + xDelta))
		y1 := int(math.Round(float64(y) + yDelta))

		screen.Set(x1, y1, purpleClr)
	}
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
