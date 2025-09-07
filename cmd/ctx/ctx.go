package main

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.WithValue(context.Background(), favoriteColorKey, "blue")
	value := ctx.Value(favoriteColorKey)
	strValue, ok := value.(string)
	if !ok {
		fmt.Println("not a string")
		return
	}

	fmt.Println(strValue)
	fmt.Println(strings.HasPrefix(strValue, "b"))
}
