package main

import (
	"fmt"
	"math"
)

func (g *Game) CheckBubblesColisions() error {

	for i, b := range g.bubbles {

		var xLayer, yLayer int
		xLayer = int(b.coordinate.x) / 32
		yLayer = (int(b.coordinate.y)-16) / 32

		
		xLayerLeft := (int(b.coordinate.x-float64(g.boardXStart))-16) / 32
		xLayerRight := (int(b.coordinate.x-float64(g.boardXStart))+16) / 32
		fmt.Printf("xLayerLeft %d xLayerRight %d yLayer %d\n",xLayerLeft, xLayerRight, yLayer)

				//top collision
				if b.coordinate.y <= float64(g.boardYStart)+16 {
					g.bubblesLayer[yLayer][xLayerLeft] = &b
					g.bubbles = popBubble(g.bubbles, i)
					return nil
				}

		yLayerToCheck := g.bubblesLayer[yLayer]
		if yLayerToCheck[xLayerLeft] != nil || yLayerToCheck[xLayerRight] != nil{
			fmt.Println("boom!")
			g.bubblesLayer[yLayer+1][xLayerLeft] = &b
			g.bubbles = popBubble(g.bubbles, i)
		}


		//TODO need better idea
		maxLayer := 20
		//check collisions in other bubbles in boards
		for y := yLayer - 1; y <= yLayer+1; y++ {
			continue
			//out of the board
			if y < 0 || y > maxLayer {
				continue 
			}
			for x := xLayer - 1; x <= xLayer+1; x++ {
				//out of the board
				if x < 0 || x > maxLayer {
					continue
				}
				bubbleToCheck := g.bubblesLayer[yLayer][xLayer]
				if bubbleToCheck == nil {
					continue
				}
				// 32 is bubble size, -16 mean the center of the bubble
			dx := b.coordinate.x - float64((x*32))
			dy := b.coordinate.y - float64((y*32))
			distance := math.Sqrt((dx * dx) + (dy * dy))
			if distance < 64 {
				fmt.Printf("boom xLayer %d yLayer %d x %f y %f\n", xLayer, yLayer, b.coordinate.x, b.coordinate.y)
				g.bubblesLayer[yLayer][xLayerLeft] = &b
				g.bubbles = popBubble(g.bubbles, i)

				return nil
			}
			}
		}

		//check colision in board
		if b.coordinate.x <= float64(g.boardXStart) || b.coordinate.x >= float64(g.boardXStart+g.boardWidth) {
			//fmt.Printf("current angle is %f\n",g.bubbles[i].angle)
			g.bubbles[i].angle = g.bubbles[i].angle - 180 - (2 * g.bubbles[i].angle)
		}

		/*var dx = circle1.x - circle2.x
		var dy = circle1.y - circle2.y
		var distance = Math.sqrt(dx*dx + dy*dy)

		if distance < circle1.radius+circle2.radius {
			// collision détectée !
		}*/
		

	}

	return nil
}

func popBubble(bSlice []bubble, i int) []bubble {
	bSlice[i] = bSlice[len(bSlice)-1]
	return bSlice[:len(bSlice)-1]
}
