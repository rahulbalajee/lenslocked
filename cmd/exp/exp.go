package main

import (
	"fmt"

	"github.com/rahulbalajee/lenslocked/models"
)

func main() {
	gs := models.GalleryService{}

	fmt.Println(gs.Images(3))
}
