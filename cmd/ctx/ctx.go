package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.WithValue(context.Background(), "favorite-color", "blue")
	value := ctx.Value("favorite-color")
	fmt.Println(value)
}
