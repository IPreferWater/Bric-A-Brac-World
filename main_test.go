package main

import (
	"fmt"
	"testing")


func TestTimeConsuming(t *testing.T) {
	width := 420
	tileSize := 48

	fmt.Println(width/tileSize)
	fmt.Println(width%tileSize)
}
