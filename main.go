package main

import (
	"fmt"
	"image"
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
}

func (g *Game) Update() error {

	switch g.state {
	case WaitPlayerAction:
		checkAction(g)
	case AnimatePlayerAction:
		animate(g)
	case WaitWorldAction:
		worldAction(g)
	case AnimateWorldAction:
		animationWorld(g)
	}

	if g.state == WaitPlayerAction {
		//keyboard
	}
	return nil
}

func worldAction(g *Game) {
	g.state = AnimateWorldAction
}

func animationWorld(g *Game) {
	g.state = WaitPlayerAction
}

func (g *Game) Draw(screen *ebiten.Image) {
	//guess coordinate in matrice
	//yM:= screenWidth % int(g.mainCharacter.position.y)
	xT := int(g.mainCharacter.position.x / tileSize)
	yT := int(g.mainCharacter.position.y / tileSize)

	//TODO can be calculated once
	xNum := screenWidth / tileSize
	yNum := screenHeight / tileSize

	// -1 for the split one
	yTileGlogal := (float64(yT) * tileSize) - g.mainCharacter.position.y
	xTileGlogal := (float64(xT) * tileSize) - g.mainCharacter.position.x

	jMax := 0
	for j := yT - (yNum / 2) - 1; j <= yT+(yNum/2)-1; j++ {
		iMax := 0
		for i := xT - (xNum / 2) - 1; i <= xT+(xNum/2)-1; i++ {
			op := &ebiten.DrawImageOptions{}
			tile := g.layers[j][i]
			op.GeoM.Translate(float64(iMax)*tileSize+xTileGlogal, float64(jMax)*tileSize+yTileGlogal)
			sx, sy := getTileCoordinate(tile)
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
			iMax++
		}
		jMax++
	}
	characterDraw, charOpts := drawCharacter(g)
	screen.DrawImage(characterDraw, charOpts)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n xT: %d, yT %d, xNum %d, yNum %d\n charX %f charY %f", ebiten.CurrentTPS(), xT, yT, xNum, yNum, g.mainCharacter.position.x, g.mainCharacter.position.y))
}

func drawCharacter(g *Game) (*ebiten.Image, *ebiten.DrawImageOptions) {
	//put the char in the center of div

	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Translate(float64(g.mainCharacter.position.x), float64(g.mainCharacter.position.y))
	//opChar.GeoM.Scale(4,4)

	return mainCharacterImage.SubImage(image.Rect(0, 0, 32, 48)).(*ebiten.Image), opChar

}

//TODO
func getTileCoordinate(i int) (int, int) {
	if i == 1 {
		return 0, 0
	}

	if i == 2 {
		return 48, 0
	}

	if i == 6 {
		return 96, 48
	}

	return 48, 48

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		state: WaitPlayerAction,
		layers: [][]int{
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 5, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 5, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 5, 2, 1, 2, 1, 2, 1, 6, 5, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 5, 5, 5, 5, 2, 1, 5, 1, 5, 5, 5, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 5, 5, 5, 5, 5, 2, 1, 5, 5, 5, 5, 5, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
		},
		mapCenterX: 9,
		mapCenterY: 6,
		mainCharacter: mainCharacter{
			position: coordinate{
				x: 9 * tileSize,
				y: 6 * tileSize,
			},
			destination: coordinate{
				x: 9 * tileSize,
				y: 6 * tileSize,
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
