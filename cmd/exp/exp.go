package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func main() {
	err := b()
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Not found error caught", err)
	}
}

func a() error {
	return ErrNotFound
}

func b() error {
	err := a()
	if err != nil {
		return fmt.Errorf("b: %w", err)
	}

	return nil
}
