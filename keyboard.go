package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func contains(keys []ebiten.Key, key ebiten.Key) bool {
	for _, v := range keys {
		if v == key {
			return true
		}
	}
	return false
}

func checkAction(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	//key w = wait
	if contains(keysPressed, ebiten.KeyW) {
		return
	}

	if contains(keysPressed, ebiten.KeyArrowDown) {
		g.mainCharacter.destination.y += tileSize
	}

	if contains(keysPressed, ebiten.KeyArrowUp) {
		g.mainCharacter.destination.y -= tileSize
	}

	if contains(keysPressed, ebiten.KeyArrowLeft) {
		g.mainCharacter.destination.x -= tileSize
	}

	if contains(keysPressed, ebiten.KeyArrowRight) {
		g.mainCharacter.destination.x += tileSize
	}
	g.state = AnimatePlayerAction
}

func animate(g *Game){
	g.mainCharacter.animate(g)
}
