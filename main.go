package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 768
	screenHeight = 480
)

const (
	tileSize = 48
)

var (
	tilesImage         *ebiten.Image
	mainCharacterImage *ebiten.Image
)

func init() {
	tilesEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/tiles/room.png")
	if err != nil {
		log.Fatal(err)
	}

	mainCharacterEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/main_character.png")
	if err != nil {
		log.Fatal(err)
	}

	tilesImage = tilesEbitenImage
	mainCharacterImage = mainCharacterEbitenImage
}

type Game struct {
	layers        [][]int
	mapCenterX    float64
	mapCenterY    float64
	mainCharacter mainCharacter
	state         State
	worldSpeed    float64
	frame      int
}

func (g *Game) Update() error {
	checkAction(g)

	g.frame++
	weaveUpdate := false
	if g.frame%12 == 0 {
		weaveUpdate = true
		g.frame = 0
	}

	m := g.mainCharacter

	//stoped to weave, do something
	if m.weave.isWeaving == false && len(m.weave.coordinates) > 0 {
		g.mainCharacter.weave.coordinates = nil
	}

	if m.weave.isWeaving {
		if weaveUpdate {
			g.mainCharacter.weave.coordinates = append(m.weave.coordinates, coordinate{
				x: m.position.x,
				y: m.position.y,
			})
		}

		
	}
	return nil
}
func drawWeaving(xChar, yChar float64, weaves []coordinate,screen *ebiten.Image) {
	//TODO why reverse ?
	/*for i := len(weaves) - 1; i >= 0; i-- {
		
		coordinate := weaves[i]
		//vertical or horizontal ?
		//width := xChar - coordinate.x
		//height := yChar - coordinate.y
		weaveToDraw := ebiten.NewImage(16, 16)
		
		opChar := &ebiten.DrawImageOptions{}
		weaveToDraw.Fill(color.RGBA{
			R: 193,
			G: 136,
			B: 70,
			A: 0,
		})
		opChar.GeoM.Translate(coordinate.x, coordinate.y)
		screen.DrawImage(weaveToDraw, opChar)
	}*/
	for _, c := range weaves {
		fmt.Printf("need to draw weave x %f y %f\n",c.x, c.y)
		weaveToDraw := ebiten.NewImage(16, 16)
		opChar := &ebiten.DrawImageOptions{}
		/*weaveToDraw.Fill(color.RGBA{
			R: 193,
			G: 136,
			B: 70,
			A: 0,
		})*/
		weaveToDraw.Fill(color.White)
		opChar.GeoM.Translate(c.x, c.y)
		screen.DrawImage(weaveToDraw, opChar)
		//TODO not working well, must construct one single image, vertex ?
	}
}

func (g *Game) Draw(screen *ebiten.Image){
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
	if len(g.mainCharacter.weave.coordinates)>0 {
		drawWeaving(g.mainCharacter.position.x, g.mainCharacter.position.y, g.mainCharacter.weave.coordinates, screen)
	}
	

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n , charX %f charY %f", ebiten.CurrentTPS(), g.mainCharacter.position.x, g.mainCharacter.position.y))
}

func drawCharacter(g *Game) (*ebiten.Image, *ebiten.DrawImageOptions) {
	//put the char in the center of div
	square := ebiten.NewImage(16, 16)

	opChar := &ebiten.DrawImageOptions{}
	square.Fill(color.White)
	opChar.GeoM.Translate(float64(g.mainCharacter.position.x), float64(g.mainCharacter.position.y))

	return square, opChar

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	g := &Game{
		state: WaitPlayerAction,
		mainCharacter: mainCharacter{
			position: coordinate{
				x: 0,
				y: 0,
			},
			speed: 1,
			weave: weave{
				isWeaving:   false,
				coordinates: make([]coordinate, 0),
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
