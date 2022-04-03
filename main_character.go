package main

type mainCharacter struct {
	position    coordinate
	destination coordinate
	speed float64
	weave weave
}

type weave struct {
	isWeaving bool
	coordinates []coordinate
}

type coordinate struct {
	x float64
	y float64
}