package main

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.WithValue(context.Background(), favoriteColorKey, "blue")
	value := ctx.Value(favoriteColorKey)
	fmt.Println(value)
}
