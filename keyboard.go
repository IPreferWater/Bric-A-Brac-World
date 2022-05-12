package main

import (
	"math"

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
	isKeyArrowUp := contains(keysPressed, ebiten.KeyArrowUp)
	//key w = wait
	if contains(keysPressed, ebiten.KeyW) {
		return
	}

	if contains(keysPressed, ebiten.KeyArrowDown) {
	//	g.mainCharacter.position.y += g.mainCharacter.speed
	}

	if isKeyArrowUp {
		g.mainCharacter.position.x += math.Cos(g.mainCharacter.angle*math.Pi/180.0)*g.mainCharacter.speed
		g.mainCharacter.position.y += math.Sin(g.mainCharacter.angle*math.Pi/180.0)*g.mainCharacter.speed
	}

	if contains(keysPressed, ebiten.KeyArrowLeft) && !isWeavingAndGoingUp(g.mainCharacter.weave.isWeaving, isKeyArrowUp) {
		g.mainCharacter.angle -= g.mainCharacter.angleSpeed
	}

	if contains(keysPressed, ebiten.KeyArrowRight) && !isWeavingAndGoingUp(g.mainCharacter.weave.isWeaving, isKeyArrowUp) {
		g.mainCharacter.angle += g.mainCharacter.angleSpeed
	}

	/*if contains(keysPressed, ebiten.KeySpace) {
		g.mainCharacter.weave.isWeaving = true
	} else {
		g.mainCharacter.weave.isWeaving = false
	}*/
	if contains(keysPressed, ebiten.KeySpace) {
		if g.bubbleShootCooldown {
			return
		}
		g.bubbles = append(g.bubbles, bubble{
			angle:      g.mainCharacter.angle,
			speed:      3,
			coordinate: coordinate{
				x: g.mainCharacter.position.x,
				y: g.mainCharacter.position.y,
			},
		})
		g.bubbleShootCooldown = true
	} 
}

func isWeavingAndGoingUp(isWeaving bool, isGoingUp bool) bool{
	return isWeaving && isGoingUp
}
