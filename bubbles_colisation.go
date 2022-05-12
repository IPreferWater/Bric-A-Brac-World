package main

func (g *Game) CheckBubblesColisions() error {

	for i, b := range g.bubbles {

		//check colision in table

		if b.coordinate.x<= float64(g.boardXStart) {
			g.bubbles[i].angle+=90
		}

		if b.coordinate.y <= 50 {
			var xLayer int
			xLayer = int(b.coordinate.x)/32
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
