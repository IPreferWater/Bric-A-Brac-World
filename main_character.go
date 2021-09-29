package main

type mainCharacter struct {
	position    coordinate
	destination coordinate
}

type coordinate struct {
	x float64
	y float64
}

//TODO if no need to moove, gresille
//TODO if g.worldSpeed to big, we could go over the destination
func (c *mainCharacter) animate(g *Game) {
	//moove
	if c.position.x == c.destination.x && c.position.y == c.destination.y {
		g.state = WaitWorldAction
		return
	}

	diffX := c.destination.x - c.position.x
	diffY := c.destination.y - c.position.y
	if diffX >= 0 {
		//its close enought
		c.position.x += g.worldSpeed
	} else {
		c.position.x -= g.worldSpeed
	}

	if diffY >= 0 {
		//its close enought
		c.position.y += g.worldSpeed
	} else {
		c.position.y -= g.worldSpeed
	}
}
