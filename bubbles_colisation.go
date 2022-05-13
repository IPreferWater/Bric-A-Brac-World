package main

import "fmt"

func (g *Game) CheckBubblesColisions() error {

	for i, b := range g.bubbles {

		//check colision in table

		if b.coordinate.x <= float64(g.boardXStart) ||  b.coordinate.x >= float64(g.boardXStart+g.boardWidth) {
			fmt.Printf("current angle is %f\n",g.bubbles[i].angle)
			g.bubbles[i].angle = g.bubbles[i].angle-180-(2*g.bubbles[i].angle)
		}

		

		if b.coordinate.y <= 50 {
			var xLayer int
			xLayer = int(b.coordinate.x) / 32
			g.bubblesLayer[0][xLayer] = &b
			g.bubbles = popBubble(g.bubbles, i)
		}
	}

	return nil
}

func popBubble(bSlice []bubble, i int) []bubble {
	bSlice[i] = bSlice[len(bSlice)-1]
	return bSlice[:len(bSlice)-1]
}
